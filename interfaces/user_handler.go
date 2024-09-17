package interfaces

import (
	"strconv"

	"github.com/OmarBader7/web-service-jayeek/application"
	"github.com/OmarBader7/web-service-jayeek/infrastructure/auth"
	"github.com/OmarBader7/web-service-jayeek/pkg/language"
	"github.com/OmarBader7/web-service-jayeek/pkg/pagination"
	"github.com/OmarBader7/web-service-jayeek/pkg/response"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
)

// Users holds the user-related application interfaces
type Users struct {
	AuthService  auth.AuthServiceInterface
	TokenService auth.TokenInterface
	UserApp      application.UserApplicationInterface
	LocationApp  application.LocationApplicationInterface
}

// NewUsers returns a new instance of Users
func NewUsers(authService auth.AuthServiceInterface, tokenService auth.TokenInterface, userApp application.UserApplicationInterface, locationApp application.LocationApplicationInterface) *Users {
	return &Users{
		AuthService:  authService,
		TokenService: tokenService,
		UserApp:      userApp,
		LocationApp:  locationApp,
	}
}

// GetAllUsers retrieves a paginated list of all users.
func (u *Users) GetAllUsers(ctx *gin.Context) {
	// Get the desired page number from the query parameters.
	page := pagination.GetPage(ctx)

	// Set the number of items per page.
	perPage := 30

	// Get the user count from the user application service.
	count, err := u.UserApp.CountUsers()
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	// Get the users from the user application service.
	users, err := u.UserApp.GetAllUsers(page, perPage)
	if err != nil {
		response.SendInternalServerError(ctx, err.Error())
		return
	}

	if page <= 1 && len(users) <= 0 {
		response.SendOK(ctx, nil, ginI18n.MustGetMessage("No users found."))
		return
	}

	// Check if the page is valid
	if page <= 0 || (len(users) <= 0 && page*perPage > int(count)) {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("Page not found."))
		return
	}

	var userPublicData []interface{}

	for _, user := range users {
		userPublicData = append(userPublicData, user.PublicData(language.GetLanguage(ctx)))
	}

	// Build response data
	data := make(map[string]interface{})
	data["data"] = userPublicData
	data["current_page"] = page
	if page*perPage < int(count) {
		data["next_page"] = page + 1
	}
	data["total"] = count

	// Send the users as a response.
	response.SendOK(ctx, data, "")
}

// GetUserByID retrieves a single user by ID.
func (u *Users) GetUserByID(ctx *gin.Context) {
	// Parse the user ID from the URL parameter.
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	if err != nil {
		response.SendBadRequest(ctx, ginI18n.MustGetMessage("Invalid user ID."))
		return
	}

	// Get the user from the user application service.
	user, err := u.UserApp.GetUserByID(userID)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("User not found."))
		return
	}

	var currentUserID uint64 = 0

	// Extract the token metadata from the request
	metadata, err := u.TokenService.ExtractTokenMetadata(ctx.Request)
	if err == nil {
		// Fetch the authenticated user's ID from the auth service
		loggedInUserID, _ := u.AuthService.FetchAuth(metadata.AccessTokenUUID)
		if err == nil {
			if loggedInUserID == userID {
				currentUserID = loggedInUserID
			}
		}
	}

	userPublicData := user.PublicData(language.GetLanguage(ctx), currentUserID)

	// Send the user as a response.
	response.SendOK(ctx, userPublicData, "")
}

// GetUserByPhone retrieves a single user by Phone.
func (u *Users) GetUserByPhone(ctx *gin.Context) {
	// Parse the user phone from the URL parameter.
	userPhone := ctx.Param("phone")

	// Get the user from the user application service.
	user, err := u.UserApp.GetUserByPhone(userPhone)
	if err != nil {
		response.SendNotFound(ctx, ginI18n.MustGetMessage("User not found."))
		return
	}

	// Send the user as a response.
	response.SendOK(ctx, user.PublicData(language.GetLanguage(ctx)), "")
}
