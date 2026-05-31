// @Title Register Routes
// @Description Handles user registration for the Expense Tracker API. Creates a new user account using name, email, and password.
// @Summary Registers user signup endpoint
// @Tags Authentication
// @Version 1.0.0
package routers

import (
	"expense-tracker-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

// RegisterRegisterRoutes registers all user registration route
func RegisterRegisterRoutes() {

	// POST /api/v1/auth/register
	// Creates a new user account in the system
	web.Router("/api/v1/auth/register", &controllers.RegisterController{}, "post:Register")
}