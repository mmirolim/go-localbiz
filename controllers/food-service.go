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

// method to process fs category requests
func (this *FoodServiceCtrl) Category() {
	// get attr, tag and city
	attr := this.Ctx.Input.Param(":attr")
	tag := this.Ctx.Input.Param(":tag")
	city := this.Ctx.Input.Param(":city")

	// get all places with cat and city
	fds, err := models.FoodServices.Find(bson.D{ { "lang", this.Lang }, {attr , tag} , { "address.city", city} })
	check("FSCtrl.Category -> ", err)
	count := len(fds)
	this.TplNames = "food-service/category.tpl"

	this.Data["Data"] = struct {
		Category string
		City	string
		FdsList []models.FoodService
		Count int
	}{
		tag,
		city,
		fds,
		count,
	}


}
