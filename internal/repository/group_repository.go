package repository

import (
	"database/sql"
	"errors"
	"socialnet/internal/model"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *model.Group) (int64, error) {
	query := `INSERT INTO groups (owner_id, title, description, avatar_url) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, group.OwnerID, group.Title, group.Description, group.AvatarURL)
	if err != nil {
		return 0, err
	}
	groupID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := r.AddMember(groupID, group.OwnerID); err != nil {
		return 0, err
	}

	return groupID, nil
}

func (r *GroupRepository) GetByID(id int64) (*model.Group, error) {
	query := `SELECT id, owner_id, title, description, COALESCE(avatar_url, ''), created_at FROM groups WHERE id = ?`
	group := &model.Group{}
	err := r.db.QueryRow(query, id).Scan(
		&group.ID, &group.OwnerID, &group.Title, &group.Description, &group.AvatarURL, &group.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("group not found")
	}
	return group, err
}

func (r *GroupRepository) Update(group *model.Group) error {
	query := `UPDATE groups SET title = ?, description = ?, avatar_url = ? WHERE id = ?`
	_, err := r.db.Exec(query, group.Title, group.Description, group.AvatarURL, group.ID)
	return err
}

func (r *GroupRepository) AddMember(groupID, userID int64) error {
	query := `INSERT INTO group_members (group_id, user_id, role) VALUES (?, ?, 'member')`
	_, err := r.db.Exec(query, groupID, userID)
	return err
}

func (r *GroupRepository) RemoveMember(groupID, userID int64) error {
	query := `DELETE FROM group_members WHERE group_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, groupID, userID)
	return err
}

func (r *GroupRepository) IsMember(groupID, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)`
	var exists bool
	err := r.db.QueryRow(query, groupID, userID).Scan(&exists)
	return exists, err
}

func (r *GroupRepository) GetMemberCount(groupID int64) (int, error) {
	query := `SELECT COUNT(*) FROM group_members WHERE group_id = ?`
	var count int
	err := r.db.QueryRow(query, groupID).Scan(&count)
	return count, err
}

func (r *GroupRepository) GetMembers(groupID int64) ([]*model.GroupMember, error) {
	query := `SELECT id, group_id, user_id, COALESCE(role, 'member'), joined_at 
			  FROM group_members WHERE group_id = ? ORDER BY joined_at ASC`
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*model.GroupMember
	for rows.Next() {
		m := &model.GroupMember{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *GroupRepository) CreatePost(post *model.GroupPost) (int64, error) {
	query := `INSERT INTO group_posts (group_id, user_id, content, media_url) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, post.GroupID, post.UserID, post.Content, post.MediaURL)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *GroupRepository) GetPostByID(id int64) (*model.GroupPost, error) {
	query := `SELECT id, group_id, user_id, content, COALESCE(media_url, ''), created_at FROM group_posts WHERE id = ?`
	post := &model.GroupPost{}
	err := r.db.QueryRow(query, id).Scan(
		&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.MediaURL, &post.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("group post not found")
	}
	return post, err
}

func (r *GroupRepository) GetPosts(groupID int64, limit int) ([]*model.GroupPost, error) {
	query := `SELECT id, group_id, user_id, content, COALESCE(media_url, ''), created_at 
			  FROM group_posts WHERE group_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.Query(query, groupID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.GroupPost
	for rows.Next() {
		post := &model.GroupPost{}
		err := rows.Scan(&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.MediaURL, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (r *GroupRepository) GetUserGroups(userID int64) ([]*model.Group, error) {
	query := `SELECT g.id, g.owner_id, g.title, g.description, COALESCE(g.avatar_url, ''), g.created_at
			  FROM groups g
			  INNER JOIN group_members gm ON gm.group_id = g.id
			  WHERE gm.user_id = ? ORDER BY g.created_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		group := &model.Group{}
		err := rows.Scan(&group.ID, &group.OwnerID, &group.Title, &group.Description, &group.AvatarURL, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, rows.Err()
}
