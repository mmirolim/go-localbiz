package routers

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/ctrl"
)

func init() {
	beego.Router("/:city/"+ctrl.FoodService{}.Slug()+"/:attr?/:tag", &ctrl.FoodService{}, "get:Category")
	beego.Router("/"+ctrl.FoodService{}.Slug()+"/:slug", &ctrl.FoodService{})
	beego.Router("/login", &ctrl.Auth{})
	beego.Router("/signup", &ctrl.User{}, "*:SignUp")
	beego.Router("/", &ctrl.Home{})

}
