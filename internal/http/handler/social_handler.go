package handler

import (
	"encoding/json"
	"net/http"
	"socialnet/internal/http/middleware"
	"socialnet/internal/model"
	"socialnet/internal/service"
	"strconv"
	"strings"
)

type SocialHandler struct {
	socialService *service.SocialService
}

func NewSocialHandler(socialService *service.SocialService) *SocialHandler {
	return &SocialHandler{socialService: socialService}
}

func (h *SocialHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req model.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.socialService.SendFriendRequest(userID, req.AddresseeID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"friend request sent"}`))
}

func (h *SocialHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid request ID", http.StatusBadRequest)
		return
	}

	requestID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid request ID", http.StatusBadRequest)
		return
	}

	if err := h.socialService.AcceptFriendRequest(requestID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"friend request accepted"}`))
}

func (h *SocialHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid request ID", http.StatusBadRequest)
		return
	}

	requestID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid request ID", http.StatusBadRequest)
		return
	}

	if err := h.socialService.BlockUser(requestID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"user blocked"}`))
}

func (h *SocialHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	friends, err := h.socialService.GetFriends(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func (h *SocialHandler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	requests, err := h.socialService.GetPendingRequests(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func (h *SocialHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	if err := h.socialService.LikePost(postID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"post liked"}`))
}

func (h *SocialHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	if err := h.socialService.UnlikePost(postID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"post unliked"}`))
}

func (h *SocialHandler) CommentOnPost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	var create model.CommentCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	comment, err := h.socialService.CommentOnPost(postID, userID, &create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *SocialHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := h.socialService.GetComments(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
