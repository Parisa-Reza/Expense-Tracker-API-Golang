package routers

import "testing"

func TestRegisterRegisterRoutes(t *testing.T) {
	tests := []struct {
		name string
		run  func()
	}{
		{name: "register register routes", run: RegisterRegisterRoutes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
		})
	}
}
