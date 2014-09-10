package controllers

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
)

type FoodServiceCtrl struct {
	baseController
}

func (this FoodServiceCtrl) Slug() string {
	return "fs"
}

func (this *FoodServiceCtrl) Get() {
	this.Data["Lang"] = this.Lang

	// get FoodService by slug
	slug := this.Ctx.Input.Param(":slug")

	foodService, err := models.FoodServices.FindOne(bson.M{"slug": slug})
	if err != nil {
		beego.Error(err)
		this.Abort("404")
	}
	near, err := models.FoodServices.FindNear(1, 1000, foodService.GeoJson)
	if err == nil {
		this.Data["Near"] = near
	} else {
		beego.Warn(err)
	}

	this.Data["Title"] = "Title - District - City | APPNAME"
	this.TplNames = "food-service/food-service.tpl"
	this.Data["Entity"] = foodService
	this.Data["CtrlSlug"] = this.Slug()

}
