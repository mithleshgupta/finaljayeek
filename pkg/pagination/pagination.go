package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPage returns the desired page number from the query parameters.
func GetPage(ctx *gin.Context) int {
	pageStr := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	return page
}
