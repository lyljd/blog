package view

import (
	"blog/config"
	"blog/dao"
	"blog/model"
	"blog/service"
	"net/http"
	"strconv"
	"strings"
)

func (*html) Index(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	isSlug := config.CheckSlugExist(r.URL.Path)
	if !(isSlug || url == "/" || (len(url) > 7 && url[:7] == "/?page=")) {
		model.WriteError(w, nil, "未找到页面", 404)
		return
	}
	t := model.Template.Index

	if err := r.ParseForm(); err != nil {
		model.WriteError(w, err, "解析负载失败")
		return
	}
	pageStr, page, pageSize := r.Form.Get("page"), 1, 10
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			model.WriteError(w, err, "page不存在", 400)
			return
		}
	}

	postNum, pageNum := 0, 1
	slug := strings.TrimPrefix(r.URL.Path, "/")
	if !isSlug {
		postNum = dao.GetPostNum()
	} else {
		postNum = dao.GetPostNumBySlug(slug)
	}
	if postNum > 0 {
		pageNum = (postNum-1)/pageSize + 1
	}
	if page < 1 || page > pageNum {
		model.WriteError(w, nil, "page不存在", 400)
		return
	}

	var data *model.HomeResponse
	var err error
	if !isSlug {
		data, err = service.GetIndexData(page, pageSize, postNum, pageNum, "")
	} else {
		data, err = service.GetIndexData(page, pageSize, postNum, pageNum, slug)
	}

	if err != nil {
		model.WriteError(w, err)
		return
	}
	t.WriteData(w, data)
}
