package model

import "time"

type FriendshipStatus string

const (
	FriendshipPending  FriendshipStatus = "pending"
	FriendshipAccepted FriendshipStatus = "accepted"
	FriendshipBlocked  FriendshipStatus = "blocked"
)

type Friendship struct {
	ID          int64            `json:"id"`
	RequesterID int64            `json:"requester_id"`
	AddresseeID int64            `json:"addressee_id"`
	Status      FriendshipStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Requester   *User            `json:"requester,omitempty"`
	Addressee   *User            `json:"addressee,omitempty"`
}

type FriendRequest struct {
	AddresseeID int64 `json:"addressee_id"`
}
