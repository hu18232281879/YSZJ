package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"pyg/controllers"
)

func init() {
	beego.InsertFilter("/user/*",beego.BeforeExec,filterFunc)
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/sendMsg", &controllers.UserController{}, "post:SendMsg")
	beego.Router("/active", &controllers.UserController{}, "get:ShowActive;post:HandelActive")
	beego.Router("/activeUser", &controllers.UserController{}, "get:ActivateTheSuccess")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/",&controllers.GoodsController{},"get:ShowIndex")
	beego.Router("/logOut",&controllers.UserController{},"get:LogOut")
	beego.Router("/user/userCenterInfo",&controllers.UserController{},"get:ShowUserCenterInfo")
	beego.Router("/user/userCenterSite",&controllers.UserController{},"get:ShowUserCenterSite")
	beego.Router("/user/submitAddress",&controllers.UserController{},"post:AddAddress")
	beego.Router("/indexSx",&controllers.GoodsController{},"get:ShowIndexSx")

}

func filterFunc(ctx *context.Context){
	userName:=ctx.Input.Session("pyg_userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
		return
	}

}
