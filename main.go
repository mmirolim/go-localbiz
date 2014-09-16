package main

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/ctrls"
	"github.com/mmirolim/yalp-go/models"
	_ "github.com/mmirolim/yalp-go/routers"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/nicksnyder/go-i18n/i18n"
)

func initialize() {
	// set constants for ctrl
	ctrls.AppVer = beego.AppConfig.String("appver")
	ctrls.APP = beego.AppConfig.String("appname")
	ctrls.IsPro = beego.RunMode == "prod"

	// init ctrl
	ctrls.InitApp()
	// init connection to mongodb
	models.InitConnection()
}

func main() {
	initialize()

	if beego.AppConfig.String("runmode") == "dev" {
		beego.EnableAdmin = true
		beego.AdminHttpAddr = "192.168.1.107"
		beego.AdminHttpPort = 8088
		ctrls.IsPro = false
	}
	// change default tpl tags
	beego.TemplateLeft = "[["
	beego.TemplateRight = "]]"

	// register i18n T func
	beego.AddFuncMap("T", i18n.IdentityTfunc)

	beego.Run()

}
