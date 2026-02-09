package http

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"socialnet/internal/http/handler"
	"socialnet/internal/http/middleware"
	"strings"
	"time"
)

type Router struct {
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	socialHandler  *handler.SocialHandler
	messageHandler *handler.MessageHandler
	groupHandler   *handler.GroupHandler
	notifHandler   *handler.NotificationHandler
	adminHandler   *handler.AdminHandler
	authMiddleware *middleware.AuthMiddleware
	rateLimiter    *middleware.RateLimiter
	uploadDir      string
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	socialHandler *handler.SocialHandler,
	messageHandler *handler.MessageHandler,
	groupHandler *handler.GroupHandler,
	notifHandler *handler.NotificationHandler,
	adminHandler *handler.AdminHandler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimiter *middleware.RateLimiter,
	uploadDir string,
) *Router {
	return &Router{
		authHandler:    authHandler,
		userHandler:    userHandler,
		postHandler:    postHandler,
		socialHandler:  socialHandler,
		messageHandler: messageHandler,
		groupHandler:   groupHandler,
		notifHandler:   notifHandler,
		adminHandler:   adminHandler,
		authMiddleware: authMiddleware,
		rateLimiter:    rateLimiter,
		uploadDir:      uploadDir,
	}
}

func (rt *Router) Setup() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", rt.authHandler.Register)
	mux.HandleFunc("/login", rt.authHandler.Login)
	mux.HandleFunc("/auth/google", rt.authHandler.GoogleLogin)
	mux.HandleFunc("/auth/password-requirements", rt.authHandler.GetPasswordRequirements)

	mux.Handle("/users/search", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.SearchUsers)))
	mux.Handle("/users/", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.GetProfile)))

	mux.Handle("/profile", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			rt.userHandler.UpdateProfile(w, r)
		case http.MethodDelete:
			rt.userHandler.DeleteAccount(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/profile/privacy", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.UpdatePrivacySettings)))
	mux.Handle("/profile/emoji", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.SetEmojiAvatar)))
	mux.Handle("/profile/status", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.UpdateOnlineStatus)))

	mux.Handle("/emojis", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.GetEmojis)))
	mux.Handle("/emojis/", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.GetUsersWithEmoji)))

	mux.Handle("/posts", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.postHandler.GetFeed(w, r)
		case http.MethodPost:
			rt.postHandler.CreatePost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/posts/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.postHandler.GetPost(w, r)
		case http.MethodPut:
			rt.postHandler.UpdatePost(w, r)
		case http.MethodDelete:
			rt.postHandler.DeletePost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/likes/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rt.socialHandler.LikePost(w, r)
		case http.MethodDelete:
			rt.socialHandler.UnlikePost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/comments/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.socialHandler.GetComments(w, r)
		case http.MethodPost:
			rt.socialHandler.CommentOnPost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/friends", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.socialHandler.GetFriends)))
	mux.Handle("/friends/requests", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.socialHandler.GetPendingRequests)))
	mux.Handle("/friends/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rt.socialHandler.SendFriendRequest(w, r)
		case http.MethodPut:
			rt.socialHandler.AcceptFriendRequest(w, r)
		case http.MethodDelete:
			rt.socialHandler.BlockUser(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/conversations", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.messageHandler.GetConversations(w, r)
		case http.MethodPost:
			rt.messageHandler.StartConversation(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/conversations/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/messages") {
			switch r.Method {
			case http.MethodGet:
				rt.messageHandler.GetMessages(w, r)
			case http.MethodPost:
				rt.messageHandler.SendMessage(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			rt.messageHandler.GetConversations(w, r)
		}
	})))

	mux.Handle("/groups", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.groupHandler.GetUserGroups(w, r)
		case http.MethodPost:
			rt.groupHandler.CreateGroup(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/groups/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasSuffix(path, "/join") {
			rt.groupHandler.JoinGroup(w, r)
			return
		}
		if strings.HasSuffix(path, "/leave") {
			rt.groupHandler.LeaveGroup(w, r)
			return
		}
		if strings.HasSuffix(path, "/posts") {
			switch r.Method {
			case http.MethodGet:
				rt.groupHandler.GetGroupPosts(w, r)
			case http.MethodPost:
				rt.groupHandler.PostToGroup(w, r)
			}
			return
		}
		if strings.HasSuffix(path, "/members") {
			rt.groupHandler.GetGroupMembers(w, r)
			return
		}
		if strings.HasSuffix(path, "/settings") {
			rt.groupHandler.UpdateGroupSettings(w, r)
			return
		}

		rt.groupHandler.GetGroup(w, r)
	})))

	mux.Handle("/notifications", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.notifHandler.GetNotifications(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/notifications/clear", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			rt.notifHandler.ClearNotifications(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	mux.Handle("/notifications/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/read") {
			rt.notifHandler.MarkAsRead(w, r)
		}
	})))

	mux.Handle("/reports", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rt.adminHandler.CreateReport(w, r)
		case http.MethodGet:
			rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.GetReports)).ServeHTTP(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/reports/", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.ReviewReport))))

	mux.Handle("/admin/delete/", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.DeleteContent))))
	mux.Handle("/admin/stats", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.GetStats))))
	mux.Handle("/admin/grant", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.GrantAdmin))))
	mux.Handle("/admin/broadcast", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rt.adminHandler.Broadcast(w, r)
		case http.MethodGet:
			rt.adminHandler.GetBroadcasts(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/admin/broadcast/emoji/", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.BroadcastToEmoji))))

	mux.Handle("/upload", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.handleUpload)))

	if rt.uploadDir != "" {
		mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(rt.uploadDir))))
	}

	return rt.rateLimiter.Limit(mux)
}

func (rt *Router) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if rt.uploadDir == "" {
		http.Error(w, `{"error":"uploads not configured"}`, http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, `{"error":"file too large"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"no file provided"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := time.Now().Format("20060102150405") + "_" + header.Filename
	savePath := filepath.Join(rt.uploadDir, filename)

	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[strings.ToLower(ext)] {
		http.Error(w, `{"error":"invalid file type"}`, http.StatusBadRequest)
		return
	}

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, `{"error":"failed to save file"}`, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, `{"error":"failed to save file"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"url":"/uploads/` + filename + `"}`))
}
