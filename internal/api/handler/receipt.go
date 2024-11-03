package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"receiptProcessor/internal/models"
	service "receiptProcessor/internal/services"

	"github.com/gorilla/mux"
)

type ReceiptHandler struct {
	service *service.ReceiptService
}

func NewReceiptHandler(service *service.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{service: service}
}

func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		log.Println("error is: ", err)
		http.Error(w, "Invalid receipt data", http.StatusBadRequest)
		return
	}
	// create a id
	id, err := h.service.ProcessReceipt(receipt)
	if err != nil {
		log.Println("error is: ", err)
		http.Error(w, "Failed to process Receipt", http.StatusBadRequest)
		return
	}

	response := models.ReceiptResult{Id: id}
	w.Header().Set("Content-Type", "applicaiton/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ReceiptHandler) ReceiptPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	points, error := h.service.ReceiptPoints(id)
	if error != nil {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := models.PointsResult{Points: points}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
