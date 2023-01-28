package model

import "blog/config"

type PigeonholeResponse struct {
	config.Viewer
	Categorys []Category
	Lines     map[string][]PostPigeonhole
}
