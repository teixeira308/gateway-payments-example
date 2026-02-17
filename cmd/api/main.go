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
	"gateway-payments/internal/usecase"
)

func main() {

	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		log.Fatal(err)
	}

	paymentRepo := mysqlRepo.NewPaymentRepository(db)

	createPayment := &usecase.CreatePayment{Repo: paymentRepo}

	router := httpRouter.NewRouter(createPayment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
	log.Println("Server running on :8080")
}
