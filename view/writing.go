package view

import (
	"blog/config"
	"blog/dao"
	"blog/model"
	"net/http"
)

func (*html) Writing(w http.ResponseWriter, _ *http.Request) {
	t := model.Template.Writing

	cgs, err := dao.GetCategorys()
	if err != nil {
		model.WriteError(w, err)
		return
	}

	data := &model.WritingRes{
		Title:     config.Cfg.Viewer.Title,
		Categorys: cgs,
	}
	t.WriteData(w, data)
}
