package model

type Backend struct {
	ID     uint   `json:"id"`
	ApiId  uint   `json:"api_id"`
	Url    string `json:"url"`
	Weight uint   `json:"weight"`
}
