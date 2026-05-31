// @APIVersion 1.0.0
// @Title Expense Tracker APIs
// @Description A RESTful API for managing user authentication and expense tracking, including user registration, login, and system health checks.
// @BasePath /api/v1
package routers

func init() {

	// Register all application routes

	// Health check endpoints (system status)
	RegisterHealthRoutes()

	// User authentication routes (register, login)
	RegisterRegisterRoutes()
	RegisterLoginRoutes()

	// Expense routes
	RegisterExpenseRoutes()
}
