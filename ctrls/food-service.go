package ctrls

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
)

type List struct {
	Id    string `bson:"_id"`
	Count uint16 `bson:"count"`
}

type FoodService struct {
	baseController
}

func (c FoodService) Slug() string {
	return "fs"
}

func (c *FoodService) Get() {
	c.Data["Lang"] = c.Lang
	// get FoodService by slug
	slug := c.Ctx.Input.Param(":slug")
	var fd models.FoodService
	err := models.DocFindOne(bson.M{"slug": slug}, bson.M{}, &fd, 60)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}
	var near models.Near
	err = models.DocFindNear(1, 1000, &fd, &near, 60)
	check("FS Get -> ", err)
	var fds []struct {
		Dis float32
		Obj models.FoodService
	}
	err = near.Results.Unmarshal(&fds)
	check("FS Get raw unmarshal -> ", err)

	c.Data["Near"] = fds
	c.Data["Title"] = "Title - District - City | APPNAME"
	c.TplNames = "food-service/food-service.tpl"
	c.Data["Entity"] = fd
	c.Data["CtrlSlug"] = c.Slug()

}

// method to process fs category requests
func (c *FoodService) Category() {
	var err error
	beego.Error("Cat")
	// get attr, tag and city
	attr := c.Ctx.Input.Param(":attr")
	tag := strings.Replace(c.Ctx.Input.Param(":tag"), "_", " ", -1)
	city := c.Ctx.Input.Param(":city")
	beego.Warn(tag)
	// $regex query used to match case sensitive index
	q := bson.D{
		{"lang", c.Lang},
		{attr, bson.M{"$regex": bson.RegEx{Pattern: `^` + tag, Options: "i"}}},
		{"city", city},
	}

	// cache category list
	var fds []models.FoodService
	// get all places with cat and city
	err = models.DocFind(q, bson.M{"name": 1, "slug": 1}, &models.FoodService{}, &fds, 60)
	check("Category FInd ->", err)

	count := len(fds)

	var catList []List
	err = models.DocCountDistinct(&models.FoodService{}, bson.M{"lang": c.Lang}, "types", &catList, 60)
	check("FS->category DocCountDistinct -> ", err)
	c.TplNames = "food-service/category.tpl"

	c.Data["Data"] = struct {
		Category string
		City     string
		FdsList  []models.FoodService
		Count    int
		CatList  []List
	}{
		tag,
		city,
		fds,
		count,
		catList,
	}

}
