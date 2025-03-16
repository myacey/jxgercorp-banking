package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/shared/ctxkeys"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/opentracing/opentracing-go"
)

func (h *Handler) AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.lg.Info(c.Request.Method)
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		authToken, err := c.Cookie("authToken")
		if err != nil {
			h.JSONError(c, cstmerr.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		req := &tokenpb.ValidateTokenRequest{
			Token: authToken,
		}
		resp, err := h.tokenSrv.ValidateToken(c, req)
		if err != nil {
			h.JSONError(c, cstmerr.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		username, valid := resp.Username, resp.Valid
		if !valid || username == "" {
			h.JSONError(c, cstmerr.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		h.lg.Debugw("auth success",
			"authToken", authToken,
			"username", username)
		c.Set(string(ctxkeys.UsernameKey), username)

		c.Next()
	}
}

func TracingMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := tracer.StartSpan(c.FullPath())
		tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		println(span)
		c.Next()
		span.Finish()
	}
}

// func RequestIDMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 	}
// }
