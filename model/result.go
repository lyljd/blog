package model

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	Code  int    `json:"code"`
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func RespJson(w http.ResponseWriter, code int, data any, err string) {
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal(Result{
		Code:  code,
		Data:  data,
		Error: err,
	})
	_, _ = w.Write(res)
}

func RespJsonSucc(w http.ResponseWriter, data any) {
	RespJson(w, 200, data, "")
}
