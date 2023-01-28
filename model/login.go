package model

type LoginRes struct {
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
}
