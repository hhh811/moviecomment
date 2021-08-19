package auth

import (
	"moviecomment/models"
	"moviecomment/setting"
	"moviecomment/utils"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/session"
)

// get login redirect url from cookie
func GetLoginRedirect(ctx *context.Context) string {
	loginRedirect := strings.TrimSpace(ctx.GetCookie("login_to"))
	if !utils.IsMatchHost(loginRedirect) {
		loginRedirect = "/"
	} else {
		ctx.SetCookie("login_to", "", -1, "/")
	}
	return loginRedirect
}

// login user
func LoginUser(user *models.User, ctx *context.Context, remember bool) {
	// beego session regenerate id
	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
	ctx.Input.CruSession = beego.GlobalSessions.SessionRegenerateID(ctx.ResponseWriter, ctx.Request)
	ctx.Input.CruSession.Set("auth_user_id", user.Id)

	if remember {
		// 重新写 cookie 刷新记录时间
		WriteRememberCookie(user, ctx)
	}
}

func WriteRememberCookie(user *models.User, ctx *context.Context) {
	secret := utils.EncodeMd5(user.Rands + user.Password)
	days := 86400 * setting.LoginRememberDays
	ctx.SetCookie(setting.CookieUserName, user.UserName, days)
	ctx.SetSecureCookie(secret, setting.CookieRememberName, user.UserName, days)
}

func DeleteRememberCookie(ctx *context.Context) {
	ctx.SetCookie(setting.CookieUserName, "", -1)
	ctx.SetCookie(setting.CookieRememberName, "", -1)
}

// 如果浏览器记录了 cookie 直接登录
func LoginUserFromRememberCookie(user *models.User, ctx *context.Context) (success bool) {
	userName := ctx.GetCookie(setting.CookieUserName)
	if len(userName) == 0 {
		return false
	}

	// 后面的操作如果返回 false
	defer func() {
		if !success {
			DeleteRememberCookie(ctx)
		}
	}()

	user.UserName = userName
	// 没有该用户
	if err := user.Read("UserName"); err != nil {
		return false
	}

	// 检查 cookie 记录的密码
	secret := utils.EncodeMd5(user.Rands + user.Password)
	value, _ := ctx.GetSecureCookie(secret, setting.CookieRememberName)
	if value != userName {
		return false
	}

	LoginUser(user, ctx, true)

	return true
}

func GetUserIdFromSession(sess session.Store) int {
	if id, ok := sess.Get("auth_user_id").(int); ok && id > 0 {
		return id
	}
	return 0
}

// get user if key exists in session
func GetUserFromSession(user *models.User, sess session.Store) bool {
	id := GetUserIdFromSession(sess)
	if id > 0 {
		u := models.User{Id: id}
		if u.Read() == nil {
			*user = u
			return true
		}
	}

	return false
}
