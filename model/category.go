package model

import "time"

type Category struct {
	Cid      int    `json:"cid"`
	Name     string `json:"name"`
	CreateAt string `json:"createAt"`
	UpdateAt string `json:"updateAt"`
}

type CategoryDB struct {
	Cid      int
	Name     string
	CreateAt time.Time
	UpdateAt time.Time
}

type CategoryResponse struct {
	HomeResponse
	CategoryName string
}
