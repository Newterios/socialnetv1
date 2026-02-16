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
	frontendDir    string
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
	frontendDir string,
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
		frontendDir:    frontendDir,
	}
}

func (rt *Router) Setup() http.Handler {
	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/register", rt.authHandler.Register)
	apiMux.HandleFunc("/login", rt.authHandler.Login)
	apiMux.HandleFunc("/auth/google", rt.authHandler.GoogleLogin)
	apiMux.HandleFunc("/auth/password-requirements", rt.authHandler.GetPasswordRequirements)

	apiMux.Handle("/users/search", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.SearchUsers)))
	apiMux.Handle("/users/", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.GetProfile)))

	apiMux.Handle("/profile", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			rt.userHandler.UpdateProfile(w, r)
		case http.MethodDelete:
			rt.userHandler.DeleteAccount(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	apiMux.Handle("/profile/privacy", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.UpdatePrivacySettings)))
	apiMux.Handle("/profile/emoji", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.SetEmojiAvatar)))
	apiMux.Handle("/profile/status", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.UpdateOnlineStatus)))

	apiMux.Handle("/emojis", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.GetEmojis)))
	apiMux.Handle("/emojis/", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.userHandler.GetUsersWithEmoji)))

	apiMux.Handle("/posts", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.postHandler.GetFeed(w, r)
		case http.MethodPost:
			rt.postHandler.CreatePost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	apiMux.Handle("/posts/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	apiMux.Handle("/likes/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rt.socialHandler.LikePost(w, r)
		case http.MethodDelete:
			rt.socialHandler.UnlikePost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	apiMux.Handle("/comments/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.socialHandler.GetComments(w, r)
		case http.MethodPost:
			rt.socialHandler.CommentOnPost(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	apiMux.Handle("/friends", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.socialHandler.GetFriends)))
	apiMux.Handle("/friends/requests", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.socialHandler.GetPendingRequests)))
	apiMux.Handle("/friends/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	apiMux.Handle("/conversations", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.messageHandler.GetConversations(w, r)
		case http.MethodPost:
			rt.messageHandler.StartConversation(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	apiMux.Handle("/conversations/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	apiMux.Handle("/groups", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.groupHandler.GetUserGroups(w, r)
		case http.MethodPost:
			rt.groupHandler.CreateGroup(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	apiMux.Handle("/groups/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	apiMux.Handle("/notifications", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rt.notifHandler.GetNotifications(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	apiMux.Handle("/notifications/clear", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			rt.notifHandler.ClearNotifications(w, r)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	apiMux.Handle("/notifications/", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/read") {
			rt.notifHandler.MarkAsRead(w, r)
		}
	})))

	apiMux.Handle("/reports", rt.authMiddleware.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rt.adminHandler.CreateReport(w, r)
		case http.MethodGet:
			rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.GetReports)).ServeHTTP(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	apiMux.Handle("/reports/", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.ReviewReport))))

	apiMux.Handle("/admin/delete/", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.DeleteContent))))
	apiMux.Handle("/admin/stats", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.GetStats))))
	apiMux.Handle("/admin/grant", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.GrantAdmin))))
	apiMux.Handle("/admin/broadcast", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			rt.adminHandler.Broadcast(w, r)
		case http.MethodGet:
			rt.adminHandler.GetBroadcasts(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	apiMux.Handle("/admin/broadcast/emoji/", rt.authMiddleware.Authenticate(rt.authMiddleware.RequireAdmin(http.HandlerFunc(rt.adminHandler.BroadcastToEmoji))))

	apiMux.Handle("/upload", rt.authMiddleware.Authenticate(http.HandlerFunc(rt.handleUpload)))

	// Top-level mux: /api → API, /uploads → static, everything else → SPA
	topMux := http.NewServeMux()

	// API routes under /api/ prefix
	topMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// Serve uploaded files
	if rt.uploadDir != "" {
		topMux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(rt.uploadDir))))
	}

	// SPA: serve frontend static files, fallback to index.html
	if rt.frontendDir != "" {
		topMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Try to serve the file directly
			filePath := filepath.Join(rt.frontendDir, r.URL.Path)
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				http.ServeFile(w, r, filePath)
				return
			}
			// Fallback to index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(rt.frontendDir, "index.html"))
		})
	}

	return rt.rateLimiter.Limit(topMux)
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
