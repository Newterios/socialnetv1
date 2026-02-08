package service

import (
	"socialnet/internal/model"
	"socialnet/internal/repository"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) CreateNotification(notification *model.Notification) error {
	_, err := s.notifRepo.Create(notification)
	return err
}

func (s *NotificationService) GetNotifications(userID int64) ([]*model.Notification, error) {
	return s.notifRepo.GetByUser(userID, 50)
}

func (s *NotificationService) MarkAsRead(notificationID int64) error {
	return s.notifRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) GetUnreadCount(userID int64) (int, error) {
	return s.notifRepo.GetUnreadCount(userID)
}

func (s *NotificationService) ClearNotifications(userID int64) error {
	return s.notifRepo.DeleteByUser(userID)
}
