package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/achsanalfitra/gopayslip/internal/app"
	"github.com/achsanalfitra/gopayslip/internal/router"
	"github.com/achsanalfitra/gopayslip/internal/services/empl"
	"github.com/google/uuid"
)

type OvertimeRequest struct {
	Interval     float64 `json:"overtime_duration"`
	OvertimeDate string  `json:"overtime_date"`
}

type ReimbursementRequest struct {
	Amount      float64 `json:"reimbursement_amount"`
	Description string  `json:"description"`
}

type EmplHandler struct {
	EmplService empl.Empl
	UserService empl.User
	App         *app.App
}

func NewEmplHandler(emplSvc empl.Empl, userSvc empl.User, a *app.App) *EmplHandler {
	return &EmplHandler{
		EmplService: emplSvc,
		UserService: userSvc,
		App:         a,
	}
}

func (e *EmplHandler) AttendanceHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(router.CtxUserKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context or invalid type", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format in context", http.StatusInternalServerError)
		return
	}

	requestID, ok := r.Context().Value(router.CtxRequestKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Request ID not found in context or invalid type", http.StatusInternalServerError)
		return
	}

	err = e.UserService.CheckIn(userID, requestID, r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process check-in: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Check-in successful", "request_id": requestID.String()})
}

func (e *EmplHandler) OvertimeHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(router.CtxUserKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context or invalid type", http.StatusInternalServerError)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format in context", http.StatusInternalServerError)
		return
	}

	requestID, ok := r.Context().Value(router.CtxRequestKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Request ID not found in context or invalid type", http.StatusInternalServerError)
		return
	}

	var reqBody OvertimeRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	overtimeDuration := time.Duration(reqBody.Interval * float64(time.Hour))

	overtimeDate, err := time.Parse(time.RFC3339, reqBody.OvertimeDate)
	if err != nil {
		http.Error(w, "Invalid overtime date format. Expected RFC3339 (e.g., 2006-01-02T15:04:05Z07:00)", http.StatusBadRequest)
		return
	}

	err = e.UserService.ProposeOvertime(userID, requestID, overtimeDuration, overtimeDate, r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to propose overtime: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Overtime proposal successful", "request_id": requestID.String()})
}

func (e *EmplHandler) ReimbursementHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(router.CtxUserKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context or invalid type", http.StatusInternalServerError)
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID format in context", http.StatusInternalServerError)
		return
	}

	requestID, ok := r.Context().Value(router.CtxRequestKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Request ID not found in context or invalid type", http.StatusInternalServerError)
		return
	}

	var reqBody ReimbursementRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = e.UserService.ProposeReimbursement(userID, requestID, reqBody.Amount, reqBody.Description, r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to propose reimbursement: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reimbursement proposal successful", "request_id": requestID.String()})
}
