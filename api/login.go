package api

import (
	"blog/common"
	"blog/model"
	"blog/service"
	"net/http"
)

func (*api) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		model.WriteError(w, nil, "未找到页面", 404)
		return
	}

	jp := common.GetRequestJsonParam(r)
	un := jp["username"].(string)
	pw := jp["password"].(string)

	data, err := service.GetLoginData(un, pw)
	if err != "" {
		model.RespJson(w, 400, data, err)
		return
	}
	model.RespJsonSucc(w, *data)
}
