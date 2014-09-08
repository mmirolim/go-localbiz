package routers

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeCtrl{})
	beego.Router("/fs/:slug", &controllers.FoodServiceCtrl{})
}
