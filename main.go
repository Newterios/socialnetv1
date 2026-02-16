package main

import (
	"log"
	"net/http"
	"os"
	"socialnet/internal/config"
	"socialnet/internal/database"
	httpRouter "socialnet/internal/http"
	httpHandler "socialnet/internal/http/handler"
	httpMiddleware "socialnet/internal/http/middleware"
	"socialnet/internal/model"
	"socialnet/internal/repository"
	"socialnet/internal/security"
	"socialnet/internal/service"
	"socialnet/internal/worker"
	"time"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database initialized successfully")

	if cfg.UploadDir != "" {
		if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
			log.Printf("Warning: Failed to create upload directory: %v", err)
		}
	}

	var mongodb *database.MongoDB
	if cfg.MongoDBURI != "" {
		mongodb, err = database.NewMongoDB(cfg.MongoDBURI)
		if err != nil {
			log.Printf("Warning: Failed to connect to MongoDB: %v", err)
		} else {
			defer mongodb.Close()
			log.Println("MongoDB connected successfully")
		}
	}

	var firebaseAuth *security.FirebaseAuth
	if cfg.FirebaseKeyPath != "" {
		firebaseAuth, err = security.NewFirebaseAuth(cfg.FirebaseKeyPath)
		if err != nil {
			log.Printf("Warning: Failed to initialize Firebase: %v", err)
		} else {
			log.Println("Firebase initialized successfully")
		}
	}

	if cfg.UploadDir != "" {
		if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
			log.Printf("Warning: Failed to create upload directory: %v", err)
		}
	}

	userRepo := repository.NewUserRepository(db.DB)
	postRepo := repository.NewPostRepository(db.DB)
	commentRepo := repository.NewCommentRepository(db.DB)
	likeRepo := repository.NewLikeRepository(db.DB)
	friendRepo := repository.NewFriendshipRepository(db.DB)
	messageRepo := repository.NewMessageRepository(db.DB)
	groupRepo := repository.NewGroupRepository(db.DB)
	notifRepo := repository.NewNotificationRepository(db.DB)
	reportRepo := repository.NewReportRepository(db.DB)
	statsRepo := repository.NewStatsRepository(db.DB)

	notifQueue := make(chan *model.Notification, 100)

	authService := service.NewAuthService(userRepo, firebaseAuth, cfg.InitialAdmins)
	userService := service.NewUserService(userRepo, friendRepo)
	postService := service.NewPostService(postRepo, likeRepo, userRepo)
	socialService := service.NewSocialService(friendRepo, likeRepo, commentRepo, postRepo, userRepo, notifQueue)
	messageService := service.NewMessageService(messageRepo, friendRepo, userRepo, notifQueue)
	groupService := service.NewGroupService(groupRepo, userRepo, notifQueue)
	notifService := service.NewNotificationService(notifRepo)
	adminService := service.NewAdminService(reportRepo, postRepo, commentRepo, userRepo, statsRepo, notifQueue)

	authHandler := httpHandler.NewAuthHandler(authService, cfg.JWTSecret, cfg.SessionDuration)
	userHandler := httpHandler.NewUserHandler(userService)
	postHandler := httpHandler.NewPostHandler(postService)
	socialHandler := httpHandler.NewSocialHandler(socialService)
	messageHandler := httpHandler.NewMessageHandler(messageService)
	groupHandler := httpHandler.NewGroupHandler(groupService)
	notifHandler := httpHandler.NewNotificationHandler(notifService)
	adminHandler := httpHandler.NewAdminHandler(adminService)

	authMiddleware := httpMiddleware.NewAuthMiddleware(cfg.JWTSecret)
	rateLimiter := httpMiddleware.NewRateLimiter(cfg.RateLimitPerMin, time.Minute)

	router := httpRouter.NewRouter(
		authHandler, userHandler, postHandler, socialHandler,
		messageHandler, groupHandler, notifHandler, adminHandler,
		authMiddleware, rateLimiter, cfg.UploadDir, cfg.FrontendDir,
	)

	notifWorker := worker.NewNotificationWorker(notifQueue, notifService)
	notifWorker.Start()

	cleanupWorker := worker.NewCleanupWorker(notifService, cfg.CleanupInterval, 7*24*time.Hour)
	cleanupWorker.Start()

	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router.Setup()))
}
