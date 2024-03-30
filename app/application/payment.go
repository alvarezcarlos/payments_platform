package application

import (
	"errors"
	"fmt"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"github.com/google/uuid"
	"log/slog"
	"math/rand"
	"time"
)

const (
	errorProcessingPayment = "error processing payment please try again later"
	errorInvalidState      = "error invalid payment state for processing"
)

type paymentUseCase struct {
	repository repository.PaymentRepository
	logger     *slog.Logger
}

func NewPaymentUseCase(paymentRepository repository.PaymentRepository, logger *slog.Logger) PaymentUseCaseInterface {
	return &paymentUseCase{repository: paymentRepository, logger: logger}
}

func (p *paymentUseCase) Create(payment *entity.Payment) (*entity.Payment, error) {
	payment.ID = uuid.New()
	payment.States = append(payment.States, entity.SetState(entity.Pending))
	payment.CreatedAt, payment.UpdatedAt = time.Now(), time.Now()
	if err := p.repository.Create(payment); err != nil {
		p.logger.Error(err.Error())
		return nil, errors.New("error creating payment")
	}
	return payment, nil
}

func (p *paymentUseCase) GetByID(uuid uuid.UUID) (*entity.Payment, error) {
	payment, err := p.repository.GetByID(uuid)
	if err != nil {
		p.logger.Error(err.Error())
		return nil, errors.New("error fetching payment")
	}
	return payment, nil
}

func (p *paymentUseCase) ProcessPayment(
	payment *entity.Payment,
	card *entity.Card) (*entity.Payment, error) {
	err := p.repository.CreateCard(card)
	if err != nil {
		p.logger.Error(err.Error())
		return nil, errors.New(errorProcessingPayment)
	}

	pay, err := p.repository.GetByID(payment.ID)
	if err != nil {
		p.logger.Error(err.Error())
		return nil, errors.New(errorProcessingPayment)
	}

	pay.CardNumber = card.Number

	for _, state := range pay.States {
		if state == entity.SetState(entity.Succeeded) ||
			state == entity.SetState(entity.Refunded) {
			return nil, errors.New(fmt.Sprintf("%s (%s)", errorInvalidState, state.Name))
		}
	}

	p.BankSimulator(pay)

	updatedPayment, err := p.repository.Update(pay)
	if err != nil {
		return nil, errors.New(errorProcessingPayment)
	}

	return updatedPayment, nil
}

func (p *paymentUseCase) BankSimulator(payment *entity.Payment) {
	randomInt := rand.Intn(2)
	if randomInt == 0 {
		payment.States = append(payment.States, entity.SetState(entity.Rejected))
	} else {
		payment.States = append(payment.States, entity.SetState(entity.Succeeded))
	}
}
