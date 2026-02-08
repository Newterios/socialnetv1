package model

import "time"

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	MediaURL  string    `json:"media_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    *User     `json:"author,omitempty"`
	LikeCount int       `json:"like_count"`
	Liked     bool      `json:"liked"`
}

type PostCreate struct {
	Content  string `json:"content"`
	MediaURL string `json:"media_url,omitempty"`
}

type PostUpdate struct {
	Content  string `json:"content"`
	MediaURL string `json:"media_url,omitempty"`
}
