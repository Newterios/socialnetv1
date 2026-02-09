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

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var create model.ReportCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := h.adminService.CreateReport(userID, &create); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"report created"}`))
}

func (h *AdminHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	statusStr := r.URL.Query().Get("status")
	if statusStr == "" {
		statusStr = string(model.ReportStatusPending)
	}

	status := model.ReportStatus(statusStr)

	reports, err := h.adminService.GetReports(status)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}

func (h *AdminHandler) ReviewReport(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, `{"error":"invalid report ID"}`, http.StatusBadRequest)
		return
	}

	reportID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid report ID"}`, http.StatusBadRequest)
		return
	}

	var req struct {
		Status model.ReportStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := h.adminService.ReviewReport(reportID, req.Status); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"report reviewed"}`))
}

func (h *AdminHandler) DeleteContent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	targetType := model.ReportTargetType(parts[3])
	targetID, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		http.Error(w, `{"error":"invalid target ID"}`, http.StatusBadRequest)
		return
	}

	if err := h.adminService.DeleteContent(targetType, targetID); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"content deleted"}`))
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.adminService.GetSiteStats()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *AdminHandler) GrantAdmin(w http.ResponseWriter, r *http.Request) {
	var req model.GrantAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := h.adminService.GrantAdmin(req.UserID); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"admin rights granted"}`))
}

func (h *AdminHandler) Broadcast(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r)

	var req model.BroadcastCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := h.adminService.BroadcastNotification(adminID, req.Message); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"broadcast sent"}`))
}

func (h *AdminHandler) GetBroadcasts(w http.ResponseWriter, r *http.Request) {
	broadcasts, err := h.adminService.GetBroadcasts(20)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(broadcasts)
}

func (h *AdminHandler) BroadcastToEmoji(w http.ResponseWriter, r *http.Request) {
	adminID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, `{"error":"emoji ID required"}`, http.StatusBadRequest)
		return
	}

	emojiID := parts[4]

	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := h.adminService.BroadcastToEmoji(adminID, emojiID, req.Message); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"broadcast sent to emoji users"}`))
}
