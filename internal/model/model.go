package model

import "time"

type Post struct {
	ID        string    `json:"id"`
	Caption   string    `json:"caption"`
	Author    string    `json:"author"`
	Pics      []Pic     `json:"pics"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

type Pic struct {
	ID            string `json:"id"`
	HDPath        string `json:"path"`
	ThumbnailPath string `json:"thumbnail"`
}
