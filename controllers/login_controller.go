package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// LoginController handles user login requests.
type LoginController struct {
	web.Controller
}

// Request Structure
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates a user with email and password.
func (c *LoginController) Login() {
	var request loginRequest

	// Parsing JSON to Struct
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		beego.Warn("invalid login request body:", err)
		writeLoginJSON(&c.Controller, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	// Remove Extra Spaces
	email := strings.TrimSpace(request.Email)
	password := strings.TrimSpace(request.Password)

	// Check Authentication
	user, err := models.GetUserByEmail(email)
	if err != nil || user.Password != password {
		beego.Warn("failed login attempt for email:", email)
		writeLoginJSON(&c.Controller, http.StatusUnauthorized, false, "Invalid email or password", nil)
		return
	}

	// Login Success
	beego.Info("user logged in:", user.ID)
	writeLoginJSON(&c.Controller, http.StatusOK, true, "Login successful", map[string]interface{}{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	})
}

func writeLoginJSON(controller *web.Controller, statusCode int, success bool, message string, data interface{}) {
	controller.Ctx.Output.SetStatus(statusCode)
	controller.Data["json"] = map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}
	controller.ServeJSON()
}
