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

func Restrict(ip string) (int, error) {
	log.Debug("Checking internal address " + ip)

	if isInternal(ip) {
		log.Error("Address " + ip + " forbidden for use of internal API")
		return http.StatusForbidden, fmt.Errorf("Forbidden")
	} else {
		log.Info("Address " + ip + " authorized for use of internal API")
		return http.StatusOK, nil
	}
}
