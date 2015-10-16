package routers

import (
	"admin/controllers"
	"admin/src/admin"
	"github.com/astaxie/beego"
)

func init() {
	admin.Run()
	beego.Router("/", &controllers.MainController{})
}
