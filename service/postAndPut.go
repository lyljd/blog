package service

import (
	"blog/dao"
	"blog/model"
	"log"
	"time"
)

func PostAndPut(pr *model.PostReq) (int, string) {
	if pr.Type == 1 && pr.Slug == "" {
		return -1, "请输入自定义链接"
	}
	pm := &model.Post{
		Title:      pr.Title,
		Slug:       pr.Slug,
		Content:    pr.Content,
		Markdown:   pr.Markdown,
		CategoryId: pr.CategoryId,
		ViewCount:  0,
		Type:       pr.Type,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	}

	var pid int
	var err error
	if pr.Pid > 0 {
		pm.Pid = pr.Pid
		pid, err = dao.SetPost(pm)
	} else {
		pid, err = dao.SavePost(pm)
	}
	if err != nil {
		log.Println(err)
		return pid, "提交失败"
	}
	return pid, ""
}
