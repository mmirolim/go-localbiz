package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["Title"] = "Yalp.uz"
	this.Data["Name"] = "My name is Mirolim!"
	this.Layout = "layout.tpl"
	this.TplNames = "index.tpl"
    this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
}
