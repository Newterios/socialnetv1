package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int64          `json:"id"`
	Email             string         `json:"email"`
	Username          string         `json:"username"`
	PasswordHash      string         `json:"-"`
	FullName          string         `json:"full_name"`
	Bio               string         `json:"bio"`
	AvatarURL         string         `json:"avatar_url"`
	EmojiAvatar       string         `json:"emoji_avatar"`
	IsAdmin           bool           `json:"is_admin"`
	IsOnline          bool           `json:"is_online"`
	LastSeen          sql.NullTime   `json:"last_seen"`
	ShowLastSeen      string         `json:"show_last_seen"`
	AllowMessagesFrom string         `json:"allow_messages_from"`
	FirebaseUID       sql.NullString `json:"-"`
	CreatedAt         time.Time      `json:"created_at"`
}

type UserRegistration struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	FullName    string `json:"full_name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	EmojiAvatar string `json:"emoji_avatar"`
}

type UserPrivacySettings struct {
	ShowLastSeen      string `json:"show_last_seen"`
	AllowMessagesFrom string `json:"allow_messages_from"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token"`
}

type UserPublic struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	EmojiAvatar string    `json:"emoji_avatar"`
	IsOnline    bool      `json:"is_online"`
	LastSeen    string    `json:"last_seen,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (u *User) ToPublic(viewerID int64, isFriend bool) *UserPublic {
	pub := &UserPublic{
		ID:          u.ID,
		Username:    u.Username,
		FullName:    u.FullName,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
		EmojiAvatar: u.EmojiAvatar,
		IsOnline:    u.IsOnline,
		CreatedAt:   u.CreatedAt,
	}

	showLastSeen := false
	switch u.ShowLastSeen {
	case "all":
		showLastSeen = true
	case "friends":
		showLastSeen = isFriend || viewerID == u.ID
	case "nobody":
		showLastSeen = viewerID == u.ID
	}

	if showLastSeen && u.LastSeen.Valid {
		pub.LastSeen = u.LastSeen.Time.Format(time.RFC3339)
	}

	return pub
}
