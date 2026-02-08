package repository

import (
	"database/sql"
	"errors"
	"socialnet/internal/model"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) (int64, error) {
	query := `INSERT INTO users (email, username, password_hash, full_name, bio, avatar_url, emoji_avatar, is_admin, show_last_seen, allow_messages_from, firebase_uid) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Email, user.Username, user.PasswordHash,
		user.FullName, user.Bio, user.AvatarURL, user.EmojiAvatar, user.IsAdmin,
		"all", "all", user.FirebaseUID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	query := `SELECT id, email, username, password_hash, full_name, bio, avatar_url, 
			  COALESCE(emoji_avatar, ''), is_admin, COALESCE(is_online, 0), last_seen,
			  COALESCE(show_last_seen, 'all'), COALESCE(allow_messages_from, 'all'), firebase_uid, created_at 
			  FROM users WHERE id = ?`
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FullName, &user.Bio, &user.AvatarURL, &user.EmojiAvatar, &user.IsAdmin,
		&user.IsOnline, &user.LastSeen, &user.ShowLastSeen, &user.AllowMessagesFrom,
		&user.FirebaseUID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `SELECT id, email, username, password_hash, full_name, bio, avatar_url,
			  COALESCE(emoji_avatar, ''), is_admin, COALESCE(is_online, 0), last_seen,
			  COALESCE(show_last_seen, 'all'), COALESCE(allow_messages_from, 'all'), firebase_uid, created_at 
			  FROM users WHERE email = ?`
	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FullName, &user.Bio, &user.AvatarURL, &user.EmojiAvatar, &user.IsAdmin,
		&user.IsOnline, &user.LastSeen, &user.ShowLastSeen, &user.AllowMessagesFrom,
		&user.FirebaseUID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	query := `SELECT id, email, username, password_hash, full_name, bio, avatar_url,
			  COALESCE(emoji_avatar, ''), is_admin, COALESCE(is_online, 0), last_seen,
			  COALESCE(show_last_seen, 'all'), COALESCE(allow_messages_from, 'all'), firebase_uid, created_at 
			  FROM users WHERE username = ?`
	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FullName, &user.Bio, &user.AvatarURL, &user.EmojiAvatar, &user.IsAdmin,
		&user.IsOnline, &user.LastSeen, &user.ShowLastSeen, &user.AllowMessagesFrom,
		&user.FirebaseUID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) GetByFirebaseUID(uid string) (*model.User, error) {
	query := `SELECT id, email, username, password_hash, full_name, bio, avatar_url,
			  COALESCE(emoji_avatar, ''), is_admin, COALESCE(is_online, 0), last_seen,
			  COALESCE(show_last_seen, 'all'), COALESCE(allow_messages_from, 'all'), firebase_uid, created_at 
			  FROM users WHERE firebase_uid = ?`
	user := &model.User{}
	err := r.db.QueryRow(query, uid).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.FullName, &user.Bio, &user.AvatarURL, &user.EmojiAvatar, &user.IsAdmin,
		&user.IsOnline, &user.LastSeen, &user.ShowLastSeen, &user.AllowMessagesFrom,
		&user.FirebaseUID, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) Update(user *model.User) error {
	query := `UPDATE users SET full_name = ?, bio = ?, avatar_url = ?, emoji_avatar = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.FullName, user.Bio, user.AvatarURL, user.EmojiAvatar, user.ID)
	return err
}

func (r *UserRepository) UpdateOnlineStatus(id int64, isOnline bool) error {
	query := `UPDATE users SET is_online = ?, last_seen = ? WHERE id = ?`
	_, err := r.db.Exec(query, isOnline, time.Now(), id)
	return err
}

func (r *UserRepository) UpdatePrivacySettings(id int64, settings *model.UserPrivacySettings) error {
	query := `UPDATE users SET show_last_seen = ?, allow_messages_from = ? WHERE id = ?`
	_, err := r.db.Exec(query, settings.ShowLastSeen, settings.AllowMessagesFrom, id)
	return err
}

func (r *UserRepository) UpdateEmojiAvatar(id int64, emoji string) error {
	query := `UPDATE users SET emoji_avatar = ? WHERE id = ?`
	_, err := r.db.Exec(query, emoji, id)
	return err
}

func (r *UserRepository) SetAdmin(id int64, isAdmin bool) error {
	query := `UPDATE users SET is_admin = ? WHERE id = ?`
	_, err := r.db.Exec(query, isAdmin, id)
	return err
}

func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) Search(searchTerm string, limit int) ([]*model.User, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, COALESCE(emoji_avatar, ''), is_admin, created_at 
			  FROM users WHERE username LIKE ? OR full_name LIKE ? LIMIT ?`
	pattern := "%" + searchTerm + "%"
	rows, err := r.db.Query(query, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FullName,
			&user.Bio, &user.AvatarURL, &user.EmojiAvatar, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) GetByEmojiAvatar(emoji string) ([]*model.User, error) {
	query := `SELECT id, email, username, full_name, bio, avatar_url, COALESCE(emoji_avatar, ''), is_admin, created_at 
			  FROM users WHERE emoji_avatar = ?`
	rows, err := r.db.Query(query, emoji)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FullName,
			&user.Bio, &user.AvatarURL, &user.EmojiAvatar, &user.IsAdmin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) CountTotal() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (r *UserRepository) CountOnline() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE is_online = 1`).Scan(&count)
	return count, err
}

func (r *UserRepository) GetAllUserIDs() ([]int64, error) {
	rows, err := r.db.Query(`SELECT id FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *UserRepository) LinkFirebaseUID(id int64, uid string) error {
	query := `UPDATE users SET firebase_uid = ? WHERE id = ?`
	_, err := r.db.Exec(query, uid, id)
	return err
}

func (r *UserRepository) GetUsersWithEmoji(emojiID string) ([]*model.User, error) {
	return r.GetByEmojiAvatar(emojiID)
}
