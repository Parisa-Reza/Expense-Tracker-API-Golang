package routers

import "testing"

func TestRegisterExpenseRoutes(t *testing.T) {
	tests := []struct {
		name string
		run  func()
	}{
		{name: "register expense routes", run: RegisterExpenseRoutes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
		})
	}
}
