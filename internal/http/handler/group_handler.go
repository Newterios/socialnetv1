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

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var create model.GroupCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	group, err := h.groupService.CreateGroup(userID, &create)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	group, err := h.groupService.GetGroup(groupID, userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

func (h *GroupHandler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	groups, err := h.groupService.GetUserGroups(userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.groupService.JoinGroup(groupID, userID); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"joined group"}`))
}

func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.groupService.LeaveGroup(groupID, userID); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"left group"}`))
}

func (h *GroupHandler) PostToGroup(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	var create model.GroupPostCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	post, err := h.groupService.PostToGroup(groupID, userID, &create)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *GroupHandler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	posts, err := h.groupService.GetGroupPosts(groupID, userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	members, err := h.groupService.GetGroupMembers(groupID, userID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func (h *GroupHandler) UpdateGroupSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	groupID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid group ID"}`, http.StatusBadRequest)
		return
	}

	var settings model.GroupSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := h.groupService.UpdateGroupSettings(groupID, userID, &settings); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"group settings updated"}`))
}
