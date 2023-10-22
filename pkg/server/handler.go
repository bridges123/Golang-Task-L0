package server

import (
	"encoding/json"
	"net/http"
)

func (s *HTTPServer) HandleGetOrderById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUid := r.URL.Query().Get("uid")
		if len(orderUid) == 0 {
			w.WriteHeader(400)
			w.Write([]byte("Error! Argument 'uid' is empty"))
			return
		}
		order, exists := s.repo.GetById(orderUid)
		if !exists {
			w.WriteHeader(500)
			w.Write([]byte("Sorry. Order " + orderUid + " doesn't exists"))
			return
		}
		orderJson, err := json.Marshal(order)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Sorry. Error with creating response"))
			return
		}
		w.Write(orderJson)
	}
}
