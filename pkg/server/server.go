package server

import (
	"L0/pkg/storage"
	"html/template"
	"net/http"
)

// HTTPServer - http сервер приложения, отдающий заказы по запросам
type HTTPServer struct {
	server http.Server
	repo   storage.OrderRepo
}

func InitServer(repo storage.OrderRepo) *HTTPServer {
	return &HTTPServer{server: http.Server{Addr: ":8080"}, repo: repo}
}

// Start - привязка всех хендлеров и запуск сервера
func (s *HTTPServer) Start() {
	var tpl = template.Must(template.ParseFiles("templates/index.html"))
	http.HandleFunc("/order", s.HandleGetOrderById())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	})
	s.server.ListenAndServe()
}
