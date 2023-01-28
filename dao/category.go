package dao

import (
	"blog/model"
	"log"
)

func GetCategorys() ([]model.Category, error) {
	data, err := DB.Query("select * from category")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var cgs []model.Category
	for data.Next() {
		var cg model.Category
		err := data.Scan(&cg.Cid, &cg.Name, &cg.CreateAt, &cg.UpdateAt)
		if err != nil {
			return nil, err
		}
		cgs = append(cgs, cg)
	}
	return cgs, nil
}

func GetCategoryNameById(cid int) string {
	name := "null"
	data := DB.QueryRow("select name from category where cid = ?", cid)
	if data.Err() != nil {
		log.Println(data.Err())
		return name
	}
	_ = data.Scan(&name)
	return name
}

func NewCategory(c *model.CategoryDB) (int, error) {
	r, err := DB.Exec("insert into category(name, create_at, update_at) "+
		"values(?, ?, ?)", c.Name, c.CreateAt, c.UpdateAt)
	if err != nil {
		return -1, err
	}
	pid, err := r.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(pid), nil
}

func CheckCategoryExistByName(cn string) bool {
	data, err := DB.Query("select cid from category where category.name = ?", cn)
	if err != nil {
		log.Println(err)
		return false
	}
	return data.Next()
}

func CheckCategoryExistByCid(cid int) bool {
	data, err := DB.Query("select cid from category where cid = ?", cid)
	if err != nil {
		log.Println(err)
		return false
	}
	return data.Next()
}

func UpdCategory(c *model.CategoryDB) error {
	_, err := DB.Exec("update category set name = ?, update_at = ? where cid = ?",
		c.Name, c.CreateAt, c.Cid)
	if err != nil {
		return err
	}
	return nil
}

func DelCategory(cid int) error {
	_, err := DB.Exec("delete from category where cid = ?", cid)
	if err != nil {
		return err
	}
	return nil
}
