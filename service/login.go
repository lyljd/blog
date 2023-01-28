package service

import (
	"blog/common"
	"blog/config"
	"blog/model"
)

func GetLoginData(un, pw string) (*model.LoginRes, string) {
	if un != config.Cfg.System.Username || pw != config.Cfg.System.Password {
		return nil, "账号或密码错误"
	}
	return &model.LoginRes{
		Nickname: config.Cfg.System.Nickname,
		Token:    common.NewToken(),
	}, ""
}
