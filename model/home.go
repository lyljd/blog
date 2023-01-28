package model

import "blog/config"

type HomeResponse struct {
	config.Viewer
	Categorys []Category
	Posts     []PostMore
	PostNum   int
	PageNum   int
	Page      int
	Pages     []int
	PageEnd   bool
}
