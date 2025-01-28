package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/shared/ctxkeys"
)

func (h *Controller) CheckAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if strings.HasPrefix(authorization, "Bearer ") {
			splits := strings.Split(authorization, " ")
			if len(splits) != 2 {
				h.JSONError(c, cstmerr.New(http.StatusUnauthorized, "invalid auth", nil))
				return
			}
			authToken := splits[1]
			if authToken == "valid" {
				c.Set(string(ctxkeys.UsernameKey), "some_username") // TODO: gRPC
				c.Next()
			} else {
				h.lg.Debug("invalid token: ", authToken)
			}
		} else {
			h.lg.Debug("unknown token type: ", authorization)
		}
	}
}
