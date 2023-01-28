package api

import (
	"blog/common"
	"blog/model"
	"blog/service"
	"net/http"
	"strconv"
	"strings"
)

func (*api) NewCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		model.WriteError(w, nil, "未找到页面", 404)
		return
	}

	token := r.Header.Get("Token")
	ok, err := common.CheckToken(token)
	if !ok {
		model.RespJson(w, 401, nil, err)
		return
	}

	jp := common.GetRequestJsonParam(r)
	if jp["cn"] == nil {
		model.RespJson(w, 400, nil, "提交的数据缺少cn")
		return
	}

	cid, err := service.NewCategory(strings.TrimSpace(jp["cn"].(string)))
	if err != "" {
		model.RespJson(w, 400, nil, err)
		return
	}
	model.RespJsonSucc(w, struct {
		Cid int `json:"cid"`
	}{cid})
}

func (*api) UpdAndDelCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodDelete {
		model.WriteError(w, nil, "未找到页面", 404)
		return
	}

	token := r.Header.Get("Token")
	ok, errStr := common.CheckToken(token)
	if !ok {
		model.RespJson(w, 401, nil, errStr)
		return
	}

	cid, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/category/"))
	if err != nil {
		model.WriteError(w, err, "cid不存在", 400)
		return
	}

	var newCn string
	if r.Method == http.MethodPut {
		jp := common.GetRequestJsonParam(r)
		if jp["cn"] == nil {
			model.RespJson(w, 400, nil, "提交的数据缺少cn")
			return
		}
		newCn = strings.TrimSpace(jp["cn"].(string))
	}

	errStr = service.UpdAndDelCategoryByCid(cid, newCn)
	if errStr != "" {
		model.RespJson(w, 400, nil, errStr)
		return
	}
	model.RespJsonSucc(w, nil)
}
