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

// @Summary Process a receipt
// @Description Processes a receipt and returns a unique identifier
// @Tags receipts
// @Accept json
// @Produce json
// @Param receipt body models.Receipt true "Receipt JSON"
// @Success 200 {object} models.ReceiptResult "Successfully processed receipt"
// @Failure 400 {string} string "Invalid receipt data or Failed to process Receipt"
// @Router /receipt/process [post]
func (h *ReceiptHandler) ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	log.Println("ReceiptHandler::ProcessReceipt: started processing the receipt.")
	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		log.Println("ReceiptHandler::ProcessReceipt: error decoding the receipt.")
		http.Error(w, "Invalid receipt data", http.StatusBadRequest)
		return
	}
	// create a id
	id, err := h.service.ProcessReceipt(receipt)
	if err != nil {
		log.Println("ReceiptHandler::ProcessReceipt: error processing the receipt.")
		http.Error(w, "Failed to process Receipt", http.StatusBadRequest)
		return
	}

	response := models.ReceiptResult{Id: id}
	w.Header().Set("Content-Type", "applicaiton/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get points for a receipt id
// @Description Retrieves the points associated with a specific receipt
// @Tags receipts
// @Produce json
// @Param id path string true "Receipt ID"
// @Success 200 {object} models.PointsResult "Successfully retrieved points"
// @Failure 404 {string} string "Receipt not found"
// @Router /receipt/{id}/points [get]
func (h *ReceiptHandler) ReceiptPoints(w http.ResponseWriter, r *http.Request) {
	log.Println("ReceiptHandler::ReceiptPoints: started processing the receipt points.")
	vars := mux.Vars(r)
	id := vars["id"]

	points, err := h.service.ReceiptPoints(id)
	if err != nil {
		log.Printf("ReceiptHandler::ReceiptPoints: error receiving points for the receipt ID %s: %v\n", id, err)
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	response := models.PointsResult{Points: points}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ReceiptHandler::ReceiptPoints: error Encoding the reponse %s: %v\n", id, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
