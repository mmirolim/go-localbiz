package controllers

import (
	"github.com/astaxie/beego"
)

type RestaurantCtrl struct {
	beego.Controller
}

func (this *RestaurantCtrl) Get() {
	this.Data["Title"] = "Restaurant Title"
	this.Data["Desc"] = "Restaurant info"
	this.TplNames = "index.tpl"
}
