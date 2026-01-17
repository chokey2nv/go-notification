package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/chokey2nv/go-notification/service"
)

type Handler struct {
	svc *service.NotificationService
}

func NewHandler(svc *service.NotificationService) *Handler {
	return &Handler{svc: svc}
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	idemKey := r.Header.Get("Idempotency-Key")

	n, err := h.svc.Add(r.Context(), service.AddInput{
		IdempotencyKey: idemKey,
		UserID:         req.UserID,
		Title:          req.Title,
		Message:        req.Message,
		TemplateID:     req.TemplateID,
		Channels:       req.Channels,
		Metadata:       req.Metadata,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		Channels:  n.Channels,
		Metadata:  n.Metadata,
		CreatedAt: n.CreatedAt.Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	n, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Title:     n.Title,
		Message:   n.Message,
		Channels:  n.Channels,
		Metadata:  n.Metadata,
		CreatedAt: n.CreatedAt.Unix(),
	}

	_ = json.NewEncoder(w).Encode(resp)
}
func (h *Handler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")

	items, err := h.svc.GetByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]NotificationResponse, 0, len(items))
	for _, n := range items {
		resp = append(resp, NotificationResponse{
			ID:        n.ID,
			UserID:    n.UserID,
			Title:     n.Title,
			Message:   n.Message,
			Channels:  n.Channels,
			Metadata:  n.Metadata,
			CreatedAt: n.CreatedAt.Unix(),
		})
	}

	_ = json.NewEncoder(w).Encode(resp)
}
