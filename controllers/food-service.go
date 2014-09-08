package controllers

import (
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/astaxie/beego"
)

type FoodServiceCtrl struct {
	baseController
}

func (this *FoodServiceCtrl) Get() {
	this.Data["Lang"] = this.Lang

	// get FoodService by slug
	slug := this.Ctx.Input.Param(":slug")

	foodService, err := models.FoodServices.FindOne(bson.M{ "slug" : slug })
	if err != nil {
		beego.Error(err)
		this.Abort("404")
	}
	nearResult, err := models.FoodServices.FindNear(1, 1000, foodService.GeoJson)
	if err == nil {
		this.Data["NearResult"] = nearResult
	} else {
		beego.Warn(err)
	}

	this.Data["Title"] = "Title - District - City | APPNAME"
	this.TplNames = "food-service/food-service.tpl"
	this.Data["Entity"] = foodService

}
