package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/achsanalfitra/gopayslip/internal/app"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"user_role"`
}

type LoginResponse struct {
	// Message string `json:"message"` // for testing only <----
	Access  string `json:"access"` // use it when tokenizer is already online
	Refresh string `json:"refresh"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type AuthHandler struct {
	AuthService AuthService
	Tokenizer   *Tokenizer
	App         *app.App
}

func NewAuthHandler(svc AuthService, a *app.App) *AuthHandler {
	return &AuthHandler{
		AuthService: svc,
		Tokenizer:   NewTokenizer(),
		App:         a,
	}
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// decode body to json
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// inject DB
	newCtx := context.WithValue(r.Context(), app.PQ, ah.App.DB)

	// update context with injected DB
	r = r.WithContext(newCtx)

	// run login service
	err := ah.AuthService.Login(req.Username, req.Password, req.Role, r.Context())
	if err != nil {
		// check for unauthorized
		if errors.Is(err, errors.New("user not found")) || errors.Is(err, errors.New("invalid password")) {
			log.Printf("internal server error during login for user %s: %v", req.Username, err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid username/password"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "an unexpected error occured"})
		return
	}

	// get token
	access, refresh, err := ah.Tokenizer.GenerateToken(req.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to generate token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Access: access, Refresh: refresh})
}

func (ah *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// decode body to json
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
		return
	}

	// inject DB
	newCtx := context.WithValue(r.Context(), app.PQ, ah.App.DB)

	// update context with injected DB
	r = r.WithContext(newCtx)

	// run login service
	err := ah.AuthService.Login(req.Username, req.Password, req.Role, r.Context())
	if err != nil {
		// check for unauthorized
		if errors.Is(err, errors.New("user not found")) || errors.Is(err, errors.New("invalid password")) {
			log.Printf("internal server error during login for user %s: %v", req.Username, err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid username/password"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "an unexpected error occured"})
		return
	}

	// get token
	access, refresh, err := ah.Tokenizer.GenerateToken(req.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to generate token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Access: access, Refresh: refresh})
}
