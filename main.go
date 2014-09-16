package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"github.com/mmirolim/yalp-go/ctrls"
	"github.com/mmirolim/yalp-go/models"
	_ "github.com/mmirolim/yalp-go/routers"
	_ "github.com/astaxie/beego/session/redis"
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
	}
	// change default tpl tags
	beego.TemplateLeft = "[["
	beego.TemplateRight = "]]"

	// register a i18n template func
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.AddFuncMap("getUrl", ctrls.GetUrl)

	beego.Run()

}
