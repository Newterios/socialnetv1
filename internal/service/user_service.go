package service

import (
	"errors"
	"socialnet/internal/model"
	"socialnet/internal/repository"
)

type UserService struct {
	userRepo   *repository.UserRepository
	friendRepo *repository.FriendshipRepository
}

func NewUserService(userRepo *repository.UserRepository, friendRepo *repository.FriendshipRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		friendRepo: friendRepo,
	}
}

func (s *UserService) GetProfile(userID int64) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) GetPublicProfile(targetID, viewerID int64) (*model.UserPublic, error) {
	user, err := s.userRepo.GetByID(targetID)
	if err != nil {
		return nil, err
	}

	isFriend := false
	if viewerID != targetID {
		friendship, _ := s.friendRepo.GetFriendship(viewerID, targetID)
		if friendship != nil && friendship.Status == model.FriendshipAccepted {
			isFriend = true
		}
	}

	return user.ToPublic(viewerID, isFriend), nil
}

func (s *UserService) UpdateProfile(userID int64, profile *model.UserProfile) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	user.FullName = profile.FullName
	user.Bio = profile.Bio
	user.AvatarURL = profile.AvatarURL
	if profile.EmojiAvatar != "" {
		user.EmojiAvatar = profile.EmojiAvatar
	}

	return s.userRepo.Update(user)
}

func (s *UserService) UpdatePrivacySettings(userID int64, settings *model.UserPrivacySettings) error {
	validLastSeen := map[string]bool{"all": true, "friends": true, "nobody": true}
	validMessages := map[string]bool{"all": true, "friends": true}

	if !validLastSeen[settings.ShowLastSeen] {
		return errors.New("invalid show_last_seen value")
	}
	if !validMessages[settings.AllowMessagesFrom] {
		return errors.New("invalid allow_messages_from value")
	}

	return s.userRepo.UpdatePrivacySettings(userID, settings)
}

func (s *UserService) SetEmojiAvatar(userID int64, emojiID string) error {
	if !model.IsValidEmoji(emojiID) {
		return errors.New("invalid emoji")
	}

	return s.userRepo.UpdateEmojiAvatar(userID, emojiID)
}

func (s *UserService) GetUsersWithEmoji(emojiID string) ([]*model.EmojiUserInfo, error) {
	if !model.IsValidEmoji(emojiID) {
		return nil, errors.New("invalid emoji")
	}

	users, err := s.userRepo.GetByEmojiAvatar(emojiID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.EmojiUserInfo, len(users))
	for i, u := range users {
		result[i] = &model.EmojiUserInfo{
			UserID:      u.ID,
			Username:    u.Username,
			FullName:    u.FullName,
			EmojiAvatar: u.EmojiAvatar,
		}
	}
	return result, nil
}

func (s *UserService) UpdateOnlineStatus(userID int64, isOnline bool) error {
	return s.userRepo.UpdateOnlineStatus(userID, isOnline)
}

func (s *UserService) DeleteAccount(userID int64) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(userID)
}

func (s *UserService) SearchUsers(searchTerm string) ([]*model.User, error) {
	if searchTerm == "" {
		return nil, errors.New("search term is required")
	}
	return s.userRepo.Search(searchTerm, 20)
}

func (s *UserService) CanMessageUser(senderID, recipientID int64) (bool, error) {
	recipient, err := s.userRepo.GetByID(recipientID)
	if err != nil {
		return false, err
	}

	if recipient.AllowMessagesFrom == "all" {
		return true, nil
	}

	friendship, _ := s.friendRepo.GetFriendship(senderID, recipientID)
	if friendship != nil && friendship.Status == model.FriendshipAccepted {
		return true, nil
	}

	return false, errors.New("user only accepts messages from friends")
}

func (s *UserService) GetAllEmojis() []model.Emoji {
	return model.PredefinedEmojis
}
