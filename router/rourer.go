package router

import (
	"blog/api"
	"blog/view"
	"net/http"
)

func Router() {
	http.HandleFunc("/", view.HTML.Index)
	http.HandleFunc("/c/", view.HTML.Category)
	http.HandleFunc("/p/", view.HTML.Post)
	http.HandleFunc("/login", view.HTML.Login)
	http.HandleFunc("/writing", view.HTML.Writing)
	http.HandleFunc("/pigeonhole", view.HTML.Pigeonhole)

	http.HandleFunc("/api/v1/login", api.API.Login)
	http.HandleFunc("/api/v1/image", api.API.Image)
	http.HandleFunc("/api/v1/post", api.API.Post)
	http.HandleFunc("/api/v1/post/", api.API.PostGet)
	http.HandleFunc("/api/v1/post/search", api.API.PostSearch)
	http.HandleFunc("/api/v1/category", api.API.NewCategory)
	http.HandleFunc("/api/v1/category/", api.API.UpdAndDelCategory)

	http.Handle("/resource/image/", http.StripPrefix("/resource/image/", http.FileServer(http.Dir("blog_image/"))))
	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("public/resource/"))))
}
