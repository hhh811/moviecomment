package main

import (
	"fmt"
	"moviecomment/controllers/auth"
	_ "moviecomment/routers"
	"moviecomment/setting"

	"github.com/astaxie/beego"
)

func initialize() {
	setting.LoadConfig()

	// TODO sphinx

}

func main() {
	beego.SetLogFuncCall(true)

	initialize()

	beego.Info("AppPath:", beego.AppPath)
	// TODO 版本信息，url信息等 log

	// TODO 静态目录设置

	// TODO filters

	// 注册 routers
	login := new(auth.LoginRouter)
	beego.Router("/login", login, "get:Get;post:Login")
	beego.Router("/logout", login, "get:Logout")

	// beego.Run()

}

func demo() (success bool) {
	if false {
		return false
	}

	defer func() {
		if !success {
			fmt.Print("not")
		}
	}()

	return true
}
