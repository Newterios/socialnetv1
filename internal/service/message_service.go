package service

import (
	"errors"
	"socialnet/internal/model"
	"socialnet/internal/repository"
	"socialnet/internal/security"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	friendRepo  *repository.FriendshipRepository
	userRepo    *repository.UserRepository
	notifQueue  chan *model.Notification
}

func NewMessageService(messageRepo *repository.MessageRepository, friendRepo *repository.FriendshipRepository,
	userRepo *repository.UserRepository, notifQueue chan *model.Notification) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		friendRepo:  friendRepo,
		userRepo:    userRepo,
		notifQueue:  notifQueue,
	}
}

func (s *MessageService) StartConversation(user1ID, user2ID int64) (*model.Conversation, error) {
	if user1ID == user2ID {
		return nil, errors.New("cannot message yourself")
	}

	areFriends, _ := s.friendRepo.AreFriends(user1ID, user2ID)
	if !areFriends {
		return nil, errors.New("can only message friends")
	}

	conversation, err := s.messageRepo.GetConversationBetween(user1ID, user2ID)
	if err != nil {
		return nil, err
	}

	if conversation != nil {
		return conversation, nil
	}

	convID, err := s.messageRepo.CreateConversation()
	if err != nil {
		return nil, err
	}

	if err := s.messageRepo.AddMember(convID, user1ID); err != nil {
		return nil, err
	}

	if err := s.messageRepo.AddMember(convID, user2ID); err != nil {
		return nil, err
	}

	return &model.Conversation{ID: convID}, nil
}

func (s *MessageService) SendMessage(conversationID, userID int64, create *model.MessageCreate) (*model.Message, error) {
	if err := security.ValidateContent(create.Body, 2000); err != nil {
		return nil, err
	}

	isMember, _ := s.messageRepo.IsMember(conversationID, userID)
	if !isMember {
		return nil, errors.New("not a member of this conversation")
	}

	message := &model.Message{
		ConversationID: conversationID,
		UserID:         userID,
		Body:           create.Body,
	}

	id, err := s.messageRepo.CreateMessage(message)
	if err != nil {
		return nil, err
	}

	message.ID = id

	sender, _ := s.userRepo.GetByID(userID)
	notifMessage := sender.Username + " sent you a message"

	s.notifQueue <- &model.Notification{
		UserID:   0,
		Type:     model.NotificationMessage,
		TargetID: conversationID,
		Message:  notifMessage,
	}

	return message, nil
}

func (s *MessageService) GetMessages(conversationID, userID int64) ([]*model.Message, error) {
	isMember, _ := s.messageRepo.IsMember(conversationID, userID)
	if !isMember {
		return nil, errors.New("not a member of this conversation")
	}

	messages, err := s.messageRepo.GetMessages(conversationID, 100)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		author, _ := s.userRepo.GetByID(message.UserID)
		message.Author = author
	}

	return messages, nil
}

func (s *MessageService) GetConversations(userID int64) ([]*model.Conversation, error) {
	return s.messageRepo.GetUserConversations(userID)
}
