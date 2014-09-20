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
	var fd models.FoodService
	err := models.DocFindOne(bson.M{"slug": slug}, &fd)
	if err != nil {
		beego.Error(err)
		this.Abort("404")
	}
	var near models.Near
	err = models.DocFindNear(1, 1000, fd, &near)
	check("FS Get -> ", err)
	var fds []struct {
		Dis float32
		Obj models.FoodService
	}
	err = near.Results.Unmarshal(&fds)
	check("FS Get raw unmarshal -> ", err)

	this.Data["Near"] = fds
	this.Data["Title"] = "Title - District - City | APPNAME"
	this.TplNames = "food-service/food-service.tpl"
	this.Data["Entity"] = fd
	this.Data["CtrlSlug"] = this.Slug()

}

// method to process fs category requests
func (this *FoodService) Category() {
	var err error

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

	// cache category list
	cacheKey := city + attr + tag
	var fds []models.FoodService
	if !cacheEnabled || !cacheIsExist(cacheKey) {
		// get all places with cat and city
		err = models.DocFind(q, models.FoodService{}, &fds)
		check("Category FInd ->", err)
		if cacheEnabled && err == nil {
			cachePut(cacheKey, fds, 60)
		}
	} else {
		cacheGet(cacheKey, &fds)
	}
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
