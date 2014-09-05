package main

import (
	_ "github.com/mmirolim/yalp-go/routers"
	"github.com/astaxie/beego"
)

func main() {
	if beego.AppConfig.String("runmode") == "dev" {
		beego.EnableAdmin = true
		beego.AdminHttpAddr = "192.168.1.107"
		beego.AdminHttpPort = 8088
	}
	// change default tpl tags
	beego.TemplateLeft = "[["
	beego.TemplateRight = "]]"
	beego.Run()
}

