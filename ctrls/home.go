package ctrls

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
)

type Home struct {
	baseController
}

func (c *Home) Get() {
	c.TplNames = "index.tpl"

	c.Data["csrf"] = c.XsrfToken()
	c.Data["Lang"] = c.Lang
	c.Data["Title"] = "Yalp.uz"
	c.Data["Name"] = "My name is Mirolim!"
	c.Data["Person"] = map[string]interface{}{"Person": "Mirolim"}
	isAuth := c.GetSession("isAuth")
	if isAuth != nil {
		c.Data["isAuth"] = isAuth.(bool)
	}
}

func (c *Home) Category() {

	var err error
	// get attr, tag and city
	city := c.Ctx.Input.Param(":city")
	bizType := c.Ctx.Input.Param(":bizType")

	var catList []List
	foodService := new(models.FoodService)
	switch bizType {
	case foodService.GetC():
		beego.Warn("foodservice type")
		err = models.DocCountDistinct(foodService, bson.M{"lang": c.Lang, "city": city},
			"types",
			&catList,
			60)
	default:
		c.Abort("404")
		return
	}
	check("FS->category DocCountDistinct -> ", err)
	c.TplNames = "food-service/category.tpl"

	c.Data["Data"] = struct {
		Category string
		City     string
		CatList  []List
	}{
		foodService.GetC(),
		city,
		catList,
	}
}
