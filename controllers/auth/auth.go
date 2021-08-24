package auth

import (
	"moviecomment/controllers/base"
	"moviecomment/modules/auth"
	"moviecomment/utils"
	"strings"
)

// LoginRouter servers login page
type LoginRouter struct {
	base.BaseRouter
}

// Get implemented login page
func (rt *LoginRouter) Get() {
	rt.Data["IsLoginPage"] = true
	rt.TplName = "auth/login.html"

	// url = "/login?to=..." loginRedirect 获取 to= 后面的字符串
	loginRedirect := strings.TrimSpace(rt.GetString("to"))
	if !utils.IsMatchHost(loginRedirect) {
		loginRedirect = "/"
	}

	// no need login
	// CheckLoginRedirect 函数逻辑就是如果已经登录，重定向到 /
	if rt.CheckLoginRedirect(false, loginRedirect) {
		return
	}

	if len(loginRedirect) > 0 {
		rt.Ctx.SetCookie("login_to", loginRedirect, 0, "/")
	}

	form := auth.LoginForm{}
	rt.SetFormSets(&form)

}
