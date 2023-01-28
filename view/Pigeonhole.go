package view

import (
	"blog/model"
	"blog/service"
	"net/http"
)

func (*html) Pigeonhole(w http.ResponseWriter, _ *http.Request) {
	t := model.Template.Pigeonhole

	data, err := service.GetPigeonholeData()
	if err != nil {
		model.WriteError(w, err)
		return
	}
	t.WriteData(w, data)
}
