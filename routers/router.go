package routers

import (
	"github.com/mmirolim/yalp-go/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeCtrl{})
	beego.Router("/fs/:biz", &controllers.FoodServiceCtrl{})
}
