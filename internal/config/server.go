package config

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	Addr   string
	Router *http.ServeMux
	Srv    *http.Server
}

const (
	ReadTimeout    = 10 * time.Second
	WriteTimeout   = 10 * time.Second
	MaxHeaderBytes = 1 << 20
)

// custom server building
func CreateServer(addr string, rtr *http.ServeMux) *Server {
	srv := Server{
		Addr:   addr,
		Router: rtr,
	}

	srv.Srv = &http.Server{
		Addr:           srv.Addr,
		Handler:        srv.Router,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		MaxHeaderBytes: MaxHeaderBytes,
	}

	return &srv
}

func (s *Server) Start() {
	log.Printf("running server on port: %s", s.Addr)

	// do error check when the server fails to run
	if err := s.Srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

	log.Print("server is running")
}
