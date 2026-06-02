// @Title Login Routes
// @Description Handles user authentication (login) for the Expense Tracker API.
// @Summary Registers login endpoint for user authentication.
// @Tags Authentication
// @Version 1.0.0
package routers

import (
	"expense-tracker-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

// RegisterLoginRoutes registers all login-related routes
func RegisterLoginRoutes() {

	// POST /api/v1/auth/login
	// Authenticates a user using email and password
	web.Router("/api/v1/auth/login", &controllers.LoginController{}, "post:Login")
}
