package service

import (
	"errors"
	"socialnet/internal/model"
	"socialnet/internal/repository"
	"socialnet/internal/security"
)

type PostService struct {
	postRepo *repository.PostRepository
	likeRepo *repository.LikeRepository
	userRepo *repository.UserRepository
}

func NewPostService(postRepo *repository.PostRepository, likeRepo *repository.LikeRepository,
	userRepo *repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		likeRepo: likeRepo,
		userRepo: userRepo,
	}
}

func (s *PostService) CreatePost(userID int64, create *model.PostCreate) (*model.Post, error) {
	if err := security.ValidateContent(create.Content, 5000); err != nil {
		return nil, err
	}

	post := &model.Post{
		UserID:   userID,
		Content:  create.Content,
		MediaURL: create.MediaURL,
	}

	id, err := s.postRepo.Create(post)
	if err != nil {
		return nil, err
	}

	post.ID = id
	return s.GetPost(id, userID)
}

func (s *PostService) GetPost(postID, currentUserID int64) (*model.Post, error) {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	author, _ := s.userRepo.GetByID(post.UserID)
	post.Author = author

	count, _ := s.likeRepo.GetCountByPostID(postID)
	post.LikeCount = count

	liked, _ := s.likeRepo.HasUserLiked(postID, currentUserID)
	post.Liked = liked

	return post, nil
}

func (s *PostService) UpdatePost(postID, userID int64, update *model.PostUpdate) error {
	if err := security.ValidateContent(update.Content, 5000); err != nil {
		return err
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("unauthorized")
	}

	post.Content = update.Content
	post.MediaURL = update.MediaURL

	return s.postRepo.Update(post)
}

func (s *PostService) DeletePost(postID, userID int64, isAdmin bool) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	if post.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	return s.postRepo.Delete(postID)
}

func (s *PostService) GetFeed(userID int64) ([]*model.Post, error) {
	posts, err := s.postRepo.GetFeed(userID, 50)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		author, _ := s.userRepo.GetByID(post.UserID)
		post.Author = author

		count, _ := s.likeRepo.GetCountByPostID(post.ID)
		post.LikeCount = count

		liked, _ := s.likeRepo.HasUserLiked(post.ID, userID)
		post.Liked = liked
	}

	return posts, nil
}
