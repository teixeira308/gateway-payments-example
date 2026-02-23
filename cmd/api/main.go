package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"gateway-payments/internal/infrastructure/config"
	mysqlRepo "gateway-payments/internal/infrastructure/database/mysql"
	httpRouter "gateway-payments/internal/interface/http"
	httpHandler "gateway-payments/internal/interface/http/handler"
	"gateway-payments/internal/usecase"
)

func main() {

	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		log.Fatal(err)
	}

	paymentRepo := mysqlRepo.NewPaymentRepository(db)

	createPayment := usecase.NewCreatePaymentUseCase(paymentRepo)
	updatePayment := usecase.NewUpdatePaymentUseCase(paymentRepo)
	getPayment := usecase.NewGetPaymentUseCase(paymentRepo)
	getAllPayments := usecase.NewGetAllPaymentsUseCase(paymentRepo)
	deletePayment := usecase.NewDeletePaymentUseCase(paymentRepo)

	paymentHandler := httpHandler.NewPaymentHandler(
		createPayment,
		updatePayment,
		getPayment,
		getAllPayments,
		deletePayment,
	)

	router := httpRouter.NewRouter(
		paymentHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
