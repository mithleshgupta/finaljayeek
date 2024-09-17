package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response struct defines the response format of the API
type Response struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// SendUnprocessableEntity sends a 422 Unprocessable Entity response with the provided errors and message
func SendUnprocessableEntity(c *gin.Context, errors interface{}, message string) {
	code := http.StatusUnprocessableEntity
	response := Response{
		Success: false,
		Status:  code,
		Errors:  errors,
		Message: message,
	}

	c.JSON(code, response)
}

// SendCreated sends a 201 Created response with the provided data and message
func SendCreated(c *gin.Context, data interface{}, message string) {
	code := http.StatusCreated
	response := Response{
		Success: true,
		Status:  code,
		Data:    data,
		Message: message,
	}

	c.JSON(code, response)
}

// SendOK sends a 200 Ok response with the provided data and message
func SendOK(c *gin.Context, data interface{}, message string) {
	code := http.StatusCreated
	response := Response{
		Success: true,
		Status:  code,
		Data:    data,
		Message: message,
	}

	c.JSON(code, response)
}

// SendBadRequest sends a 400 Bad Request response with the provided message
func SendBadRequest(c *gin.Context, message string) {
	code := http.StatusBadRequest
	response := Response{
		Success: false,
		Status:  code,
		Message: message,
	}

	c.JSON(code, response)
}

// SendInternalServerError sends a 500 Internal Server Error response with the provided message
func SendInternalServerError(c *gin.Context, message string) {
	code := http.StatusInternalServerError
	response := Response{
		Success: false,
		Status:  code,
		Message: message,
	}

	c.JSON(code, response)
}

// SendOK sends a 401 Unauthorized response with the provided message
func SendUnauthorized(c *gin.Context, message string) {
	code := http.StatusUnauthorized
	response := Response{
		Success: false,
		Status:  code,
		Message: message,
	}

	c.JSON(code, response)
}

// SendNotFound sends a 404 Not Found response with the provided message
func SendNotFound(c *gin.Context, message string) {
	code := http.StatusNotFound
	response := Response{
		Success: false,
		Status:  code,
		Message: message,
	}

	c.JSON(code, response)
}

// SendForbidden sends a 403 Forbidden response with the provided message
func SendForbidden(c *gin.Context, message string) {
	code := http.StatusForbidden
	response := Response{
		Success: false,
		Status:  code,
		Message: message,
	}

	c.JSON(code, response)
}
