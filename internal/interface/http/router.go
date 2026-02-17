package http

import (
	"gateway-payments/internal/interface/http/handler"
	"gateway-payments/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	createPayment *usecase.CreatePayment,
	updatePayment *usecase.UpdatePayment,
	getPayment *usecase.GetPayment,
	getAllPayments *usecase.GetAllPayments,
	deletePayment *usecase.DeletePayment,
) *chi.Mux {
	paymentHandler := &handler.PaymentHandler{
		CreatePayment:  createPayment,
		UpdatePayment:  updatePayment,
		GetPayment:     getPayment,
		GetAllPayments: getAllPayments,
		DeletePayment:  deletePayment,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/payments", paymentHandler.Create)
	router.Put("/payments/{id}", paymentHandler.Update)
	router.Get("/payments/{id}", paymentHandler.Get)
	router.Get("/payments", paymentHandler.List)
	router.Delete("/payments/{id}", paymentHandler.Delete)

	return router
}
