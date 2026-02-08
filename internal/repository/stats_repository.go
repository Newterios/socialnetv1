package repository

import (
	"database/sql"
	"socialnet/internal/model"
)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetSiteStats() (*model.SiteStats, error) {
	stats := &model.SiteStats{}

	r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	r.db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&stats.TotalPosts)
	r.db.QueryRow(`SELECT COUNT(*) FROM groups`).Scan(&stats.TotalGroups)
	r.db.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&stats.TotalMessages)
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_online = 1`).Scan(&stats.OnlineUsers)
	r.db.QueryRow(`SELECT COUNT(*) FROM reports WHERE status = 'pending'`).Scan(&stats.PendingReports)
	r.db.QueryRow(`SELECT COUNT(*) FROM notifications`).Scan(&stats.TotalNotifications)

	return stats, nil
}

func (r *StatsRepository) CreateBroadcast(adminID int64, message string) (*model.Broadcast, error) {
	result, err := r.db.Exec(`INSERT INTO broadcasts (admin_id, message) VALUES (?, ?)`, adminID, message)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	var broadcast model.Broadcast
	err = r.db.QueryRow(`SELECT id, admin_id, message, created_at FROM broadcasts WHERE id = ?`, id).Scan(
		&broadcast.ID, &broadcast.AdminID, &broadcast.Message, &broadcast.CreatedAt,
	)
	return &broadcast, err
}

func (r *StatsRepository) GetBroadcasts(limit int) ([]*model.Broadcast, error) {
	rows, err := r.db.Query(`SELECT id, admin_id, message, created_at FROM broadcasts ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var broadcasts []*model.Broadcast
	for rows.Next() {
		b := &model.Broadcast{}
		if err := rows.Scan(&b.ID, &b.AdminID, &b.Message, &b.CreatedAt); err != nil {
			return nil, err
		}
		broadcasts = append(broadcasts, b)
	}
	return broadcasts, rows.Err()
}
