package controllers

import (

)

type HomeController struct {
	baseController
}

func (this *HomeController) Get() {
	this.Data["Lang"] = this.Lang
	this.Data["Title"] = "Yalp.uz"
	this.Data["Name"] = "My name is Mirolim!"
	this.Layout = "layout.tpl"
	this.TplNames = "index.tpl"
    this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"

}
