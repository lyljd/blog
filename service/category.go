package service

import (
	"blog/config"
	"blog/dao"
	"blog/model"
	"log"
	"time"
)

func GetCategoryData(page, pageSize, postNum, pageNum, cid int, name string) (*model.CategoryResponse, error) {
	categorys, err := dao.GetCategorys()
	if err != nil {
		return nil, err
	}

	posts, err := dao.GetPostsByCidDescByUpdateAt(page, pageSize, cid)
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
	var cr = &model.CategoryResponse{
		HomeResponse: model.HomeResponse{
			Viewer:    config.Cfg.Viewer,
			Categorys: categorys,
			Posts:     postsMore,
			PostNum:   postNum,
			PageNum:   pageNum,
			Page:      page,
			Pages:     pages,
			PageEnd:   page != pageNum,
		},
		CategoryName: name,
	}

	return cr, nil
}

func NewCategory(cn string) (int, string) {
	if len(cn) == 0 {
		return -1, "分类名不能为空"
	}
	if len(cn) > 10 {
		return -1, "分类名最长10位"
	}
	if dao.CheckCategoryExistByName(cn) {
		return -1, "分类名已存在"
	}

	cid, err := dao.NewCategory(&model.CategoryDB{
		Name:     cn,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return cid, "新增失败"
	}
	return cid, ""
}

func UpdAndDelCategoryByCid(cid int, newCn string) string {
	if !dao.CheckCategoryExistByCid(cid) {
		return "欲操作的分类id不存在"
	}

	var err error
	if newCn != "" {
		if len(newCn) == 0 {
			return "分类名不能为空"
		}
		if len(newCn) > 10 {
			return "分类名最长10位"
		}
		if dao.CheckCategoryExistByName(newCn) {
			return "分类名已存在"
		}

		err = dao.UpdCategory(&model.CategoryDB{
			Cid:      cid,
			Name:     newCn,
			CreateAt: time.Now(),
		})
		if err != nil {
			log.Println(err)
			return "修改失败"
		}
	} else {
		err = dao.DelCategory(cid)
		if err != nil {
			log.Println(err)
			return "删除失败"
		}
		err = dao.RemovePostCid(cid)
		if err != nil {
			log.Println(err)
			return "删除成功，但所属此分类的所有文章将被取消分类失败，请手动更新"
		}
	}

	return ""
}
