package main

import (
	_ "github.com/mmirolim/yalp-go/routers"
	"github.com/mmirolim/yalp-go/controllers"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

const (
	APP_VER = "0.0.1"
)
func initialize() {
	controllers.InitApp()
}
func main() {
	initialize()

	if beego.AppConfig.String("runmode") == "dev" {
		beego.EnableAdmin = true
		beego.AdminHttpAddr = "192.168.1.107"
		beego.AdminHttpPort = 8088
	}
	// change default tpl tags
	beego.TemplateLeft = "[["
	beego.TemplateRight = "]]"


	// register a i18n template func
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.Run()
}

