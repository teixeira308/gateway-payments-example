package usecase

import (
	"context"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/infrastructure/broker"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type CreatePayment struct {
	Repo   repository.PaymentRepository
	Broker *broker.RabbitMQClient
}

func NewCreatePaymentUseCase(repo repository.PaymentRepository, broker *broker.RabbitMQClient) *CreatePayment {
	return &CreatePayment{
		Repo:   repo,
		Broker: broker,
	}
}

func (pc *CreatePayment) Execute(ctx context.Context, paymentRequested event.PaymentRequested) (*entity.Payment, error) {
	// Check idempotency
	existingPayment, err := pc.Repo.FindByOrderID(paymentRequested.OrderID)
	if err != nil && err.Error() != fmt.Sprintf("payment with order ID %s not found", paymentRequested.OrderID) { // Check for actual error not "not found"
		return nil, fmt.Errorf("error checking existing payment for order %s: %w", paymentRequested.OrderID, err)
	}
	if existingPayment != nil {
		fmt.Printf("Payment for order %s already exists, ignoring. Status: %s\n", paymentRequested.OrderID, existingPayment.Status)
		return existingPayment, nil // Idempotent: payment already processed
	}

	// Simulate payment processing (random success/failure)
	rand.Seed(time.Now().UnixNano())
	isApproved := rand.Intn(100) < 80 // 80% chance of approval

	var paymentStatus string
	if isApproved {
		paymentStatus = entity.StatusApproved
	} else {
		paymentStatus = entity.StatusRejected
	}

	// Persist payment record
	payment := entity.NewPayment(uuid.NewString(), paymentRequested.OrderID, paymentRequested.Amount, "Credit Card") // Assuming method
	payment.Status = paymentStatus

	err = pc.Repo.Save(payment)
	if err != nil {
		return nil, fmt.Errorf("error saving payment for order %s: %w", paymentRequested.OrderID, err)
	}

	// Publish payment.processed event
	paymentProcessedEvent := event.PaymentProcessed{
		Event:       "payment.processed",
		OrderID:     payment.OrderID,
		Status:      payment.Status,
		ProcessedAt: time.Now(),
	}

	err = pc.Broker.Publish(ctx, "payments.exchange", "payment.processed", paymentProcessedEvent)
	if err != nil {
		return nil, fmt.Errorf("error publishing payment.processed event for order %s: %w", paymentRequested.OrderID, err)
	}

	return payment, nil
}
