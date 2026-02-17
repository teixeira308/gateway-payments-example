package http

import (
	"net/http"

	"gateway-payments/internal/interface/http/handler"
	"gateway-payments/internal/usecase"
)

func NewRouter(createPayment *usecase.CreatePayment) http.Handler {

	paymentHandler := &handler.PaymentHandler{
		CreatePayment: createPayment,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /payments", paymentHandler.Create)

	return mux
}
