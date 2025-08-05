package model

import "time"

type Post struct {
	ID        string    `json:"id"`
	Caption   string    `json:"caption"`
	Pics      []string  `json:"pics"`
	Author    string    `json:"author"`
	Hearts    int       `json:"hearts"`
	CreatedAt time.Time `json:"created_at"`
}
