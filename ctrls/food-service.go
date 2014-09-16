package ctrls

import (
	s "strings"

	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
)

type FoodService struct {
	baseController
}

func (this FoodService) Slug() string {
	return "fs"
}

func (this *FoodService) Get() {
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
func (this *FoodService) Category() {
	// get attr, tag and city
	attr := this.Ctx.Input.Param(":attr")
	tag := s.Replace(this.Ctx.Input.Param(":tag"), "_", " ", -1)
	city := this.Ctx.Input.Param(":city")

	// $regex query used to match case sensitive index
	q := bson.D{
		{"lang", this.Lang},
		{attr, bson.M{"$regex": bson.RegEx{`^` + tag, "i"}}},
		{"address.city", bson.M{"$regex": bson.RegEx{`^` + city, "i"}}},
	}
	// get all places with cat and city
	fds, err := models.FoodServices.Find(q)
	check("FSCtrl.Category -> ", err)
	count := len(fds)
	this.TplNames = "food-service/category.tpl"

	this.Data["Data"] = struct {
		Category string
		City     string
		FdsList  []models.FoodService
		Count    int
	}{
		tag,
		city,
		fds,
		count,
	}

}
