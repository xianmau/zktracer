package routers

import (
	"github.com/astaxie/beego"
	"zktracer/controllers"
)

func init() {

	beego.Router("/", &controllers.HomeController{})
	beego.Router("/home", &controllers.HomeController{})

	beego.Router("/broker", &controllers.BrokerController{})
	beego.Router("/broker/detail", &controllers.BrokerController{}, "get:Detail")
	beego.Router("/broker/getlatestdata", &controllers.BrokerController{}, "get:GetLatestData")

	beego.Router("/logger", &controllers.LoggerController{})
	beego.Router("/logger/detail", &controllers.LoggerController{}, "get:Detail")
	beego.Router("/logger/getlatestdata", &controllers.LoggerController{}, "get:GetLatestData")

	beego.Router("/app", &controllers.AppController{})

	beego.Router("/topic", &controllers.TopicController{})

	beego.Router("/login", &controllers.LoginController{})

	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/admin/node", &controllers.AdminController{})
	beego.Router("/admin/zone", &controllers.AdminController{}, "get:Zone")
	beego.Router("/admin/zone/create", &controllers.AdminController{}, "post:CreateZone")
	beego.Router("/admin/zone/delete", &controllers.AdminController{}, "post:DeleteZone")
	beego.Router("/admin/getdata", &controllers.AdminController{}, "get:GetData")
	beego.Router("/admin/getnodedata", &controllers.AdminController{}, "get:GetNodeData")
	beego.Router("/admin/createnode", &controllers.AdminController{}, "post:CreateNode")
	beego.Router("/admin/deletenode", &controllers.AdminController{}, "post:DeleteNode")
}
