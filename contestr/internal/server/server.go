package server

import (
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("Server starting on port 8080")
	return s.httpServer.ListenAndServe()
}
