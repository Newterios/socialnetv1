package service

import (
	"errors"
	"socialnet/internal/model"
	"socialnet/internal/repository"
)

type AdminService struct {
	reportRepo  *repository.ReportRepository
	postRepo    *repository.PostRepository
	commentRepo *repository.CommentRepository
	userRepo    *repository.UserRepository
	statsRepo   *repository.StatsRepository
	notifQueue  chan *model.Notification
}

func NewAdminService(
	reportRepo *repository.ReportRepository,
	postRepo *repository.PostRepository,
	commentRepo *repository.CommentRepository,
	userRepo *repository.UserRepository,
	statsRepo *repository.StatsRepository,
	notifQueue chan *model.Notification,
) *AdminService {
	return &AdminService{
		reportRepo:  reportRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		statsRepo:   statsRepo,
		notifQueue:  notifQueue,
	}
}

func (s *AdminService) CreateReport(reporterID int64, create *model.ReportCreate) error {
	report := &model.Report{
		ReporterID: reporterID,
		TargetType: create.TargetType,
		TargetID:   create.TargetID,
		Reason:     create.Reason,
		Status:     model.ReportStatusPending,
	}

	_, err := s.reportRepo.Create(report)
	return err
}

func (s *AdminService) GetReports(status model.ReportStatus) ([]*model.Report, error) {
	reports, err := s.reportRepo.GetAll(status, 100)
	if err != nil {
		return nil, err
	}

	for _, report := range reports {
		reporter, _ := s.userRepo.GetByID(report.ReporterID)
		report.Reporter = reporter
	}

	return reports, nil
}

func (s *AdminService) ReviewReport(reportID int64, status model.ReportStatus) error {
	if status != model.ReportStatusReviewed && status != model.ReportStatusResolved {
		return errors.New("invalid status")
	}

	return s.reportRepo.UpdateStatus(reportID, status)
}

func (s *AdminService) DeleteContent(targetType model.ReportTargetType, targetID int64) error {
	switch targetType {
	case model.ReportTargetPost:
		return s.postRepo.Delete(targetID)
	case model.ReportTargetComment:
		return s.commentRepo.Delete(targetID)
	default:
		return errors.New("unsupported target type")
	}
}

func (s *AdminService) GetSiteStats() (*model.SiteStats, error) {
	return s.statsRepo.GetSiteStats()
}

func (s *AdminService) GrantAdmin(targetUserID int64) error {
	user, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.IsAdmin {
		return errors.New("user is already an admin")
	}

	return s.userRepo.SetAdmin(targetUserID, true)
}

func (s *AdminService) RevokeAdmin(targetUserID int64) error {
	user, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.IsAdmin {
		return errors.New("user is not an admin")
	}

	return s.userRepo.SetAdmin(targetUserID, false)
}

func (s *AdminService) BroadcastNotification(adminID int64, message string) error {
	if message == "" {
		return errors.New("message is required")
	}

	broadcast, err := s.statsRepo.CreateBroadcast(adminID, message)
	if err != nil {
		return err
	}

	userIDs, err := s.userRepo.GetAllUserIDs()
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		notification := &model.Notification{
			UserID:   userID,
			Type:     "broadcast",
			TargetID: broadcast.ID,
			Message:  message,
		}

		select {
		case s.notifQueue <- notification:
		default:
		}
	}

	return nil
}

func (s *AdminService) GetBroadcasts(limit int) ([]*model.Broadcast, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.statsRepo.GetBroadcasts(limit)
}

func (s *AdminService) BroadcastToEmoji(adminID int64, emojiID string, message string) error {
	if message == "" {
		return errors.New("message is required")
	}

	if emojiID == "" {
		return errors.New("emoji ID is required")
	}

	broadcast, err := s.statsRepo.CreateBroadcast(adminID, message)
	if err != nil {
		return err
	}

	users, err := s.userRepo.GetUsersWithEmoji(emojiID)
	if err != nil {
		return err
	}

	for _, user := range users {
		notification := &model.Notification{
			UserID:   user.ID,
			Type:     "broadcast",
			TargetID: broadcast.ID,
			Message:  message,
		}

		select {
		case s.notifQueue <- notification:
		default:
		}
	}

	return nil
}
