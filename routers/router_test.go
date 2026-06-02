package routers

import "testing"

func TestRouterPackageInit(t *testing.T) {
	tests := []struct {
		name string
		run  func()
	}{
		{name: "register all routes", run: func() {
			RegisterHealthRoutes()
			RegisterRegisterRoutes()
			RegisterLoginRoutes()
			RegisterExpenseRoutes()
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
		})
	}
}
