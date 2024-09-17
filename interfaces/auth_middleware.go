// Package interfaces contains the authentication middleware for the web service.
package interfaces

import (
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a gin middleware that handles the authorization of incoming requests.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the request has a valid token.
		err := auth.TokenValid(ctx.Request)
		if err != nil {
			// If the token is not valid, send a 401 Unauthorized response.
			response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
			ctx.Abort()
			return
		}
		// If the token is valid, continue with processing the request.
		ctx.Next()
	}
}
