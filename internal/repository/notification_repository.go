package repository

import (
	"database/sql"
	"socialnet/internal/model"
	"time"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *model.Notification) (int64, error) {
	query := `INSERT INTO notifications (user_id, type, target_id, message) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, notification.UserID, notification.Type,
		notification.TargetID, notification.Message)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *NotificationRepository) GetByUser(userID int64, limit int) ([]*model.Notification, error) {
	query := `SELECT id, user_id, type, target_id, message, read, created_at 
			  FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		notification := &model.Notification{}
		var targetID sql.NullInt64
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.Type,
			&targetID, &notification.Message, &notification.Read, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}
		if targetID.Valid {
			notification.TargetID = targetID.Int64
		}
		notifications = append(notifications, notification)
	}
	return notifications, rows.Err()
}

func (r *NotificationRepository) MarkAsRead(id int64) error {
	query := `UPDATE notifications SET read = TRUE WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *NotificationRepository) DeleteByUser(userID int64) error {
	query := `DELETE FROM notifications WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *NotificationRepository) DeleteOld(olderThan time.Duration) error {
	query := `DELETE FROM notifications WHERE read = TRUE AND created_at < ?`
	cutoff := time.Now().Add(-olderThan)
	_, err := r.db.Exec(query, cutoff)
	return err
}

func (r *NotificationRepository) GetUnreadCount(userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read = FALSE`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}
