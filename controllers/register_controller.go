package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/server/web"
)

// Register Controller
type RegisterController struct {
	web.Controller
}

// Request Structure
type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register Method
func (c *RegisterController) Register() {

	// Create Empty Request Object
	var request registerRequest

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
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	// Check Required Fields
	if request.Name == "" || request.Email == "" || request.Password == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Name, email, and password are required",
		}
		c.ServeJSON()
		return
	}

	// Check Whether Email Already Exists
	if _, err := models.GetUserByEmail(request.Email); err == nil {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Email already exists",
		}
		c.ServeJSON()
		return
	} else if !errors.Is(err, models.ErrUserNotFound) {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Could not check email",
		}
		c.ServeJSON()
		return
	}

	// Create User Object
	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	// Save User
	
	if err := models.CreateUser(user); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "Could not register user",
		}
		c.ServeJSON()
		return
	}

	// Registration Success
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
	}

	// Return JSON Response to Client
	c.ServeJSON()
}
