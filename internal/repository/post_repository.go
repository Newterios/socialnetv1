package repository

import (
	"database/sql"
	"errors"
	"socialnet/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) (int64, error) {
	query := `INSERT INTO posts (user_id, content, media_url) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, post.UserID, post.Content, post.MediaURL)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *PostRepository) GetByID(id int64) (*model.Post, error) {
	query := `SELECT id, user_id, content, media_url, created_at, updated_at FROM posts WHERE id = ?`
	post := &model.Post{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.UserID, &post.Content, &post.MediaURL, &post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("post not found")
	}
	return post, err
}

func (r *PostRepository) Update(post *model.Post) error {
	query := `UPDATE posts SET content = ?, media_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := r.db.Exec(query, post.Content, post.MediaURL, post.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *PostRepository) Delete(id int64) error {
	query := `DELETE FROM posts WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *PostRepository) GetUserPosts(userID int64, limit int) ([]*model.Post, error) {
	query := `SELECT id, user_id, content, media_url, created_at, updated_at 
			  FROM posts WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) GetFeed(userID int64, limit int) ([]*model.Post, error) {
	query := `SELECT DISTINCT p.id, p.user_id, p.content, p.media_url, p.created_at, p.updated_at
			  FROM posts p
			  LEFT JOIN friendships f ON (f.requester_id = ? OR f.addressee_id = ?)
			  WHERE (p.user_id = ? OR p.user_id = f.requester_id OR p.user_id = f.addressee_id)
			  AND (f.status = 'accepted' OR p.user_id = ?)
			  ORDER BY p.created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, userID, userID, userID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r *PostRepository) scanPosts(rows *sql.Rows) ([]*model.Post, error) {
	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.MediaURL,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}
