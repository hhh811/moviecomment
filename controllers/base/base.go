package base

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"moviecomment/models"
	"moviecomment/modules/auth"
	"moviecomment/setting"
	"moviecomment/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
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
	} else if auth.LoginUserFromRememberCookie(&rt.User, rt.Ctx) {
		rt.IsLogin = true
	}

	if rt.IsLogin {
		rt.IsLogin = true
		rt.Data["user"] = &rt.User
		rt.Data["IsLogin"] = rt.IsLogin

		// 用户被禁用的处理
	}

	// 配置数据
	rt.Data["AppName"] = setting.AppName
	rt.Data["AppVer"] = setting.AppVer
	rt.Data["AppUrl"] = setting.AppUrl

	// TODO set lang and redirect to make url clean

	// read flash message
	beego.ReadFromRequest(&rt.Controller)

	// TODO xsrf

	// TODO form once

}

// on router finished
func (rt *BaseRouter) Finish() {

}

// check if not login when redirect
func (rt *BaseRouter) CheckLoginRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			// custom redirect url
			code = v
		}
	}

	// if need login then redirect
	if needLogin && !rt.IsLogin {
		if len(redirect_to) == 0 {
			req := rt.Ctx.Request
			scheme := "http"
			if req.TLS != nil {
				scheme += "s"
			}
			redirect_to = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
		}
		redirect_to = "/login?to=" + url.QueryEscape(redirect_to)
		rt.Redirect(redirect_to, code)
		return true
	}

	// if not need login then redirect
	if !needLogin && rt.IsLogin {
		if len(redirect_to) == 0 {
			redirect_to = "/"
		}
		rt.Redirect(redirect_to, code)
		return true
	}
	return false
}

func (rt *BaseRouter) SetFormSets(form interface{}, names ...string) *utils.FormSets {
	return rt.setFormSets(form, nil, names...)
}

func (rt *BaseRouter) setFormSets(form interface{}, errs map[string]*validation.Error, names ...string) *utils.FormSets {
	formSets := utils.NewFormSets(form, errs)
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	name += "Sets"
	rt.Data[name] = formSets

	return formSets
}
