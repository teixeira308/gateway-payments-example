package handler

import (
	"encoding/json"
	"gateway-payments/internal/interface/dto"
	"gateway-payments/internal/usecase"
	"net/http"
)

type PaymentHandler struct {
	CreatePayment *usecase.CreatePayment
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	payment, err := h.CreatePayment.Execute(input.Method, input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := dto.CreatePaymentResponse(payment)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
