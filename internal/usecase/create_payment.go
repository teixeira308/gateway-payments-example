package usecase

import (
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/repository"

	"github.com/google/uuid"
)

type CreatePayment struct {
	Repo repository.PaymentRepository
}

func NewCreatePaymentUseCase(repo repository.PaymentRepository) *CreatePayment {
	return &CreatePayment{
		Repo: repo,
	}
}

func (pc *CreatePayment) Execute(method string, amount float64, orderID string) (*entity.Payment, error) {
	payment := entity.NewPayment(uuid.NewString(), orderID, amount, method)

	err := pc.Repo.Save(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
