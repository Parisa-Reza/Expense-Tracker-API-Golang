package routers

import "testing"

func TestRegisterHealthRoutes(t *testing.T) {
	tests := []struct {
		name string
		run  func()
	}{
		{name: "register health routes", run: RegisterHealthRoutes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
		})
	}
}
