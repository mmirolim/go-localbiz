package routers

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/controllers"
)

func init() {
	beego.Router("/:city/" + controllers.FoodServiceCtrl{}.Slug() +"/:attr?/:tag", &controllers.FoodServiceCtrl{}, "*:Category")
	beego.Router("/" + controllers.FoodServiceCtrl{}.Slug() +"/:slug", &controllers.FoodServiceCtrl{})
	beego.Router("/", &controllers.HomeCtrl{})

}
