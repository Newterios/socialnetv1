package model

import "time"

type Broadcast struct {
	ID        int64     `json:"id"`
	AdminID   int64     `json:"admin_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type BroadcastCreate struct {
	Message string `json:"message"`
}

type SiteStats struct {
	TotalUsers         int64 `json:"total_users"`
	TotalPosts         int64 `json:"total_posts"`
	TotalGroups        int64 `json:"total_groups"`
	TotalMessages      int64 `json:"total_messages"`
	OnlineUsers        int64 `json:"online_users"`
	PendingReports     int64 `json:"pending_reports"`
	TotalNotifications int64 `json:"total_notifications"`
}

type GrantAdminRequest struct {
	UserID int64 `json:"user_id"`
}
