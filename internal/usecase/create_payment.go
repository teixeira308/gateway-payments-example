package usecase

import (
	"context"
	"fmt"
	"gateway-payments/internal/domain/entity"
	"gateway-payments/internal/domain/event"
	"gateway-payments/internal/domain/repository"
	"gateway-payments/internal/infrastructure/broker"
	"math/rand"
	"os"
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

	autoApprove := os.Getenv("AUTO_APPROVE_PAYMENTS") == "true"

	var paymentStatus string
	if autoApprove {
		// Lógica atual: Simular processamento (random success/failure)
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) < 80 {
			paymentStatus = entity.StatusApproved
		} else {
			paymentStatus = entity.StatusRejected
		}
	} else {
		// Nova lógica: Nasce pendente para aprovação manual via PUT
		paymentStatus = "PENDING"
	}
	// Persist payment record
	payment := entity.NewPayment(uuid.NewString(), paymentRequested.OrderID, paymentRequested.Amount, "Credit Card")
	payment.Status = paymentStatus

	err = pc.Repo.Save(payment)
	if err != nil {
		return nil, fmt.Errorf("error saving payment: %w", err)
	}

	// 2. Só dispara o evento se o pagamento já estiver decidido (Approved ou Rejected)
	if paymentStatus != "PENDING" {
		paymentProcessedEvent := event.PaymentProcessed{
			Event:       "payment.processed",
			OrderID:     payment.OrderID,
			Status:      payment.Status,
			ProcessedAt: time.Now(),
		}

		err = pc.Broker.Publish(ctx, "payments.exchange", "payment.processed", paymentProcessedEvent)
		if err != nil {
			return nil, fmt.Errorf("error publishing event: %w", err)
		}
	}

	return payment, nil
}
