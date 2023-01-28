package dao

import (
	"blog/common"
	"blog/model"
	"errors"
	"html/template"
	"log"
	"strconv"
)

func GetPostsAsPageDescByUpdateAt(page, pageSize int) ([]model.Post, error) {
	page = (page - 1) * pageSize
	data, err := DB.Query("select * from post order by update_at desc limit ?, ?", page, pageSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ps []model.Post
	for data.Next() {
		var p model.Post
		err := data.Scan(
			&p.Pid,
			&p.Title,
			&p.Content,
			&p.Markdown,
			&p.CategoryId,
			&p.ViewCount,
			&p.Type,
			&p.Slug,
			&p.CreateAt,
			&p.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func GetPostsAsPageBySlugDescByUpdateAt(page, pageSize int, slug string) ([]model.Post, error) {
	page = (page - 1) * pageSize
	data, err := DB.Query("select * from post where type = 1 and slug = ? order by update_at desc limit ?, ?", slug, page, pageSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ps []model.Post
	for data.Next() {
		var p model.Post
		err := data.Scan(
			&p.Pid,
			&p.Title,
			&p.Content,
			&p.Markdown,
			&p.CategoryId,
			&p.ViewCount,
			&p.Type,
			&p.Slug,
			&p.CreateAt,
			&p.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func GetPostsDescByCreateAt() ([]model.Post, error) {
	data, err := DB.Query("select * from post order by create_at desc")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ps []model.Post
	for data.Next() {
		var p model.Post
		err := data.Scan(
			&p.Pid,
			&p.Title,
			&p.Content,
			&p.Markdown,
			&p.CategoryId,
			&p.ViewCount,
			&p.Type,
			&p.Slug,
			&p.CreateAt,
			&p.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func GetPostsByCidDescByUpdateAt(page, pageSize, cid int) ([]model.Post, error) {
	page = (page - 1) * pageSize
	data, err := DB.Query("select * from post where category_id = ? order by update_at desc limit ?, ?", cid, page, pageSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var ps []model.Post
	for data.Next() {
		var p model.Post
		err := data.Scan(
			&p.Pid,
			&p.Title,
			&p.Content,
			&p.Markdown,
			&p.CategoryId,
			&p.ViewCount,
			&p.Type,
			&p.Slug,
			&p.CreateAt,
			&p.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func GetPostMores(posts []model.Post) []model.PostMore {
	var pms []model.PostMore
	for _, p := range posts {
		content := []rune(common.TrimHtml(p.Content))
		if len(content) > 150 {
			content = content[:150]
		}
		pm := model.PostMore{
			Pid:          p.Pid,
			Title:        p.Title,
			Slug:         p.Slug,
			Content:      template.HTML(content),
			CategoryId:   p.CategoryId,
			CategoryName: GetCategoryNameById(p.CategoryId),
			ViewCount:    p.ViewCount,
			Type:         p.Type,
			CreateAt:     common.DateDay(p.CreateAt),
			UpdateAt:     common.DateDay(p.UpdateAt),
		}
		pms = append(pms, pm)
	}
	return pms
}

func GetPostNum() (num int) {
	data := DB.QueryRow("select count(1) from post")
	err := data.Scan(&num)
	if err != nil {
		log.Println(err)
	}
	return
}

func GetPostNumBySlug(slug string) (num int) {
	data := DB.QueryRow("select count(1) from post where type = 1 and slug = ?", slug)
	err := data.Scan(&num)
	if err != nil {
		log.Println(err)
	}
	return
}

func GetPostNumByCid(cid int) (num int) {
	data := DB.QueryRow("select count(1) from post where category_id = ?", cid)
	err := data.Scan(&num)
	if err != nil {
		log.Println(err)
	}
	return
}

func GetPostByPid(pid int) (*model.Post, error) {
	data, err := DB.Query("select * from post where pid = ?", pid)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if !data.Next() {
		return nil, errors.New("在GetPostByPid中没有找到pid为" + strconv.Itoa(pid) + "的记录")
	}
	var p model.Post
	err = data.Scan(
		&p.Pid,
		&p.Title,
		&p.Content,
		&p.Markdown,
		&p.CategoryId,
		&p.ViewCount,
		&p.Type,
		&p.Slug,
		&p.CreateAt,
		&p.UpdateAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPostMore(p *model.Post) *model.PostMore {
	PostViewCountInc(p.Pid)
	return &model.PostMore{
		Pid:          p.Pid,
		Title:        p.Title,
		Slug:         p.Slug,
		Content:      template.HTML(p.Content),
		CategoryId:   p.CategoryId,
		CategoryName: GetCategoryNameById(p.CategoryId),
		ViewCount:    p.ViewCount + 1,
		Type:         p.Type,
		CreateAt:     common.DateDay(p.CreateAt),
		UpdateAt:     common.DateDay(p.UpdateAt),
	}
}

func SavePost(p *model.Post) (int, error) {
	r, err := DB.Exec("insert into post"+
		"(title, content, markdown, category_id, view_count, type, slug, create_at, update_at)"+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", p.Title, p.Content, p.Markdown, p.CategoryId,
		p.ViewCount, p.Type, p.Slug, p.CreateAt, p.UpdateAt)
	if err != nil {
		return -1, err
	}
	pid, err := r.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(pid), nil
}

func SetPost(p *model.Post) (int, error) {
	_, err := DB.Exec("update post set "+
		"title = ?, content = ?, markdown = ?, category_id = ?, type = ?, slug = ?, "+
		"update_at = ? where pid = ?", p.Title, p.Content, p.Markdown, p.CategoryId,
		p.Type, p.Slug, p.UpdateAt, p.Pid)
	if err != nil {
		return -1, err
	}
	return p.Pid, nil
}

func PostViewCountInc(pid int) {
	_, _ = DB.Exec("update post set view_count = view_count+1 where pid = ?", pid)
}

func SearchPosts(k string) ([]model.SearchResp, error) {
	data, err := DB.Query("select pid, title from post where title like ?", "%"+k+"%")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var srs []model.SearchResp
	for data.Next() {
		var sr model.SearchResp
		err := data.Scan(&sr.Pid, &sr.Title)
		if err != nil {
			return nil, err
		}
		srs = append(srs, sr)
	}

	return srs, nil
}

func DeletePostById(pid int) error {
	_, err := DB.Exec("delete from post where pid = ?", pid)
	if err != nil {
		return err
	}
	return nil
}

func RemovePostCid(cid int) error {
	_, err := DB.Exec("update post set category_id = 0 where category_id = ?", cid)
	if err != nil {
		return err
	}
	return nil
}
