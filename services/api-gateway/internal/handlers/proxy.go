package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
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

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Set("Access-Control-Allow-Origin", "http://localhost:8080")
			resp.Header.Set("Access-Control-Allow-Credentials", "true")
			return nil
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
