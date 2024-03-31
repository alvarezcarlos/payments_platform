package application

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"github.com/alvarezcarlos/payment/app/utils"
	"github.com/google/uuid"
)

const (
	paymentConst      = "payment"
	refundConst       = "refund"
	errorProcessing   = "error processing %s, please try again later"
	errorInvalidState = "error invalid payment state for processing"
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
	card.Balance = utils.RandomFloat()
	err := p.repository.CreateCard(card)
	if err != nil {
		p.logger.Error(err.Error())
		return nil, fmt.Errorf(errorProcessing, paymentConst)
	}

	pay, err := p.repository.GetByID(payment.ID)
	if err != nil {
		p.logger.Error(err.Error())
		return nil, fmt.Errorf(errorProcessing, paymentConst)
	}

	pay.CardNumber = card.Number

	for _, state := range pay.States {
		if state == entity.SetState(entity.Succeeded) ||
			state == entity.SetState(entity.Refunded) {
			return nil, errors.New(fmt.Sprintf("%s (%s)", errorInvalidState, pay.States[len(pay.States)-1].Name))
		}
	}

	err = p.BankSimulator(pay, paymentConst)
	if err != nil {
		return nil, errors.New("from bank" + err.Error())
	}

	updatedPayment, err := p.repository.Update(pay)
	if err != nil {
		return nil, fmt.Errorf(errorProcessing, paymentConst)
	}

	return updatedPayment, nil
}

func (p *paymentUseCase) ProcessRefund(uuid uuid.UUID) error {
	pay, err := p.repository.GetByID(uuid)
	if err != nil {
		p.logger.Error(err.Error())
		return fmt.Errorf(errorProcessing, refundConst)
	}

	if pay.States[len(pay.States)-1] != entity.SetState(entity.Succeeded) {
		return errors.New(fmt.Sprintf("%s (%s)", errorInvalidState, pay.States[len(pay.States)-1].Name))
	}

	err = p.BankSimulator(pay, refundConst)
	if err != nil {
		p.logger.Error(err.Error())
		return fmt.Errorf(errorProcessing, refundConst)
	}

	_, err = p.repository.Update(pay)
	if err != nil {
		return fmt.Errorf(errorProcessing, refundConst)
	}

	return nil
}

func (p *paymentUseCase) BankSimulator(payment *entity.Payment, op string) error {
	card, err := p.repository.GetCardByNumber(payment.CardNumber)
	if err != nil {
		return err
	}

	merch, err := p.repository.GetMerchantByID(payment.MerchantID)
	if err != nil {
		return err
	}

	if op == paymentConst {
		card.Balance = card.Balance - payment.Amount
		merch.Balance = merch.Balance + payment.Amount
		if card.Balance < 0 {
			p.logger.Error("insufficient founds in card balance")
			payment.States = append(payment.States, entity.SetState(entity.Rejected))
			return nil
		}
		payment.States = append(payment.States, entity.SetState(entity.Succeeded))
	} else {
		card.Balance = card.Balance + payment.Amount
		merch.Balance = merch.Balance - payment.Amount
		if merch.Balance < 0 {
			p.logger.Error("insufficient founds in merchant balance")
			return fmt.Errorf(errorProcessing, refundConst)
		}
		payment.States = append(payment.States, entity.SetState(entity.Refunded))
	}

	err = p.repository.UpdateCardAndMerchant(card, merch)
	if err != nil {
		p.logger.Error(err.Error())
		return errors.New("error from mock bank api")
	}
	//randomInt := rand.Intn(5)
	//if randomInt == 0 {
	//	payment.States = append(payment.States, entity.SetState(entity.Rejected))
	//} else {
	//	payment.States = append(payment.States, entity.SetState(entity.Succeeded))
	//}

	return nil
}
