package controllers

import (
)

type RestaurantCtrl struct {
	baseController
}

func (this *RestaurantCtrl) Get() {
	this.Data["Lang"] = this.Lang
	this.Data["Title"] = "Restaurant Title"
	this.Data["Desc"] = "Restaurant info"
	this.TplNames = "index.tpl"
}
