package ctrls

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
)

type Home struct {
	baseController
}

func (this *Home) Get() {
	this.TplNames = "index.tpl"

	this.Data["Lang"] = this.Lang
	this.Data["Title"] = "Yalp.uz"
	this.Data["Name"] = "My name is Mirolim!"
	isAuth := this.GetSession("isAuth")
	if isAuth != nil {
		this.Data["isAuth"] = isAuth.(bool)
	}
}

func (this *Home) Category() {

	var err error
	// get attr, tag and city
	city := this.Ctx.Input.Param(":city")
	bizType := this.Ctx.Input.Param(":bizType")

	var catList []List
	foodService := new(models.FoodService)
	switch bizType {
	case foodService.GetC():
		beego.Warn("foodservice type")
		err = models.DocCountDistinct(foodService, bson.M{"lang": this.Lang, "city": city},
			"types",
			&catList,
			60)
	default:
		this.Abort("404")
		return
	}
	check("FS->category DocCountDistinct -> ", err)
	this.TplNames = "food-service/category.tpl"

	this.Data["Data"] = struct {
		Category string
		City     string
		CatList  []List
	}{
		foodService.GetC(),
		city,
		catList,
	}
}
