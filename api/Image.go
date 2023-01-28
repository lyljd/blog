package api

import (
	"blog/common"
	"blog/model"
	"blog/service"
	"net/http"
	"path"
)

func (*api) Image(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		model.WriteError(w, nil, "未找到页面", 404)
		return
	}

	token := r.Header.Get("Token")
	ok, errStr := common.CheckToken(token)
	if !ok {
		model.RespJson(w, 401, nil, errStr)
		return
	}

	file, head, err := r.FormFile("image")
	if err != nil {
		model.RespJson(w, 500, nil, "解析图片失败")
		return
	}

	fn, errStr := service.SaveImage(file, path.Ext(path.Base(head.Filename)))
	if errStr != "" {
		model.RespJson(w, 500, nil, errStr)
		return
	}
	model.RespJsonSucc(w, struct {
		Fn string `json:"fn"`
	}{fn})
}
