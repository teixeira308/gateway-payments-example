package http

import (
	"gateway-payments/internal/interface/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	paymentHandler *handler.PaymentHandler,
) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Put("/payments/{id}", paymentHandler.Update)
	router.Get("/payments/{id}", paymentHandler.Get)
	router.Get("/payments", paymentHandler.List)
	router.Delete("/payments/{id}", paymentHandler.Delete)

	return router
}
