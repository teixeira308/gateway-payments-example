package usecase

import (
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/repository"

	"github.com/google/uuid"
)

type CreatePayment struct {
	Repo repository.PaymentRepository
}

func (pc *CreatePayment) Execute(method string, amount float64) (*entity.Payment, error) {
	payment := &entity.Payment{
		ID:     uuid.NewString(),
		Amount: amount,
		Method: method,
	}

	err := pc.Repo.Save(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
