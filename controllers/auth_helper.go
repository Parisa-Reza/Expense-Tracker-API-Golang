package controllers

import (
	"net/http"
	"strconv"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web/context"
)

// RequireUserID validates the X-User-ID request header.
func RequireUserID(ctx *context.Context) bool {
	_, ok := GetAuthenticatedUserID(ctx)
	return ok
}

// GetAuthenticatedUserID validates the X-User-ID request header and returns the user ID.
func GetAuthenticatedUserID(ctx *context.Context) (int, bool) {

	// Reads HTTP header "X-User-ID" from request and stores it as a string
	userIDHeader := ctx.Input.Header("X-User-ID")

	// Convert the string to an integer.
	userID, err := strconv.Atoi(userIDHeader)

	//  here validation checks if the conversion was successful and if the user ID is valid according to the user model. If either check fails, we return an unauthorized response.
	if err != nil || !models.IsValidUserID(userID) {
		beego.Warn("unauthorized request with X-User-ID:", userIDHeader)
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
			"data":    nil,
		}, true, false) // here false, false means: do not indent JSON and do not escape HTML characters
		return 0, false
	}

	// the user ID is valid. We can return it along with true to indicate successful authentication.
	return userID, true
}
