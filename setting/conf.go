package setting

import (
	"fmt"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
)

const (
	APP_VER = "0.1"
)

var (
	AppName string
	AppVer  string
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

	// load config file
	Cfg, err = goconfig.LoadConfigFile(AppConfPath)
	if err != nil {
		fmt.Println("Fail to load configuration file: " + err.Error())
		os.Exit(2)
	} else {
		Cfg.AppendFiles(AppConfPath)
	}

	// app version
	AppVer = strings.Split(APP_VER, ".")[0]

	// TODO captcha
}
