package routes

import (
	"receipt_processor/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/receipts/process", handlers.ProcessReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", handlers.GetPoints).Methods("GET")
}
