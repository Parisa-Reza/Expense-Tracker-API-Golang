package controllers

import (
	"net/http"
	"strconv"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/server/web/context"
)

// RequireUserID validates the X-User-ID request header.
func RequireUserID(ctx *context.Context) bool {
	_, ok := GetAuthenticatedUserID(ctx)
	return ok
}

// GetAuthenticatedUserID validates the X-User-ID request header and returns the user ID.
func GetAuthenticatedUserID(ctx *context.Context) (int, bool) {
	userIDHeader := ctx.Input.Header("X-User-ID")
	userID, err := strconv.Atoi(userIDHeader)
	if err != nil || !models.IsValidUserID(userID) {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		}, false, false)
		return 0, false
	}

	return userID, true
}
