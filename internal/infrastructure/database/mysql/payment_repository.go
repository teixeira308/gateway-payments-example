package mysql

import (
	"database/sql"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"time"
)

type PaymentRepository struct {
	DB *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

func (r *PaymentRepository) Save(payment *entity.Payment) error {
	// Definimos o tempo no Go antes de salvar
	// Assim temos o valor exato para retornar na API
	payment.CreatedAt = time.Now()

	// Se o status estiver vazio na entidade, garantimos o PENDING
	if payment.Status == "" {
		payment.Status = entity.StatusPending
	}

	query := `INSERT INTO payments (id, method, amount, status, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(
		query,
		payment.ID,
		payment.Method,
		payment.Amount,
		payment.Status,
		payment.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error persisting payment [%s]: %w", payment.ID, err)
	}

	return nil
}
