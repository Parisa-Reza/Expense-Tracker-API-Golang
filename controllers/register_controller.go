package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// RegisterController handles user registration requests.
type RegisterController struct {
	web.Controller
}

// Request Structure
type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register creates a new user account after validating the request.
func (c *RegisterController) Register() {

	// Create Empty Request Object
	var request registerRequest

	// Parsing JSON to Struct
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		beego.Warn("invalid register request body:", err)
		writeRegisterJSON(&c.Controller, http.StatusBadRequest, false, "Invalid request body", nil)
		return
	}

	// Remove Extra Spaces
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	// Check Required Fields
	if request.Name == "" {
		writeRegisterJSON(&c.Controller, http.StatusBadRequest, false, "Name is required", nil)
		return
	}
	if request.Email == "" {
		writeRegisterJSON(&c.Controller, http.StatusBadRequest, false, "Email is required", nil)
		return
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		writeRegisterJSON(&c.Controller, http.StatusBadRequest, false, "Email must be valid", nil)
		return
	}
	if request.Password == "" {
		writeRegisterJSON(&c.Controller, http.StatusBadRequest, false, "Password is required", nil)
		return
	}
	if len(request.Password) < 6 {
		writeRegisterJSON(&c.Controller, http.StatusBadRequest, false, "Password must be at least 6 characters", nil)
		return
	}

	// Check Whether Email Already Exists
	if _, err := models.GetUserByEmail(request.Email); err == nil {
		writeRegisterJSON(&c.Controller, http.StatusConflict, false, "Email already exists", nil)
		return
	} else if !errors.Is(err, models.ErrUserNotFound) {
		beego.Error("failed to check email:", err)
		writeRegisterJSON(&c.Controller, http.StatusInternalServerError, false, "Could not check email", nil)
		return
	}

	// Create User Object
	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := models.CreateUser(user); err != nil {
		beego.Error("failed to register user:", err)
		writeRegisterJSON(&c.Controller, http.StatusInternalServerError, false, "Could not register user", nil)
		return
	}

	// Registration Success
	beego.Info("user registered:", user.ID)
	writeRegisterJSON(&c.Controller, http.StatusCreated, true, "User registered successfully", nil)
}

func writeRegisterJSON(controller *web.Controller, statusCode int, success bool, message string, data interface{}) {
	controller.Ctx.Output.SetStatus(statusCode)
	controller.Data["json"] = map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	}
	controller.ServeJSON()
}
