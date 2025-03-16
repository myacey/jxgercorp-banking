package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/ctxkeys"
)

func ProxyHandler(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":    http.StatusInternalServerError,
				"message": "invalid target",
			}})
			return
		}

		username, exists := c.Get(string(ctxkeys.UsernameKey)) // move ctx username to HTTP Header
		if exists {
			if usernameStr, ok := username.(string); ok {
				c.Request.Header.Set("X-User-Username", usernameStr) // set username when authentificated
			}
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
