package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway-payments/internal/interface/dto"
	"gateway-payments/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	CreatePayment  *usecase.CreatePayment
	UpdatePayment  *usecase.UpdatePayment
	GetPayment     *usecase.GetPayment
	GetAllPayments *usecase.GetAllPayments
	DeletePayment  *usecase.DeletePayment
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}

	var input dto.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usecaseInput := usecase.UpdatePaymentInput{
		ID:     paymentID,
		Status: input.Status,
	}

	err := h.UpdatePayment.Execute(usecaseInput)
	if err != nil {
		if errors.Is(err, errors.New("payment not found")) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PaymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}

	usecaseInput := usecase.GetPaymentInput{
		ID: paymentID,
	}

	payment, err := h.GetPayment.Execute(usecaseInput)
	if err != nil {
		// It's better to check the specific error returned by the use case (e.g., from repository)
		// For now, a generic "payment not found" check is used.
		if err.Error() == fmt.Sprintf("payment with ID %s not found", paymentID) { // Specific error check
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	usecaseInput := usecase.GetAllPaymentsInput{
		Page:  page,
		Limit: limit,
	}

	paymentsOutput, err := h.GetAllPayments.Execute(usecaseInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paymentsOutput.Payments)
}

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	paymentID := chi.URLParam(r, "id")
	if paymentID == "" {
		http.Error(w, "payment ID is required", http.StatusBadRequest)
		return
	}

	usecaseInput := usecase.DeletePaymentInput{
		ID: paymentID,
	}

	err := h.DeletePayment.Execute(usecaseInput)
	if err != nil {
		if err.Error() == fmt.Sprintf("payment with ID %s not found for deletion", paymentID) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
