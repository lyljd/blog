package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

type TomlConfig struct {
	Viewer Viewer
	System SystemConfig
}

type Viewer struct {
	Title       string
	Description string
	Navigation  []string
	Bilibili    string
	Github      string
	Avatar      string
	UserName    string
	UserDesc    string
}

type SystemConfig struct {
	AppName      string
	Version      float32
	CurrentDir   string
	Username     string
	Password     string
	Nickname     string
	Valine       bool
	ValineAppid  string
	ValineAppkey string
}

var Cfg *TomlConfig

func init() {
	Cfg = new(TomlConfig)
	var err error
	Cfg.System.CurrentDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	Cfg.System.AppName = "blog"
	Cfg.System.Version = 1.0
	_, err = toml.DecodeFile("config/config.toml", &Cfg)
	if err != nil {
		panic("读取配置失败，" + err.Error())
	}
}

func CheckSlugExist(slug string) bool {
	if slug == "/" {
		return false
	}
	for i := 1; i < len(Cfg.Viewer.Navigation); i += 2 {
		if slug == Cfg.Viewer.Navigation[i] {
			return true
		}
	}
	return false
}
