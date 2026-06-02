package main

import (
	_ "expense-tracker-api/routers"

	beego "github.com/beego/beego/v2/server/web"
)

var runServer = beego.Run

func main() {
	configureServer()
	runServer()
}

func configureServer() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
}
