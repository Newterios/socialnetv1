package model

import "time"

type Group struct {
	ID          int64     `json:"id"`
	OwnerID     int64     `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	Owner       *User     `json:"owner,omitempty"`
	MemberCount int       `json:"member_count"`
	IsMember    bool      `json:"is_member"`
}

type GroupMember struct {
	ID       int64     `json:"id"`
	GroupID  int64     `json:"group_id"`
	UserID   int64     `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	User     *User     `json:"user,omitempty"`
}

type GroupPost struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	MediaURL  string    `json:"media_url"`
	CreatedAt time.Time `json:"created_at"`
	Author    *User     `json:"author,omitempty"`
}

type GroupCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GroupPostCreate struct {
	Content  string `json:"content"`
	MediaURL string `json:"media_url"`
}

type GroupSettings struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}
