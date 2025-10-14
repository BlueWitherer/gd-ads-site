package access

import (
	"fmt"
	"net/http"
	"strings"

	"service/log"
)

func isInternal(ip string) bool {
	return strings.HasPrefix(ip, "127.") ||
		strings.HasPrefix(ip, "::1") ||
		strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.")
}

// Check if the request was received internally (for testing with sensitive data)
func Restrict(ip string) (int, error) {
	log.Debug("Checking internal address %s", ip)

	if isInternal(ip) {
		log.Error("Address %s forbidden for use of internal API", ip)
		return http.StatusForbidden, fmt.Errorf("Forbidden")
	} else {
		log.Info("Address %s authorized for use of internal API", ip)
		return http.StatusOK, nil
	}
}
