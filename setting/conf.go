package setting

import (
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	APP_VER = "0.1"
)

var (
	AppName string
	AppVer  string
	AppHost string
	AppUrl  string
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

	// TODO captcha

	// TODO database
	driverName := Cfg.MustValue("orm", "driver_name", "mysql")
	dataSource := Cfg.MustValue("orm", "data_source", "root:root@/wetalk?charset=utf8&loc=UTC")
	maxIdle := Cfg.MustInt("orm", "max_idle_conn", 30)
	maxOpen := Cfg.MustInt("orm", "max_open_conn", 50)

	// set default database
	err = orm.RegisterDataBase("default", driverName, dataSource, maxIdle, maxOpen)
	if err != nil {
		beego.Error(err)
	}

	// TODO ..

}

func reloadConfig() {
	AppHost = Cfg.MustValue("app", "app_host", "127.0.0.1:8092")
	AppUrl = Cfg.MustValue("app", "app_url", "http://127.0.0.1:8092/")
}
