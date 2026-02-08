package repository

import (
	"database/sql"
	"errors"
	"socialnet/internal/model"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *model.Comment) (int64, error) {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, comment.PostID, comment.UserID, comment.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *CommentRepository) GetByID(id int64) (*model.Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = ?`
	comment := &model.Comment{}
	err := r.db.QueryRow(query, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("comment not found")
	}
	return comment, err
}

func (r *CommentRepository) GetByPostID(postID int64) ([]*model.Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at 
			  FROM comments WHERE post_id = ? ORDER BY created_at ASC`
	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		comment := &model.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

func (r *CommentRepository) Delete(id int64) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
