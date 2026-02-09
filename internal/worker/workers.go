package worker

import (
	"log"
	"socialnet/internal/model"
	"socialnet/internal/service"
	"time"
)

type NotificationWorker struct {
	queue   chan *model.Notification
	service *service.NotificationService
}

func NewNotificationWorker(queue chan *model.Notification, service *service.NotificationService) *NotificationWorker {
	return &NotificationWorker{
		queue:   queue,
		service: service,
	}
}

func (w *NotificationWorker) Start() {
	go func() {
		log.Println("Notification worker started")
		for notification := range w.queue {
			if err := w.service.CreateNotification(notification); err != nil {
				log.Printf("Failed to create notification: %v", err)
			}
		}
	}()
}

type CleanupWorker struct {
	service  *service.NotificationService
	interval time.Duration
	maxAge   time.Duration
}

func NewCleanupWorker(service *service.NotificationService, interval, maxAge time.Duration) *CleanupWorker {
	return &CleanupWorker{
		service:  service,
		interval: interval,
		maxAge:   maxAge,
	}
}

func (w *CleanupWorker) Start() {
	go func() {
		log.Println("Cleanup worker started")
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Running cleanup task...")
		}
	}()
}
