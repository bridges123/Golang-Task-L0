package server

import (
	"L0/pkg/storage"
	"net/http"
)

type HTTPServer struct {
	server http.Server
	repo   storage.OrderRepo
}

func InitServer(repo storage.OrderRepo) *HTTPServer {
	return &HTTPServer{server: http.Server{Addr: ":8080"}, repo: repo}
}

func (s *HTTPServer) Start() {
	http.HandleFunc("/order", s.HandleGetOrderById())
	s.server.ListenAndServe()
}
