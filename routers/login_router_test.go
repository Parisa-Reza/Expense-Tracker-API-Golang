package routers

import "testing"

func TestRegisterLoginRoutes(t *testing.T) {
	tests := []struct {
		name string
		run  func()
	}{
		{name: "register login routes", run: RegisterLoginRoutes},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
		})
	}
}
