package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/mmirolim/yalp-go/ctrls"
)

func IsAuth(ctx *context.Context) {
	uid := ctx.Input.Session("uid")
	beego.Warn("Filter isAuth")
	if uid == nil {
		ctx.Redirect(302, "/login")
	}
}

// most of aspects should be in filters, ctrls will be shorter
func init() {
	beego.Router("/p/:city/:bizType", &ctrls.Home{}, "get:Category")
	beego.Router("/"+ctrls.FoodService{}.Slug()+"/:slug", &ctrls.FoodService{})
	beego.Router("/login/?:socialNet", &ctrls.Auth{}, "get:Login")
	beego.Router("/logout", &ctrls.Auth{}, "get:Logout")
	beego.Router("/auth/?:socialNet", &ctrls.Auth{}, "*:Authorize")
	beego.Router("/user/:username", &ctrls.User{})
	beego.Router("/user/edit/:id", &ctrls.User{}, "*:Edit")
	beego.Router("/signup", &ctrls.User{}, "get:SignUp;post:SignUpProc")
	beego.Router("/", &ctrls.Home{})

	beego.InsertFilter("/*/edit/*", beego.BeforeRouter, IsAuth)

}
