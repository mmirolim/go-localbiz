package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/mmirolim/yalp-go/ctrls"
)

func IsAuthenticated(ctx *context.Context) {
	beego.Warn(ctx.Input.Session("number"))
	_, ok := ctx.Input.Session("uid").(int)
	if !ok {
		ctx.Redirect(302, "/login")
	}
}

func init() {
	beego.Router("/p/:city/:bizType", &ctrls.Home{}, "get:Category")
	beego.Router("/"+ctrls.FoodService{}.Slug()+"/:slug", &ctrls.FoodService{})
	beego.Router("/login/?:socialNet", &ctrls.Auth{}, "get:Login")
	beego.Router("/logout", &ctrls.Auth{}, "get:Logout")
	beego.Router("/auth/?:socialNet", &ctrls.Auth{}, "*:Authorize")
	beego.Router("/user/:id", &ctrls.User{})
	beego.Router("/signup", &ctrls.User{}, "get:SignUp;post:SignUpProcess")
	beego.Router("/", &ctrls.Home{})

	beego.InsertFilter("/*/edit", beego.BeforeRouter, IsAuthenticated)

}
