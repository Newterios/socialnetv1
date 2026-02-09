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

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) StartConversation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var create model.ConversationCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	conversation, err := h.messageService.StartConversation(userID, create.ParticipantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(conversation)
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	conversationID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	var create model.MessageCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	message, err := h.messageService.SendMessage(conversationID, userID, &create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	conversationID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	messages, err := h.messageService.GetMessages(conversationID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *MessageHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	conversations, err := h.messageService.GetConversations(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}
