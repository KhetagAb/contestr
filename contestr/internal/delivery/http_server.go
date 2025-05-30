package delivery

import (
	"log"
	"net/http"
)

type HTTPServer struct {
	httpServer *http.Server
}

func NewHTTPServer() *HTTPServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	return &HTTPServer{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}

func (s *HTTPServer) Start() error {
	log.Printf("HTTPServer starting on port 8080")
	return s.httpServer.ListenAndServe()
}

func (s *HTTPServer) Stop() error {
	return s.httpServer.Close()
}
