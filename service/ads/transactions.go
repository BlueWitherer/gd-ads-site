package ads

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"service/database"
	"service/log"
)

// request/response shapes
type createBoostReq struct {
	UserID string `json:"user_id"`
	Boosts int    `json:"boosts"`
}

type createBoostResp struct {
	OrderID    string `json:"order_id"`
	ApproveURL string `json:"approve_url"`
}

type captureReq struct {
	OrderID string `json:"order_id"`
}

// prototyped PayPal integration for ad boost purchases
func init() {
	http.HandleFunc("/api/ads/boost/create", createBoostOrderHandler)
	http.HandleFunc("/api/ads/boost/capture", captureBoostOrderHandler)
}

// config helpers (envs left for business account info)
func paypalAPIBase() string {
	if v := os.Getenv("PAYPAL_API_BASE"); v != "" {
		return v
	}
	// default to sandbox
	return "https://api-m.sandbox.paypal.com"
}

func paypalClientCredentials() (string, string) {
	return os.Getenv("PAYPAL_CLIENT_ID"), os.Getenv("PAYPAL_SECRET")
}

func pricePerBoost() float64 {
	if v := os.Getenv("PRICE_PER_BOOST"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}

	return 1.00
}

func returnURL() string {
	if v := os.Getenv("PAYPAL_RETURN_URL"); v != "" {
		return v
	}

	return "https://example.com/return"
}

func cancelURL() string {
	if v := os.Getenv("PAYPAL_CANCEL_URL"); v != "" {
		return v
	}

	return "https://example.com/cancel"
}

// Obtain PayPal access token (client credentials)
func getPayPalAccessToken() (string, error) {
	clientID, secret := paypalClientCredentials()
	if clientID == "" || secret == "" {
		return "", fmt.Errorf("missing PayPal credentials")
	}

	reqBody := "grant_type=client_credentials"
	req, err := http.NewRequest("POST", paypalAPIBase()+"/v1/oauth2/token", bytes.NewBufferString(reqBody))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(clientID, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: %s", string(b))
	}

	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.AccessToken, nil
}

// Create a PayPal order with user info encoded into purchase_unit.custom_id (userID:boosts)
func createPayPalOrder(userID string, boosts int, accessToken string) (orderID string, approveURL string, err error) {
	price := pricePerBoost() * float64(boosts)
	amountStr := fmt.Sprintf("%.2f", price)

	reqBody := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"amount": map[string]string{
					"currency_code": "USD",
					"value":         amountStr,
				},
				// encode user id so capture can credit the user
				"custom_id": fmt.Sprintf("%s:%d", userID, boosts),
			},
		},
		"application_context": map[string]string{
			"return_url": returnURL(),
			"cancel_url": cancelURL(),
		},
	}

	b, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", paypalAPIBase()+"/v2/checkout/orders", bytes.NewReader(b))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("paypal create order failed: %s", string(body))
	}

	var out struct {
		ID    string `json:"id"`
		Links []struct {
			Href   string `json:"href"`
			Rel    string `json:"rel"`
			Method string `json:"method"`
		} `json:"links"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}

	var approve string
	for _, l := range out.Links {
		if l.Rel == "approve" {
			approve = l.Href
			break
		}
	}

	return out.ID, approve, nil
}

// Capture order and extract custom_id (userID:boosts)
func capturePayPalOrder(orderID, accessToken string) (userID string, boosts uint, err error) {
	// Capture
	req, err := http.NewRequest("POST", paypalAPIBase()+"/v2/checkout/orders/"+orderID+"/capture", nil)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("paypal capture failed: %s", string(body))
	}

	// Parse response to find custom_id
	var out struct {
		PurchaseUnits []struct {
			CustomID string `json:"custom_id"`
		} `json:"purchase_units"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", 0, err
	}

	if len(out.PurchaseUnits) == 0 || out.PurchaseUnits[0].CustomID == "" {
		return "", 0, fmt.Errorf("missing custom_id in capture response")
	}

	// custom_id formatted as "userID:boosts"
	parts := out.PurchaseUnits[0].CustomID
	seg := strings.SplitN(parts, ":", 2)
	if len(seg) != 2 {
		return "", 0, fmt.Errorf("invalid custom_id format: %s", parts)
	}
	user := seg[0]
	bcount, err := strconv.ParseUint(seg[1], 10, 32)
	if err != nil {
		return "", 0, fmt.Errorf("invalid boost count in custom_id: %s", seg[1])
	}

	return user, uint(bcount), nil
}

// Handler: create payPal order
func createBoostOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Create boost order called")
	header := w.Header()
	header.Set("Access-Control-Allow-Methods", "POST")
	header.Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createBoostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Invalid create request: %s", err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// require a user id to credit boosts to and a positive boost count
	if req.UserID == "" || req.Boosts <= 0 {
		http.Error(w, "missing user_id or boosts", http.StatusBadRequest)
		return
	}

	accessToken, err := getPayPalAccessToken()
	if err != nil {
		log.Error("Failed to get PayPal token: %s", err.Error())
		http.Error(w, "payment provider error", http.StatusInternalServerError)
		return
	}

	orderID, approveURL, err := createPayPalOrder(req.UserID, req.Boosts, accessToken)
	if err != nil {
		log.Error("Failed to create PayPal order: %s", err.Error())
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	resp := createBoostResp{OrderID: orderID, ApproveURL: approveURL}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Handler: capture payPal order and credit boosts to user account
func captureBoostOrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Capture boost order called")
	header := w.Header()
	header.Set("Access-Control-Allow-Methods", "POST")
	header.Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req captureReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Invalid capture request: %s", err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.OrderID == "" {
		http.Error(w, "missing order_id", http.StatusBadRequest)
		return
	}

	accessToken, err := getPayPalAccessToken()
	if err != nil {
		log.Error("Failed to get PayPal token: %s", err.Error())
		http.Error(w, "payment provider error", http.StatusInternalServerError)
		return
	}

	userID, boosts, err := capturePayPalOrder(req.OrderID, accessToken)
	if err != nil {
		log.Error("Failed to capture order: %s", err.Error())
		http.Error(w, "failed to capture order", http.StatusInternalServerError)
		return
	}

	if boosts <= 0 || userID == "" {
		http.Error(w, "invalid order payload", http.StatusBadRequest)
		return
	}

	// credit boosts to the user's account using the preexisting function
	if err := database.AddBoostsToUser(userID, boosts); err != nil {
		log.Error("Failed to apply boosts to user %s: %s", userID, err.Error())
		http.Error(w, "failed to apply boost", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Boost order handled successfully")
}
