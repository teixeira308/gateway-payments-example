package repository

import "gateway-payments/internal/domain/entity"

type PaymentRepository interface {
	Save(payment *entity.Payment) error
}
