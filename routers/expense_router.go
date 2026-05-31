// @Title Expense Routes
// @Description Handles CRUD endpoints for authenticated user expenses.
// @Summary Registers expense endpoints
// @Tags Expenses
// @Version 1.0.0
package routers

import (
	"expense-tracker-api/controllers"

	"github.com/beego/beego/v2/server/web"
)

// RegisterExpenseRoutes registers all expense-related routes.
func RegisterExpenseRoutes() {
	web.Router("/api/v1/expenses", &controllers.ExpenseCreateController{}, "post:Create")
	web.Router("/api/v1/expenses", &controllers.ExpenseListController{}, "get:List")
	web.Router("/api/v1/expenses/:id", &controllers.ExpenseGetController{}, "get:Get")
	web.Router("/api/v1/expenses/:id", &controllers.ExpenseUpdateController{}, "put:Update")
	web.Router("/api/v1/expenses/:id", &controllers.ExpenseDeleteController{}, "delete:Delete")
}
