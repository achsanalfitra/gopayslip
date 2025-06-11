package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/achsanalfitra/gopayslip/internal/app"
	"github.com/achsanalfitra/gopayslip/internal/auth"
	"github.com/google/uuid"
)

// configure public path
var publicPath = map[string]bool{
	"/api/": true,
}

type ReqKey string
type UserKey string
type StartKey string
type Endkey string

const (
	CtxRequestKey ReqKey   = "requestkey"
	CtxUserKey    UserKey  = "userkey"
	CtxStartKey   StartKey = "startdate"
	CtxEndKey     Endkey   = "enddate"
)

type Router struct {
	Route     map[string]map[string]http.HandlerFunc // format -> path: {method: http.HandlerFunc}
	Tokenizer *auth.Tokenizer
	auth      *auth.AuthHandler
	a         *app.App
	mu        sync.RWMutex
}

type defaultAuthServiceStruct struct {
	auth.AuthService
}

func NewDefaultAuthService() auth.AuthService {
	return &defaultAuthServiceStruct{}
}

func NewRouter(a *app.App) *Router {
	router := Router{
		Route:     make(map[string]map[string]http.HandlerFunc),
		Tokenizer: auth.NewTokenizer(),
		a:         a,
		mu:        sync.RWMutex{},
	}

	// assign inherent auth functionality
	authSvc := auth.NewAuthService()
	router.auth = auth.NewAuthHandler(a, authSvc)

	return &router
}

func (r *Router) RegisterRoute(method, path string, handler http.HandlerFunc) error {
	// instantiate the path if it doesn't exist
	if _, exists := r.Route[path]; !exists {
		r.Route[path] = make(map[string]http.HandlerFunc)
	}

	// check pattern existence
	if _, exists := r.Route[path][method]; exists {
		return fmt.Errorf("this path %s with %s method already exists", path, method)
	}

	r.Route[path][method] = handler

	return nil
}

// boilerplate entry point for accessing the handler function
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	// check path existence
	if _, exists := r.Route[path]; !exists {
		http.NotFound(w, req)
		log.Printf("path %s not found", path)
		return
	}

	// check method existence
	if _, exists := r.Route[path][method]; !exists {
		http.NotFound(w, req)
		log.Printf("method %s does not exist on path %s", method, path)
		return
	}

	if !publicPath[path] {
		// parse header, look for Authorization
		access, err := r.Tokenizer.ReadToken(req)
		if err != nil {
			http.Error(w, "bad authorization header", http.StatusUnauthorized)
			return
		}
		if err := r.Tokenizer.AuthorizeToken(access); err != nil {
			http.Error(w, "token unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := r.Tokenizer.GetUserFromAccess(access)
		if err != nil {
			http.Error(w, "token unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := r.auth.UserIDFromToken(user)
		if err != nil {
			http.Error(w, "token unauthorized", http.StatusUnauthorized)
			return
		}

		// Context injection
		currentCtx := req.Context()

		// Lock this to prevent panic
		r.mu.RLock()
		currentCtx = context.WithValue(currentCtx, CtxStartKey, r.a.InitStates[string(CtxStartKey)])
		currentCtx = context.WithValue(currentCtx, CtxEndKey, r.a.InitStates[string(CtxEndKey)])
		r.mu.Unlock()

		newRequestId := uuid.New()
		currentCtx = context.WithValue(currentCtx, CtxRequestKey, newRequestId)
		currentCtx = context.WithValue(currentCtx, CtxUserKey, userID)

		// Update the request with the fully populated context
		req = req.WithContext(currentCtx)

		r.Route[path][method](w, req)
	}

	r.Route[path][method](w, req)
}
