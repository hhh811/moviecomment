package setting

import (
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego"
)

const (
	APP_VER = "0.1"
)

var (
	AppName string
	AppVer  string
	AppHost string
	AppUrl  string

	LoginRememberDays int // 登录记录天数

	CookieRememberName string
	CookieUserName     string
)

var (
	Cfg *goconfig.ConfigFile
)

var (
	GlobalConfPath = "conf/_"
	AppConfPath    = "conf/app.conf"
)

// LoadConfig loads configuration file.
func LoadConfig() *goconfig.ConfigFile {
	var err error

	if fd, _ := os.OpenFile(AppConfPath, os.O_RDONLY|os.O_CREATE, 0600); fd != nil {
		fd.Close()
	}

	// app version
	AppVer = strings.Split(APP_VER, ".")[0]

	// load additional config besides beego.BConfig
	Cfg, err = goconfig.LoadConfigFile(AppConfPath)
	if err != nil {
		fmt.Println("Fail to load configuration file: " + err.Error())
		os.Exit(2)
	}

	// session settings
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = Cfg.MustValue("session", "session_name", "moviecomment_sess")
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = Cfg.MustInt("session", "session_life_time", 0)
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = Cfg.MustInt64("session", "session_gc_time", 86400)

	// TODO captcha

	// TODO database
	// driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	// dataSource := Cfg.MustValue("orm", "data_source", "root:root@/wetalk?charset=utf8&loc=UTC")
	// maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	// maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	// set default database
	// err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	// if err != nil {
	// 	beego.Error(err)
	// }

	// TODO ..

	reloadConfig()

	return Cfg
}

func reloadConfig() {
	AppHost = Cfg.MustValue("app", "app_host", "127.0.0.1:8088")
	AppUrl = Cfg.MustValue("app", "app_url", "http://127.0.0.1:8088/")

	LoginRememberDays = Cfg.MustInt("app", "login_remember_days", 7)

	CookieRememberName = Cfg.MustValue("app", "cookie_remember_name", "moviecomment_cookie_remember_name")
	CookieUserName = Cfg.MustValue("app", "cooki_user_name", "moviecomment_cookie_user_name")
}
