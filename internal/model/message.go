package model

import "time"

type Conversation struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Members     []*User   `json:"members,omitempty"`
	Participant *User     `json:"participant,omitempty"`
	LastMessage *Message  `json:"last_message,omitempty"`
}

type Message struct {
	ID             int64      `json:"id"`
	ConversationID int64      `json:"conversation_id"`
	UserID         int64      `json:"user_id"`
	Body           string     `json:"body"`
	CreatedAt      time.Time  `json:"created_at"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	Author         *User      `json:"author,omitempty"`
}

type MessageCreate struct {
	Body string `json:"body"`
}

type ConversationCreate struct {
	ParticipantID int64 `json:"participant_id"`
}
