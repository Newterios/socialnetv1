package repository

import (
	"database/sql"
	"socialnet/internal/model"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateConversation() (int64, error) {
	query := `INSERT INTO conversations DEFAULT VALUES`
	result, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MessageRepository) AddMember(conversationID, userID int64) error {
	query := `INSERT INTO conversation_members (conversation_id, user_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, conversationID, userID)
	return err
}

func (r *MessageRepository) GetConversationBetween(user1ID, user2ID int64) (*model.Conversation, error) {
	query := `SELECT c.id, c.created_at FROM conversations c
			  INNER JOIN conversation_members cm1 ON cm1.conversation_id = c.id
			  INNER JOIN conversation_members cm2 ON cm2.conversation_id = c.id
			  WHERE cm1.user_id = ? AND cm2.user_id = ?
			  LIMIT 1`
	conversation := &model.Conversation{}
	err := r.db.QueryRow(query, user1ID, user2ID).Scan(&conversation.ID, &conversation.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return conversation, err
}

func (r *MessageRepository) CreateMessage(message *model.Message) (int64, error) {
	query := `INSERT INTO messages (conversation_id, user_id, body) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, message.ConversationID, message.UserID, message.Body)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MessageRepository) GetMessages(conversationID int64, limit int) ([]*model.Message, error) {
	query := `SELECT id, conversation_id, user_id, body, created_at, read_at 
			  FROM messages WHERE conversation_id = ? ORDER BY created_at ASC, id ASC LIMIT ?`
	rows, err := r.db.Query(query, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		message := &model.Message{}
		err := rows.Scan(&message.ID, &message.ConversationID, &message.UserID,
			&message.Body, &message.CreatedAt, &message.ReadAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *MessageRepository) GetUserConversations(userID int64) ([]*model.Conversation, error) {
	query := `SELECT c.id, c.created_at,
			  u.id, u.email, u.username, u.full_name, u.bio, u.avatar_url, u.is_admin, u.created_at,
			  m.id, m.conversation_id, m.user_id, m.body, m.created_at, m.read_at
			  FROM conversations c
			  INNER JOIN conversation_members cm ON cm.conversation_id = c.id AND cm.user_id = ?
			  LEFT JOIN conversation_members cm_other ON cm_other.conversation_id = c.id AND cm_other.user_id != ?
			  LEFT JOIN users u ON u.id = cm_other.user_id
			  LEFT JOIN messages m ON m.id = (
				SELECT id FROM messages WHERE conversation_id = c.id ORDER BY created_at DESC, id DESC LIMIT 1
			  )
			  ORDER BY COALESCE(m.created_at, c.created_at) DESC`
	rows, err := r.db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*model.Conversation
	for rows.Next() {
		conversation := &model.Conversation{}
		participant := &model.User{}
		var participantID sql.NullInt64
		var participantEmail sql.NullString
		var participantUsername sql.NullString
		var participantFullName sql.NullString
		var participantBio sql.NullString
		var participantAvatarURL sql.NullString
		var participantIsAdmin sql.NullBool
		var participantCreatedAt sql.NullTime

		lastMessage := &model.Message{}
		var messageID sql.NullInt64
		var messageConversationID sql.NullInt64
		var messageUserID sql.NullInt64
		var messageBody sql.NullString
		var messageCreatedAt sql.NullTime
		var messageReadAt sql.NullTime

		err := rows.Scan(
			&conversation.ID, &conversation.CreatedAt,
			&participantID, &participantEmail, &participantUsername, &participantFullName,
			&participantBio, &participantAvatarURL, &participantIsAdmin, &participantCreatedAt,
			&messageID, &messageConversationID, &messageUserID, &messageBody, &messageCreatedAt, &messageReadAt,
		)
		if err != nil {
			return nil, err
		}

		if participantID.Valid {
			participant.ID = participantID.Int64
			participant.Email = participantEmail.String
			participant.Username = participantUsername.String
			participant.FullName = participantFullName.String
			participant.Bio = participantBio.String
			participant.AvatarURL = participantAvatarURL.String
			participant.IsAdmin = participantIsAdmin.Bool
			if participantCreatedAt.Valid {
				participant.CreatedAt = participantCreatedAt.Time
			}
			conversation.Participant = participant
		}

		if messageID.Valid {
			lastMessage.ID = messageID.Int64
			lastMessage.ConversationID = messageConversationID.Int64
			lastMessage.UserID = messageUserID.Int64
			lastMessage.Body = messageBody.String
			if messageCreatedAt.Valid {
				lastMessage.CreatedAt = messageCreatedAt.Time
			}
			if messageReadAt.Valid {
				readAt := messageReadAt.Time
				lastMessage.ReadAt = &readAt
			}
			conversation.LastMessage = lastMessage
		}

		conversations = append(conversations, conversation)
	}
	return conversations, rows.Err()
}

func (r *MessageRepository) IsMember(conversationID, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM conversation_members WHERE conversation_id = ? AND user_id = ?)`
	var exists bool
	err := r.db.QueryRow(query, conversationID, userID).Scan(&exists)
	return exists, err
}
