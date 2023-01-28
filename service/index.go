package service

import (
	"blog/config"
	"blog/dao"
	"blog/model"
)

func GetIndexData(page, pageSize, postNum, pageNum int, slug string) (*model.HomeResponse, error) {
	categorys, err := dao.GetCategorys()
	if err != nil {
		return nil, err
	}

	var posts []model.Post
	if slug == "" {
		posts, err = dao.GetPostsAsPageDescByUpdateAt(page, pageSize)
	} else {
		posts, err = dao.GetPostsAsPageBySlugDescByUpdateAt(page, pageSize, slug)
	}

	if err != nil {
		return nil, err
	}
	postsMore := dao.GetPostMores(posts)

	var pages []int
	for i := page - 2; i <= page+2; i++ {
		if i > 0 && i <= pageNum {
			pages = append(pages, i)
		}
	}
	if page == 4 {
		pages = append([]int{1}, pages...)
	} else if page >= 5 {
		pages = append([]int{1, 0}, pages...)
	}
	if page == pageNum-3 {
		pages = append(pages, []int{pageNum}...)
	} else if page <= pageNum-3 {
		pages = append(pages, []int{0, pageNum}...)
	}
	var hr = &model.HomeResponse{
		Viewer:    config.Cfg.Viewer,
		Categorys: categorys,
		Posts:     postsMore,
		PostNum:   postNum,
		PageNum:   pageNum,
		Page:      page,
		Pages:     pages,
		PageEnd:   page != pageNum,
	}

	return hr, nil
}
