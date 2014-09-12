package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"github.com/mmirolim/yalp-go/ctrl"
	"github.com/mmirolim/yalp-go/models"
	_ "github.com/mmirolim/yalp-go/routers"
)

func initialize() {
	// set constants for ctrl
	ctrl.AppVer = beego.AppConfig.String("appver")
	ctrl.APP = beego.AppConfig.String("appname")
	ctrl.IsPro = beego.RunMode == "prod"

	// init ctrl
	ctrl.InitApp()
	// init connection to mongodb
	models.InitConnection()
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
