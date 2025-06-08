package internal

import (
	"fmt"
	"log"
	"net/http"
)

type Router struct {
	Route map[string]map[string]http.HandlerFunc // format -> path: {method: http.HandlerFunc}
}

func NewRouter() *Router {
	return &Router{
		Route: make(map[string]map[string]http.HandlerFunc),
	}
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

	r.Route[path][method](w, req)
}
