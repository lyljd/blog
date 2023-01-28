package view

import (
	"blog/dao"
	"blog/model"
	"blog/service"
	"net/http"
	"strconv"
	"strings"
)

func (*html) Category(w http.ResponseWriter, r *http.Request) {
	t := model.Template.Category

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

	cid, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/c/"))
	if err != nil {
		model.WriteError(w, err, "cid不存在", 400)
		return
	}

	name := dao.GetCategoryNameById(cid)
	if name == "null" {
		model.WriteError(w, err, "cid不存在", 400)
		return
	}

	postNum, pageNum := dao.GetPostNumByCid(cid), 1
	if postNum > 0 {
		pageNum = (postNum-1)/pageSize + 1
	}
	if page < 1 || page > pageNum {
		model.WriteError(w, nil, "page不存在", 400)
		return
	}

	data, err := service.GetCategoryData(page, pageSize, postNum, pageNum, cid, name)
	if err != nil {
		model.WriteError(w, err)
		return
	}
	data.Title = data.CategoryName + " - " + data.Title
	t.WriteData(w, data)
}
