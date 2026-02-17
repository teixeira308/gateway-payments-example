package usecase

import (
	"errors"
	"gateway-payments/internal/domain/repository"
)

type UpdatePaymentInput struct {
	ID     string
	Status string
}

type UpdatePayment struct {
	Repo repository.PaymentRepository
}

func NewUpdatePaymentUseCase(repo repository.PaymentRepository) *UpdatePayment {
	return &UpdatePayment{
		Repo: repo,
	}
}

func (up *UpdatePayment) Execute(input UpdatePaymentInput) error {
	payment, err := up.Repo.FindByID(input.ID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Validate status transition if necessary, otherwise just set
	// For this example, we'll assume any status can be set directly.
	// In a real application, you'd add logic like:
	// if !payment.CanTransitionTo(input.Status) {
	//    return errors.New("invalid status transition")
	// }

	payment.Status = input.Status

	err = up.Repo.Save(payment) // Assuming Save also handles updates
	if err != nil {
		return err
	}

	return nil
}
