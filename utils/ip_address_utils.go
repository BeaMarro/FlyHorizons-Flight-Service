package utils

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var WhitelistedIPs []string

func LoadWhitelistedIPs() {
	ips := os.Getenv("WHITELISTED_IPS")
	WhitelistedIPs = strings.Split(ips, ",")
}

func IPWhitelistingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetIPAddress(c.Request)

		for _, allowedIP := range WhitelistedIPs {
			if strings.TrimSpace(ip) == strings.TrimSpace(allowedIP) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "IP address not allowed"})
	}
}

func GetIPAddress(r *http.Request) string {
	// X-Forwarded-For first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
