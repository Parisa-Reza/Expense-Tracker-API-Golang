// @Title Health Check Routes
// @Description Provides a simple endpoint to check if the Expense Tracker API is running and healthy.
// @Summary Health check endpoint
// @Tags System
// @Version 1.0.0

package routers

import (
	"expense-tracker-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

// RegisterHealthRoutes registers all health-related API routes.
//
// Available Routes:
//   GET /api/v1/health → Checks if the server is running
func RegisterHealthRoutes() {

	// Health check endpoint
	// Used to verify that the API server is working properly
	web.Router("/api/v1/health", &controllers.HealthController{}, "get:Check")
}