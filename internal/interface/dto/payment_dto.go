package dto

import (
	"gateway-payments/internal/domain/entity"
	"time"
)

type CreatePaymentRequest struct {
	Amount float64 `json:"amount"`
	Method string  `json:"method"`
}

type UpdatePaymentRequest struct {
	Status string `json:"status"`
}

type PaymentResponse struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatePaymentResponse(payment *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:        payment.ID,
		Method:    payment.Method,
		Amount:    payment.Amount,
		CreatedAt: payment.CreatedAt,
	}
}
