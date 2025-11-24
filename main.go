package main

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, clientIPFromRequest(c))
	})

	r.Run(":8080")
}

// clientIPFromRequest attempts to return the real client IP when the app
// is running behind a reverse proxy (like Azure App Service). It checks
// `X-Forwarded-For` (first entry), `X-Real-IP`, then falls back to
// the request RemoteAddr. Candidate IPs are validated with net.ParseIP.
func clientIPFromRequest(c *gin.Context) string {
	// X-Forwarded-For may contain a comma-separated list; client is first
	if xff := c.Request.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// X-Real-IP (some proxies set this)
	if xr := strings.TrimSpace(c.Request.Header.Get("X-Real-IP")); xr != "" {
		if net.ParseIP(xr) != nil {
			return xr
		}
	}

	// Fallback to RemoteAddr (may be proxy IP)
	if host, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		if net.ParseIP(host) != nil {
			return host
		}
	}

	// Last resort: gin's helper (behaviour depends on Gin trusted proxies)
	return c.ClientIP()
}
