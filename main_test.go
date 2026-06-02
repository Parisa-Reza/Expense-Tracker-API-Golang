package main

import (
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestMainPackage(t *testing.T) {
	originalRunServer := runServer
	runServer = func(params ...string) {}
	defer func() {
		runServer = originalRunServer
	}()

	tests := []struct {
		name    string
		runMode string
	}{
		{name: "dev mode config", runMode: "dev"},
		{name: "prod mode config", runMode: "prod"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beego.BConfig.RunMode = tt.runMode
			main()
			if tt.runMode == "dev" && !beego.BConfig.WebConfig.DirectoryIndex {
				t.Fatalf("DirectoryIndex = false, want true")
			}
		})
	}
}
