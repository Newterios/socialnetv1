package repository

import (
	"database/sql"
	"socialnet/internal/model"
)

type FriendshipRepository struct {
	db *sql.DB
}

func NewFriendshipRepository(db *sql.DB) *FriendshipRepository {
	return &FriendshipRepository{db: db}
}

func (r *FriendshipRepository) CreateRequest(requesterID, addresseeID int64) (int64, error) {
	query := `INSERT INTO friendships (requester_id, addressee_id, status) VALUES (?, ?, 'pending')`
	result, err := r.db.Exec(query, requesterID, addresseeID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *FriendshipRepository) UpdateStatus(id int64, status model.FriendshipStatus) error {
	query := `UPDATE friendships SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *FriendshipRepository) GetByID(id int64) (*model.Friendship, error) {
	query := `SELECT id, requester_id, addressee_id, status, created_at, updated_at 
			  FROM friendships WHERE id = ?`
	friendship := &model.Friendship{}
	err := r.db.QueryRow(query, id).Scan(
		&friendship.ID, &friendship.RequesterID, &friendship.AddresseeID,
		&friendship.Status, &friendship.CreatedAt, &friendship.UpdatedAt,
	)
	return friendship, err
}

func (r *FriendshipRepository) GetFriends(userID int64) ([]*model.User, error) {
	query := `SELECT u.id, u.email, u.username, u.full_name, u.bio, u.avatar_url, u.is_admin, u.created_at
			  FROM users u
			  INNER JOIN friendships f ON (f.requester_id = u.id OR f.addressee_id = u.id)
			  WHERE (f.requester_id = ? OR f.addressee_id = ?) 
			  AND f.status = 'accepted' AND u.id != ?`
	rows, err := r.db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FullName,
			&user.Bio, &user.AvatarURL, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *FriendshipRepository) GetPendingRequests(userID int64) ([]*model.Friendship, error) {
	query := `SELECT id, requester_id, addressee_id, status, created_at, updated_at 
			  FROM friendships WHERE addressee_id = ? AND status = 'pending' ORDER BY created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*model.Friendship
	for rows.Next() {
		friendship := &model.Friendship{}
		err := rows.Scan(&friendship.ID, &friendship.RequesterID, &friendship.AddresseeID,
			&friendship.Status, &friendship.CreatedAt, &friendship.UpdatedAt)
		if err != nil {
			return nil, err
		}
		friendships = append(friendships, friendship)
	}
	return friendships, rows.Err()
}

func (r *FriendshipRepository) AreFriends(userID1, userID2 int64) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM friendships 
		WHERE ((requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?))
		AND status = 'accepted'
	)`
	var exists bool
	err := r.db.QueryRow(query, userID1, userID2, userID2, userID1).Scan(&exists)
	return exists, err
}

func (r *FriendshipRepository) GetFriendship(userID1, userID2 int64) (*model.Friendship, error) {
	query := `SELECT id, requester_id, addressee_id, status, created_at, updated_at 
			  FROM friendships 
			  WHERE (requester_id = ? AND addressee_id = ?) OR (requester_id = ? AND addressee_id = ?)`
	friendship := &model.Friendship{}
	err := r.db.QueryRow(query, userID1, userID2, userID2, userID1).Scan(
		&friendship.ID, &friendship.RequesterID, &friendship.AddresseeID,
		&friendship.Status, &friendship.CreatedAt, &friendship.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return friendship, nil
}
