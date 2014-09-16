package routers

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/ctrls"
)

func init() {
	beego.Router("/:city/"+ctrls.FoodService{}.Slug()+"/:attr?/:tag", &ctrls.FoodService{}, "get:Category")
	beego.Router("/"+ctrls.FoodService{}.Slug()+"/:slug", &ctrls.FoodService{})
	beego.Router("/login", &ctrls.Auth{})
	beego.Router("/signup", &ctrls.User{}, "*:SignUp")
	beego.Router("/", &ctrls.Home{})

}
