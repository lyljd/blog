package view

import (
	"blog/config"
	"blog/model"
	"net/http"
)

func (*html) Login(w http.ResponseWriter, _ *http.Request) {
	t := model.Template.Login

	cv := config.Cfg.Viewer
	cv.Title = "登录 - " + cv.Title
	t.WriteData(w, cv)
}
