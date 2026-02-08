package model

import "time"

type ReportTargetType string
type ReportStatus string

const (
	ReportTargetPost    ReportTargetType = "post"
	ReportTargetComment ReportTargetType = "comment"
	ReportTargetUser    ReportTargetType = "user"

	ReportStatusPending  ReportStatus = "pending"
	ReportStatusReviewed ReportStatus = "reviewed"
	ReportStatusResolved ReportStatus = "resolved"
)

type Report struct {
	ID         int64            `json:"id"`
	ReporterID int64            `json:"reporter_id"`
	TargetType ReportTargetType `json:"target_type"`
	TargetID   int64            `json:"target_id"`
	Reason     string           `json:"reason"`
	Status     ReportStatus     `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	Reporter   *User            `json:"reporter,omitempty"`
}

type ReportCreate struct {
	TargetType ReportTargetType `json:"target_type"`
	TargetID   int64            `json:"target_id"`
	Reason     string           `json:"reason"`
}
