package model

import "time"

type NotificationType string

const (
	NotificationFriendRequest NotificationType = "friend_request"
	NotificationLike          NotificationType = "like"
	NotificationComment       NotificationType = "comment"
	NotificationMessage       NotificationType = "message"
	NotificationGroupInvite   NotificationType = "group_invite"
)

type Notification struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"user_id"`
	Type      NotificationType `json:"type"`
	TargetID  int64            `json:"target_id,omitempty"`
	Message   string           `json:"message"`
	Read      bool             `json:"read"`
	CreatedAt time.Time        `json:"created_at"`
}
