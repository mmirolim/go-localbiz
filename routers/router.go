package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/mmirolim/yalp-go/ctrls"
	M "github.com/mmirolim/yalp-go/models"
)

// @todo most of aspects should be in filters, ctrls will be shorter
func init() {
	beego.Router("/p/:city/:bizType", &ctrls.Home{}, "get:Category")
	beego.Router("/"+ctrls.FoodService{}.Slug()+"/:slug", &ctrls.FoodService{})
	beego.Router("/login/?:provider", &ctrls.Auth{}, "get:Login")
	beego.Router("/logout", &ctrls.Auth{}, "get:Logout")
	beego.Router("/auth/?:provider", &ctrls.Auth{}, "*:Authorize")
	beego.Router("/user/:username", &ctrls.User{})
	beego.Router("/user/edit", &ctrls.User{}, "get:Edit;post:EditProc")
	beego.Router("/signup", &ctrls.User{}, "get:SignUp;post:SignUpProc")
	beego.Router("/", &ctrls.Home{})

	beego.InsertFilter("/*/edit", beego.BeforeRouter, isAuth)

	// admin namespace
	adm := beego.NewNamespace("/backend",
		beego.NSBefore(allowBackend),
		beego.NSRouter("/dashboard", &ctrls.Backend{}, "get:DashBoard"),
	)
	beego.AddNamespace(adm)

}

func isAuth(c *context.Context) {
	uid := c.Input.Session("uid")
	beego.Warn("Filter isAuth")
	if uid == nil {
		c.Redirect(302, "/login")
	}
}

// filter out simple users
func allowBackend(c *context.Context) {
	uid := c.Input.Session("uid")
	if uid == nil {
		c.Abort(403, "Sorry it is a protected area")
	}
	u := M.User{}
	if !u.AllowBackend(uid.(string)) {
		c.Abort(403, "Sorry it is a private party")
	}
}
