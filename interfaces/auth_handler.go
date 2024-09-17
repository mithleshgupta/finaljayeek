package interfaces

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/domain/entity"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/chat"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	"github.com/OmarBader7/web-service-jayeek/pkg/validator"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	AuthService          auth.AuthServiceInterface
	TokenService         auth.TokenInterface
	ChatService          chat.ChatServiceInterface
	UserApp              application.UserApplicationInterface
	LocationApp          application.LocationApplicationInterface
	OrderApp             application.OrderApplicationInterface
	PasswordResetApp     application.PasswordResetApplicationInterface
	PhoneVerificationApp application.PhoneVerificationApplicationInterface
}

func NewAuth(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, chatService chat.ChatServiceInterface, userApp application.UserApplicationInterface, locationApp application.LocationApplicationInterface, orderApp application.OrderApplicationInterface, passwordReset application.PasswordResetApplicationInterface, phoneVerification application.PhoneVerificationApplicationInterface) *Auth {
	return &Auth{
		AuthService:          authService,
		TokenService:         tokenService,
		ChatService:          chatService,
		UserApp:              userApp,
		LocationApp:          locationApp,
		OrderApp:             orderApp,
		PasswordResetApp:     passwordReset,
		PhoneVerificationApp: phoneVerification,
	}
}

func (a *Auth) Register(c *gin.Context) {
	var user entity.User

	c.ShouldBindJSON(&user)

	// Validate all fields except the ones passed in.
	if errors, _ := validator.ValidateExcept(c, &user, "Location", "Role"); errors != nil {
		response.SendUnprocessableEntity(c, errors, "")
		return
	}

	// Get a location by its ID.
	if _, err := a.LocationApp.GetLocationByID(user.LocationID); err != nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("Location not found."))
		return
	}

	// Check if a user exists by its phone.
	if isUserExists, _ := a.UserApp.UserWithFieldExists("phone", user.Phone); isUserExists {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("The phone number you've entered already exists with another account."))
		return
	}

	// Create a new user.
	user.Role = "user"

	user.AddSetting("is_available", true)
	user.AddSetting("is_dark_mode", false)
	user.AddSetting("is_24_hour_format", false)

	if _, err := a.UserApp.CreateUser(&user); err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	ts, tErr := a.TokenService.CreateToken(user.ID)
	if tErr != nil {
		response.SendUnprocessableEntity(c, nil, tErr.Error())
		return
	}

	if err := a.AuthService.CreateAuth(user.ID, ts); err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	// Get the orders from the order application service.
	orders, err := a.OrderApp.GetAllOrdersByRecipientPhoneNumber(user.Phone)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	for _, order := range orders {
		order.RecipientID = user.ID

		_, err = a.OrderApp.UpdateOrderByID(order.ID, &order)
		if err != nil {
			response.SendInternalServerError(c, err.Error())
			return
		}

		if order.Status == entity.OrderAcceptedStatus || order.Status == entity.ShipmentPickedUpStatus {
			err = a.ChatService.AddMember(fmt.Sprintf("order-%d", order.ID), fmt.Sprintf("client-%d", user.ID))

			if err != nil {
				response.SendInternalServerError(c, ginI18n.MustGetMessage("Failed to add members from the chat channel."))
				return
			}
		}
	}

	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["data"] = user.PublicData(language.GetLanguage(c))

	response.SendOK(c, userData, "")
}

// Login handles the user login request
func (a *Auth) Login(c *gin.Context) {
	var user entity.User

	if err := c.ShouldBindJSON(&user); err != nil {
		response.SendBadRequest(c, err.Error())
		return
	}

	if errors, _ := validator.ValidatePartial(c, &user, "Phone", "Password"); errors != nil {
		response.SendUnprocessableEntity(c, errors, "")
		return
	}

	// Check if a user exists by its phone.
	if isUserExists, _ := a.UserApp.UserWithFieldExists("phone", user.Phone); !isUserExists {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("These credentials do not match our records."))
		return
	}

	u, _ := a.UserApp.GetUserByPhoneAndPassword(user.Phone, user.Password)
	if u == nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("The provided password is incorrect."))
		return
	}

	ts, tErr := a.TokenService.CreateToken(u.ID)
	if tErr != nil {
		response.SendUnprocessableEntity(c, nil, tErr.Error())
		return
	}

	if err := a.AuthService.CreateAuth(u.ID, ts); err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["data"] = u.PublicData(language.GetLanguage(c))

	response.SendOK(c, userData, "")
}

