package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/mmirolim/yalp-go/ctrls"
	"github.com/mmirolim/yalp-go/models"
	_ "github.com/mmirolim/yalp-go/routers"
	"github.com/nicksnyder/go-i18n/i18n"
)

func initialize() {
	// set constants for ctrl
	ctrls.AppVer = beego.AppConfig.String("appver")
	ctrls.APP = beego.AppConfig.String("appname")
	ctrls.IsPro = beego.RunMode == "prod"

	// init ctrl
	ctrls.Initialize()
	// initialize models
	models.Initialize()

}

func main() {
	// @todo refactor with aspect oriented programming
	initialize()

	if beego.AppConfig.String("runmode") == "dev" {
		beego.EnableAdmin = true
		beego.AdminHttpAddr = "192.168.1.107"
		beego.AdminHttpPort = 8088
		ctrls.IsPro = false
	}

	// enable CSRF protection
	// @todo handle properly csrf token expire
	beego.EnableXSRF = true
	// @todo change csrf key
	beego.XSRFKEY = beego.AppConfig.String("xsrfKey")
	beego.XSRFExpire = 3600

	// change default tpl tags
	beego.TemplateLeft = "[["
	beego.TemplateRight = "]]"

	// register i18n T func
	beego.AddFuncMap("T", i18n.IdentityTfunc)

	beego.Run()

}
