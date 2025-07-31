package model

type Post struct {
	ID      string   `json:"id"`
	Caption string   `json:"caption"`
	Pics    []string `json:"pics"`
}
