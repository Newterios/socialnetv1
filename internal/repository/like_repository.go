package repository

import (
	"database/sql"
	"socialnet/internal/model"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Create(like *model.Like) error {
	query := `INSERT INTO likes (post_id, user_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, like.PostID, like.UserID)
	return err
}

func (r *LikeRepository) Delete(postID, userID int64) error {
	query := `DELETE FROM likes WHERE post_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, postID, userID)
	return err
}

func (r *LikeRepository) GetCountByPostID(postID int64) (int, error) {
	query := `SELECT COUNT(*) FROM likes WHERE post_id = ?`
	var count int
	err := r.db.QueryRow(query, postID).Scan(&count)
	return count, err
}

func (r *LikeRepository) HasUserLiked(postID, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE post_id = ? AND user_id = ?)`
	var exists bool
	err := r.db.QueryRow(query, postID, userID).Scan(&exists)
	return exists, err
}
