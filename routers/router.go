package routers

import (
	"github.com/astaxie/beego"
	"pyg/controllers"
)

func init() {
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/sendMsg", &controllers.UserController{}, "post:SendMsg")
	beego.Router("/active", &controllers.UserController{}, "get:ShowActive;post:HandelActive")
	beego.Router("/activeUser", &controllers.UserController{}, "get:ActivateTheSuccess")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/",&controllers.GoodsController{},"get:ShowIndex")
}
