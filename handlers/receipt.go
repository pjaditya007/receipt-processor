package handlers

import (
	"encoding/json"
	"net/http"
	"receipt_processor/models"
	"receipt_processor/points"
	"receipt_processor/validation"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	receipts = make(map[string]models.ReceiptData)
	mu       sync.Mutex
)

func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.SanitizeAndValidate(&receipt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	points := points.CalculatePoints(receipt)

	mu.Lock()
	receipts[id] = models.ReceiptData{Receipt: receipt, Points: points}
	mu.Unlock()

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	data, exists := receipts[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := map[string]int{"points": data.Points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
