package view

import (
	"blog/config"
	"blog/dao"
	"blog/model"
	"net/http"
	"strconv"
	"strings"
)

func (*html) Post(w http.ResponseWriter, r *http.Request) {
	t := model.Template.Detail

	pid, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/p/"))
	if err != nil {
		model.WriteError(w, err, "pid不存在", 400)
		return
	}

	post, err := dao.GetPostByPid(pid)
	if err != nil {
		model.WriteError(w, err, "pid不存在", 400)
		return
	}

	cv := config.Cfg.Viewer
	cv.Title = post.Title + " - " + cv.Title
	data := &model.PostRes{
		Viewer:       cv,
		SystemConfig: config.Cfg.System,
		Article:      *dao.GetPostMore(post),
	}
	t.WriteData(w, data)
}
