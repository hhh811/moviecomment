package base

import (
	"time"

	"moviecomment/models"
	"moviecomment/modules/auth"

	"github.com/astaxie/beego"
)

// baseRouter implemented global settings for all other routers.
type BaseRouter struct {
	beego.Controller
	User    models.User
	IsLogin bool
}

// Prepare implemented Prepare method for baseRouter
func (rt *BaseRouter) Prepare() {
	//TODO enforce redirect

	// page start time
	rt.Data["PageStartTime"] = time.Now()

	// start session
	rt.StartSession()

	// TODO flash 数据处理

	// 是否已经登录
	if auth.GetUserFromSession(&rt.User, rt.CruSession) {
		rt.IsLogin = true
	} else if auth.
}
