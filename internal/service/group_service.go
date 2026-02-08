package service

import (
	"errors"
	"socialnet/internal/model"
	"socialnet/internal/repository"
	"socialnet/internal/security"
)

type GroupService struct {
	groupRepo  *repository.GroupRepository
	userRepo   *repository.UserRepository
	notifQueue chan *model.Notification
}

func NewGroupService(groupRepo *repository.GroupRepository, userRepo *repository.UserRepository,
	notifQueue chan *model.Notification) *GroupService {
	return &GroupService{
		groupRepo:  groupRepo,
		userRepo:   userRepo,
		notifQueue: notifQueue,
	}
}

func (s *GroupService) CreateGroup(ownerID int64, create *model.GroupCreate) (*model.Group, error) {
	if err := security.ValidateContent(create.Title, 100); err != nil {
		return nil, err
	}

	group := &model.Group{
		OwnerID:     ownerID,
		Title:       create.Title,
		Description: create.Description,
	}

	id, err := s.groupRepo.Create(group)
	if err != nil {
		return nil, err
	}

	group.ID = id
	return s.GetGroup(id, ownerID)
}

func (s *GroupService) GetGroup(groupID, userID int64) (*model.Group, error) {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return nil, err
	}

	owner, _ := s.userRepo.GetByID(group.OwnerID)
	group.Owner = owner

	count, _ := s.groupRepo.GetMemberCount(groupID)
	group.MemberCount = count

	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	group.IsMember = isMember

	return group, nil
}

func (s *GroupService) JoinGroup(groupID, userID int64) error {
	if _, err := s.groupRepo.GetByID(groupID); err != nil {
		return errors.New("group not found")
	}

	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if isMember {
		return errors.New("already a member")
	}

	return s.groupRepo.AddMember(groupID, userID)
}

func (s *GroupService) LeaveGroup(groupID, userID int64) error {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}

	if group.OwnerID == userID {
		return errors.New("owner cannot leave group")
	}

	return s.groupRepo.RemoveMember(groupID, userID)
}

func (s *GroupService) PostToGroup(groupID, userID int64, create *model.GroupPostCreate) (*model.GroupPost, error) {
	if err := security.ValidateContent(create.Content, 5000); err != nil {
		return nil, err
	}

	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("must be a member to post")
	}

	post := &model.GroupPost{
		GroupID:  groupID,
		UserID:   userID,
		Content:  create.Content,
		MediaURL: create.MediaURL,
	}

	id, err := s.groupRepo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	post, err = s.groupRepo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	author, _ := s.userRepo.GetByID(post.UserID)
	post.Author = author

	return post, nil
}

func (s *GroupService) GetGroupPosts(groupID, userID int64) ([]*model.GroupPost, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("must be a member to view posts")
	}

	posts, err := s.groupRepo.GetPosts(groupID, 50)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		author, _ := s.userRepo.GetByID(post.UserID)
		post.Author = author
	}

	return posts, nil
}

func (s *GroupService) GetUserGroups(userID int64) ([]*model.Group, error) {
	groups, err := s.groupRepo.GetUserGroups(userID)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		owner, _ := s.userRepo.GetByID(group.OwnerID)
		group.Owner = owner

		count, _ := s.groupRepo.GetMemberCount(group.ID)
		group.MemberCount = count

		group.IsMember = true
	}

	return groups, nil
}

func (s *GroupService) GetGroupMembers(groupID, userID int64) ([]*model.GroupMember, error) {
	isMember, _ := s.groupRepo.IsMember(groupID, userID)
	if !isMember {
		return nil, errors.New("must be a member to view members")
	}

	members, err := s.groupRepo.GetMembers(groupID)
	if err != nil {
		return nil, err
	}

	for _, member := range members {
		user, _ := s.userRepo.GetByID(member.UserID)
		member.User = user
	}

	return members, nil
}

func (s *GroupService) UpdateGroupSettings(groupID, userID int64, settings *model.GroupSettings) error {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return errors.New("group not found")
	}

	if group.OwnerID != userID {
		return errors.New("only owner can update settings")
	}

	if settings.Title != "" {
		if err := security.ValidateContent(settings.Title, 100); err != nil {
			return err
		}
		group.Title = settings.Title
	}

	if settings.Description != "" {
		group.Description = settings.Description
	}

	if settings.AvatarURL != "" {
		group.AvatarURL = settings.AvatarURL
	}

	return s.groupRepo.Update(group)
}

func (s *GroupService) IsOwner(groupID, userID int64) bool {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return false
	}
	return group.OwnerID == userID
}
