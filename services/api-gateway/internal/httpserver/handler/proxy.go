package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var address = map[string]string{
	"http://localhost:8081": "user-service",
	"http://localhost:8082": "transfer-service",
}

func (h *Handler) ProxyHandler(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		toService := "unknown"
		if s, ok := address[target]; ok {
			toService = s
		}
		ctx, span := h.tracer.Start(c.Request.Context(), "proxy: "+toService)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
				"code":    http.StatusInternalServerError,
				"message": "invalid target",
			}})
			return
		}

		username, exists := c.Get(CtxKeyUsername) // move ctx username to HTTP Header
		if exists {
			if usernameStr, ok := username.(string); ok {
				c.Request.Header.Set(HeaderUsername, usernameStr) // set username when authentificated
			}
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.Host = remote.Host

			otel.GetTextMapPropagator().Inject(c.Request.Context(), propagation.HeaderCarrier(req.Header))
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
