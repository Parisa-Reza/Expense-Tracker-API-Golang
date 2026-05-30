package routers

import (
	"expense-tracker-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

// RegisterHealthRoutes registers all health-related API routes.
//
// Routes registered:
//
//	GET /api/v1/health → HealthController.Check
func RegisterHealthRoutes() {
	web.Router("/api/v1/health", &controllers.HealthController{}, "get:Check")
}