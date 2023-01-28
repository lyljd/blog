package service

import (
	"blog/config"
	"blog/dao"
	"blog/model"
)

func GetPigeonholeData() (*model.PigeonholeResponse, error) {
	categorys, err := dao.GetCategorys()
	if err != nil {
		return nil, err
	}

	posts, err := dao.GetPostsDescByCreateAt()
	if err != nil {
		return nil, err
	}
	lines := make(map[string][]model.PostPigeonhole)
	for _, p := range posts {
		month := p.CreateAt.Format("2006年1月")
		lines[month] = append(lines[month], model.PostPigeonhole{
			Pid:      p.Pid,
			Title:    p.Title,
			CreateAt: p.CreateAt.Format("02日 15:04:05"),
		})
	}

	var pr = &model.PigeonholeResponse{
		Viewer:    config.Cfg.Viewer,
		Categorys: categorys,
		Lines:     lines,
	}

	return pr, nil
}