// Logout handles the user logout request
func (a *Auth) Logout(c *gin.Context) {
	metadata, err := a.TokenService.ExtractTokenMetadata(c.Request)
	if err != nil {
		response.SendUnauthorized(c, ginI18n.MustGetMessage("Unauthorized"))
		return
	}

	// Verify the existence and validity of the access token. If it is valid, delete both the access token and the refresh token.
	if err := a.AuthService.DeleteTokens(metadata); err != nil {
		response.SendUnauthorized(c, err.Error())
		return
	}

	response.SendOK(c, nil, ginI18n.MustGetMessage("You have been successfully logged out."))
}

func (a *Auth) Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		response.SendUnprocessableEntity(c, nil, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	// Verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		response.SendUnauthorized(c, err.Error())
		return
	}

	// Check if the token is valid
	if !token.Valid {
		response.SendUnauthorized(c, "invalid token")
		return
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		response.SendUnauthorized(c, "invalid token claims")
		return
	}

	refreshUuid, ok := claims["refresh_token_uuid"].(string)
	if !ok {
		response.SendUnprocessableEntity(c, nil, "cannot get refresh_token_uuid")
		return
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		response.SendUnprocessableEntity(c, nil, "error occurred")
		return
	}

	// Delete the previous Refresh Token
	if err := a.AuthService.DeleteRefreshToken(refreshUuid); err != nil {
		response.SendUnauthorized(c, "unauthorized")
		return
	}

	// Create new pairs of refresh and access tokens
	ts, createErr := a.TokenService.CreateToken(userId)
	if createErr != nil {
		response.SendForbidden(c, createErr.Error())
		return
	}

	// Save the tokens metadata to Redis
	if err := a.AuthService.CreateAuth(userId, ts); err != nil {
		response.SendForbidden(c, err.Error())
		return
	}

	user, _ := a.UserApp.GetUserByID(userId)
	if user == nil {
		response.SendUnprocessableEntity(c, nil, ginI18n.MustGetMessage("error occurred"))
		return
	}

	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["data"] = user.PublicData(language.GetLanguage(c))

	response.SendCreated(c, userData, "")
}

func (a *Auth) SendPasswordResetCode(c *gin.Context) {
	var passwordResetRequest auth.PasswordResetRequest

	if err := c.ShouldBindJSON(&passwordResetRequest); err != nil {
		response.SendBadRequest(c, err.Error())
		return
	}

	// Validate all fields except the ones passed in.
	if errors, _ := validator.ValidateExcept(c, &passwordResetRequest, "VerificationCode", "Password", "PasswordConfermation"); errors != nil {
		response.SendUnprocessableEntity(c, errors, "")
		return
	}

	// Get the user from the user application service
	user, err := a.UserApp.GetUserByPhone(passwordResetRequest.Phone)
	if err != nil {
		response.SendNotFound(c, ginI18n.MustGetMessage("The phone number is not associated with any account."))
		return
	}

	var passwordReset entity.PasswordReset

	passwordReset.UserID = user.ID
	passwordReset.VerificationCode = "113438"
	passwordReset.ExpiresAt = time.Now().Add(time.Hour)

	_, err = a.PasswordResetApp.CreatePasswordReset(&passwordReset)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	response.SendOK(c, nil, ginI18n.MustGetMessage("Password reset code sent successfully."))
}

func (a *Auth) VerifyPasswordResetCode(c *gin.Context) {
	var passwordResetRequest auth.PasswordResetRequest

	if err := c.ShouldBindJSON(&passwordResetRequest); err != nil {
		response.SendBadRequest(c, err.Error())
		return
	}

	// Validate all fields except the ones passed in.
	if errors, _ := validator.ValidateExcept(c, &passwordResetRequest, "Password", "PasswordConfermation"); errors != nil {
		response.SendUnprocessableEntity(c, errors, "")
		return
	}

	// Get the user from the user application service
	user, err := a.UserApp.GetUserByPhone(passwordResetRequest.Phone)
	if err != nil {
		response.SendNotFound(c, ginI18n.MustGetMessage("The phone number is not associated with any account."))
		return
	}

	passwordReset, err := a.PasswordResetApp.GetPasswordResetByUserIDAndVerificationCode(user.ID, passwordResetRequest.VerificationCode)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid code."))
		return
	}

	if time.Now().After(passwordReset.ExpiresAt) {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Password reset request has expired."))
		return
	}

	response.SendOK(c, nil, "")
}

