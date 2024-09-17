package interfaces

import (
	"fmt"
	"log"

	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/chat"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

type Chat struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	ChatService  chat.ChatServiceInterface
}

// NewChat creates and returns a new instance of Chat.
func NewChat(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, chatService chat.ChatServiceInterface) *Chat {
	return &Chat{
		AuthService:  authService,
		TokenService: tokenService,
		ChatService:  chatService,
	}
}

func (c *Chat) GetClientStreamToken(ctx *gin.Context) {
	// Extract the token metadata from the request
	metadata, err := c.TokenService.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := c.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	token, err := c.ChatService.GetStreamToken(fmt.Sprintf("client-%d", userID))

	if err != nil {
		log.Fatalf("Error recreating user: %v", err)
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to retrieve chat token."))
		return
	}

	// Prepare the response data.
	tokenData := make(map[string]interface{})
	tokenData["token"] = token.Token

	// Send the response.
	response.SendOK(ctx, tokenData, "")
}

func (c *Chat) GetDriverStreamToken(ctx *gin.Context) {
	// Extract the token metadata from the request
	metadata, err := c.TokenService.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Fetch the authenticated user's ID from the auth service
	userID, err := c.AuthService.FetchAuth(metadata.AccessTokenUUID)
	if err != nil {
		response.SendUnauthorized(ctx, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	token, err := c.ChatService.GetStreamToken(fmt.Sprintf("driver-%d", userID))

	if err != nil {
		log.Fatalf("Error recreating user: %v", err)
		response.SendInternalServerError(ctx, ginI18n.MustGetMessage("Failed to retrieve chat token."))
		return
	}

	// Prepare the response data.
	tokenData := make(map[string]interface{})
	tokenData["token"] = token.Token

	// Send the response.
	response.SendOK(ctx, tokenData, "")
}
