package config

import (
	"net/http"
)

type Server struct {
	Addr   string
	Router *http.ServeMux
	Srv    *http.Server
}

// custom server building
func CreateServer(addr string, rtr *http.ServeMux) *Server {
	srv := Server{
		Addr:   addr,
		Router: rtr,
	}

	srv.Srv = &http.Server{
		Addr:    srv.Addr,
		Handler: srv.Router,
	}

	return &srv
}