func (a *Auth) Reset(c *gin.Context) {
	var passwordResetRequest auth.PasswordResetRequest

	if err := c.ShouldBindJSON(&passwordResetRequest); err != nil {
		response.SendBadRequest(c, err.Error())
		return
	}

	// Validate all fields except the ones passed in.
	if errors, _ := validator.ValidateExcept(c, &passwordResetRequest); errors != nil {
		response.SendUnprocessableEntity(c, errors, "")
		return
	}

	// Get the user from the user application service
	user, err := a.UserApp.GetUserByPhone(passwordResetRequest.Phone)
	if err != nil {
		response.SendNotFound(c, ginI18n.MustGetMessage("The phone number is not associated with any account."))
		return
	}

	passwordReset, err := a.PasswordResetApp.GetPasswordResetByUserIDAndVerificationCode(user.ID, passwordResetRequest.VerificationCode)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid code."))
		return
	}

	if time.Now().After(passwordReset.ExpiresAt) {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Password reset code has expired."))
		return
	}

	// Validate password and confirmation match
	if passwordResetRequest.Password != passwordResetRequest.PasswordConfirmation {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Password and confirmation do not match."))
		return
	}

	user.Password = passwordResetRequest.Password

	updatedUser, err := a.UserApp.UpdateUserByID(user.ID, user)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	passwordReset.Used = true

	// Mark the password reset as used
	_, err = a.PasswordResetApp.UpdatePasswordResetByID(passwordReset.ID, passwordReset)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	ts, tErr := a.TokenService.CreateToken(user.ID)
	if tErr != nil {
		response.SendUnprocessableEntity(c, nil, tErr.Error())
		return
	}

	if err := a.AuthService.CreateAuth(user.ID, ts); err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	userData["data"] = updatedUser.PublicData(language.GetLanguage(c))

	response.SendOK(c, userData, "")
}

func (a *Auth) SendPhoneVerificationCode(c *gin.Context) {
	var phoneVerificationRequest auth.PhoneVerificationRequest

	if err := c.ShouldBindJSON(&phoneVerificationRequest); err != nil {
		response.SendBadRequest(c, err.Error())
		return
	}

	// Validate all fields except the ones passed in.
	if errors, _ := validator.ValidateExcept(c, &phoneVerificationRequest, "Code"); errors != nil {
		response.SendUnprocessableEntity(c, errors, "")
		return
	}

	var phoneVerification entity.PhoneVerification

	phoneVerification.Phone = phoneVerificationRequest.Phone
	phoneVerification.Code = "113438"
	phoneVerification.ExpiresAt = time.Now().Add(time.Hour)

	_, err := a.PhoneVerificationApp.CreatePhoneVerification(&phoneVerification)
	if err != nil {
		response.SendInternalServerError(c, err.Error())
		return
	}

	response.SendOK(c, nil, ginI18n.MustGetMessage("Phone verification code sent successfully."))
}

func (a *Auth) VerifyPhoneVerificationCode(c *gin.Context) {
	var phoneVerificationRequest auth.PhoneVerificationRequest

	if err := c.ShouldBindJSON(&phoneVerificationRequest); err != nil {
		response.SendBadRequest(c, err.Error())
		return
	}

	// Validate all fields except the ones passed in.
	if errors, _ := validator.ValidateExcept(c, &phoneVerificationRequest); errors != nil {
		response.SendUnprocessableEntity(c, errors, "")
		return
	}

	phoneVerification, err := a.PhoneVerificationApp.GetPhoneVerificationByPhoneAndCode(phoneVerificationRequest.Phone, phoneVerificationRequest.Code)
	if err != nil {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Invalid code."))
		return
	}

	if time.Now().After(phoneVerification.ExpiresAt) {
		response.SendBadRequest(c, ginI18n.MustGetMessage("Phone verification code has expired."))
		return
	}

	response.SendOK(c, nil, "")
}
