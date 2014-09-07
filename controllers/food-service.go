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
	//@todo implement

	this.Data["Title"] = "Title - District - City | APPNAME"
	this.Data["Desc"] = "Restaurant info"
	this.TplNames = "food-service/food-service.tpl"

	foodService, err := models.FoodServices.FindOne(bson.M{ "name" : "Gary Danko" })
	if err != nil {
		beego.Error(err)
		return
	}

	this.Data["Title"] = foodService.Name
	this.Data["Desc"] = foodService.Description

}
