package controllers

import (
	"net/http"
	"strconv"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/server/web/context"
)

func RequireUserID(ctx *context.Context) bool {
	userIDHeader := ctx.Input.Header("X-User-ID")
	userID, err := strconv.Atoi(userIDHeader)
	if err != nil || !models.IsValidUserID(userID) {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		}, false, false)
		return false
	}

	return true
}
