package routers

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeCtrl{})
	beego.Router("/" + controllers.FoodServiceCtrl{}.Slug() +"/:slug", &controllers.FoodServiceCtrl{})
}
