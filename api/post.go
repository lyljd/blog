package api

import (
	"blog/common"
	"blog/dao"
	"blog/model"
	"blog/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (*api) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
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
	check := []string{"title", "slug", "content", "markdown", "categoryId", "type"}
	if r.Method == http.MethodPut {
		check = append(check, "pid")
	}
	for _, c := range check {
		if jp[c] == nil {
			model.RespJson(w, 400, nil, "提交的数据缺少"+c)
			return
		}
	}
	pr := &model.PostReq{
		Title:      jp["title"].(string),
		Slug:       jp["slug"].(string),
		Content:    jp["content"].(string),
		Markdown:   jp["markdown"].(string),
		CategoryId: int(jp["categoryId"].(float64)),
		Type:       int(jp["type"].(float64)),
	}
	if r.Method == http.MethodPut {
		pr.Pid = int(jp["pid"].(float64))
		if pr.Pid <= 0 {
			model.RespJson(w, 400, nil, "pid不存在")
			return
		}
	}

	pid, err := service.PostAndPut(pr)
	if err != "" {
		model.RespJson(w, 400, nil, err)
		return
	}
	model.RespJsonSucc(w, struct {
		Pid int `json:"pid"`
	}{pid})
}

func (*api) PostGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodDelete {
		model.WriteError(w, nil, "未找到页面", 404)
		return
	}

	token := r.Header.Get("Token")
	ok, errStr := common.CheckToken(token)
	if !ok {
		model.RespJson(w, 401, nil, errStr)
		return
	}

	pidStr := strings.TrimPrefix(r.URL.Path, "/api/v1/post/")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		model.WriteError(w, err, "pid不存在")
		return
	}

	if r.Method == http.MethodGet {
		p, err := dao.GetPostByPid(pid)
		if err != nil {
			log.Println(err)
			model.RespJson(w, 500, nil, "获取帖子数据失败")
			return
		}
		model.RespJsonSucc(w, *p)
	} else if r.Method == http.MethodDelete {
		err := dao.DeletePostById(pid)
		if err != nil {
			log.Println(err)
			model.RespJson(w, 500, nil, "删除失败")
			return
		}
		model.RespJsonSucc(w, nil)
	}
}

func (*api) PostSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		model.WriteError(w, nil, "未找到页面", 404)
		return
	}

	if err := r.ParseForm(); err != nil {
		model.WriteError(w, err, "解析负载失败")
		return
	}

	keyword := r.Form.Get("k")

	ps, err := dao.SearchPosts(keyword)
	if err != nil {
		log.Println(err)
		model.RespJson(w, 500, nil, "搜索失败")
		return
	}
	model.RespJsonSucc(w, ps)
}
