package application

import (
	"errors"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type paymentUseCase struct {
	repository repository.PaymentRepository
	logger     *slog.Logger
}

func NewPaymentUseCase(paymentRepository repository.PaymentRepository, logger *slog.Logger) PaymentUseCaseInterface {
	return &paymentUseCase{repository: paymentRepository, logger: logger}
}

func (p *paymentUseCase) Create(payment *entity.Payment) error {
	payment.ID = uuid.New()
	payment.States = append(payment.States, entity.SetState("Pending"))
	payment.CreatedAt, payment.UpdatedAt = time.Now(), time.Now()
	if err := p.repository.Create(payment); err != nil {
		p.logger.Error(err.Error())
		return errors.New("error creating payment")
	}
	return nil
}
