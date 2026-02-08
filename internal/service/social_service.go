package service

import (
	"errors"
	"socialnet/internal/model"
	"socialnet/internal/repository"
	"socialnet/internal/security"
)

type SocialService struct {
	friendRepo  *repository.FriendshipRepository
	likeRepo    *repository.LikeRepository
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
	userRepo    *repository.UserRepository
	notifQueue  chan *model.Notification
}

func NewSocialService(friendRepo *repository.FriendshipRepository, likeRepo *repository.LikeRepository,
	commentRepo *repository.CommentRepository, postRepo *repository.PostRepository,
	userRepo *repository.UserRepository, notifQueue chan *model.Notification) *SocialService {
	return &SocialService{
		friendRepo:  friendRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		notifQueue:  notifQueue,
	}
}

func (s *SocialService) SendFriendRequest(requesterID, addresseeID int64) error {
	if requesterID == addresseeID {
		return errors.New("cannot send friend request to yourself")
	}

	if _, err := s.userRepo.GetByID(addresseeID); err != nil {
		return errors.New("user not found")
	}

	areFriends, _ := s.friendRepo.AreFriends(requesterID, addresseeID)
	if areFriends {
		return errors.New("already friends")
	}

	id, err := s.friendRepo.CreateRequest(requesterID, addresseeID)
	if err != nil {
		return err
	}

	requester, _ := s.userRepo.GetByID(requesterID)
	message := requester.Username + " sent you a friend request"

	s.notifQueue <- &model.Notification{
		UserID:   addresseeID,
		Type:     model.NotificationFriendRequest,
		TargetID: id,
		Message:  message,
	}

	return nil
}

func (s *SocialService) AcceptFriendRequest(requestID, userID int64) error {
	friendship, err := s.friendRepo.GetByID(requestID)
	if err != nil {
		return err
	}

	if friendship.AddresseeID != userID {
		return errors.New("unauthorized")
	}

	if friendship.Status != model.FriendshipPending {
		return errors.New("request already processed")
	}

	return s.friendRepo.UpdateStatus(requestID, model.FriendshipAccepted)
}

func (s *SocialService) BlockUser(requestID, userID int64) error {
	friendship, err := s.friendRepo.GetByID(requestID)
	if err != nil {
		return err
	}

	if friendship.AddresseeID != userID {
		return errors.New("unauthorized")
	}

	return s.friendRepo.UpdateStatus(requestID, model.FriendshipBlocked)
}

func (s *SocialService) GetFriends(userID int64) ([]*model.User, error) {
	return s.friendRepo.GetFriends(userID)
}

func (s *SocialService) GetPendingRequests(userID int64) ([]*model.Friendship, error) {
	friendships, err := s.friendRepo.GetPendingRequests(userID)
	if err != nil {
		return nil, err
	}

	for _, friendship := range friendships {
		requester, _ := s.userRepo.GetByID(friendship.RequesterID)
		friendship.Requester = requester
	}

	return friendships, nil
}

func (s *SocialService) LikePost(postID, userID int64) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	liked, _ := s.likeRepo.HasUserLiked(postID, userID)
	if liked {
		return errors.New("already liked")
	}

	like := &model.Like{
		PostID: postID,
		UserID: userID,
	}

	if err := s.likeRepo.Create(like); err != nil {
		return err
	}

	if post.UserID != userID {
		liker, _ := s.userRepo.GetByID(userID)
		message := liker.Username + " liked your post"

		s.notifQueue <- &model.Notification{
			UserID:   post.UserID,
			Type:     model.NotificationLike,
			TargetID: postID,
			Message:  message,
		}
	}

	return nil
}

func (s *SocialService) UnlikePost(postID, userID int64) error {
	return s.likeRepo.Delete(postID, userID)
}

func (s *SocialService) CommentOnPost(postID, userID int64, create *model.CommentCreate) (*model.Comment, error) {
	if err := security.ValidateContent(create.Content, 1000); err != nil {
		return nil, err
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	comment := &model.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: create.Content,
	}

	id, err := s.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	comment, err = s.commentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	author, _ := s.userRepo.GetByID(comment.UserID)
	comment.Author = author

	if post.UserID != userID {
		commenter, _ := s.userRepo.GetByID(userID)
		message := commenter.Username + " commented on your post"

		s.notifQueue <- &model.Notification{
			UserID:   post.UserID,
			Type:     model.NotificationComment,
			TargetID: postID,
			Message:  message,
		}
	}

	return comment, nil
}

func (s *SocialService) GetComments(postID int64) ([]*model.Comment, error) {
	comments, err := s.commentRepo.GetByPostID(postID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		author, _ := s.userRepo.GetByID(comment.UserID)
		comment.Author = author
	}

	return comments, nil
}
