package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/server/web"
)

// Login Controller
type LoginController struct {
	web.Controller
}

// Request Structure
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


// Login Method
func (c *LoginController) Login() {
	var request loginRequest

	
	// Parsing JSON to Struct
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
		}
		c.ServeJSON()
		return
	}

	// Remove Extra Spaces
	email := strings.TrimSpace(request.Email)
	password := strings.TrimSpace(request.Password)

	// Check Authentication
	user, err := models.GetUserByEmail(email)
	if err != nil || user.Password != password {
		c.Ctx.Output.SetStatus(http.StatusUnauthorized)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Invalid email or password",
		}
		c.ServeJSON()
		return
	}

	// Login Success
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"user_id": user.ID,
			"name":    user.Name,
			"email":   user.Email,
		},
	}
	c.ServeJSON()
}
