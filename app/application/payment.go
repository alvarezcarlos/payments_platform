package application

import (
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"log/slog"
	"strconv"
	"time"

	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"github.com/alvarezcarlos/payment/app/utils"
	"github.com/google/uuid"
)

const (
	paymentConst         = "payment"
	refundConst          = "refund"
	errorProcessing      = "error processing %s, please try again later"
	errorInvalidState    = "error invalid payment state for processing"
	errorCreatingPayment = "error creating payment"
)

type paymentUseCase struct {
	repository repository.PaymentRepository
	logger     *slog.Logger
}

func NewPaymentUseCase(paymentRepository repository.PaymentRepository, logger *slog.Logger) PaymentUseCaseInterface {
	return &paymentUseCase{repository: paymentRepository, logger: logger}
}

// Create payment can only be accessed by a Merchant, that will partially populate it with fields like
// Amount and other merchant information and the Customer should be redirected with the payment_id for processing.
func (p *paymentUseCase) Create(payment *entity.Payment) (*entity.Payment, error) {
	payment.ID = uuid.New()
	payment.States = append(payment.States, entity.SetState(entity.Pending))
	payment.CreatedAt, payment.UpdatedAt = time.Now(), time.Now()
	if err := p.repository.Create(payment); err != nil {
		p.logger.Error(err.Error())
		return nil, errors.New(errorCreatingPayment)
	}
	p.logger.Info("payment created ", payment)
	return payment, nil
}

// GetByID retrieve a payment details by Id
func (p *paymentUseCase) GetByID(uuid uuid.UUID) (*entity.Payment, error) {
	payment, err := p.repository.GetByID(uuid)
	if err != nil {
		p.logger.Error(err.Error())
		return nil, errors.New("error fetching payment")
	}
	return payment, nil
}

// ProcessPayment the business core functionality, it allows the customer to complete the payment,
// also stores the information of the card assigning random founds to it and perform the transaction with the bank
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

	//only pay pending or rejected operations
	statesMap := statesToMap(pay.States)
	if _, isSucceeded := statesMap[string(entity.Succeeded)]; isSucceeded {
		return nil, fmt.Errorf(errorInvalidState)
	}

	p.logger.Info("initializing payment process with mock Bank")
	err = p.BankSimulator(pay, paymentConst)
	if err != nil {
		return nil, errors.New("from bank" + err.Error())
	}

	updatedPayment, err := p.repository.Update(pay)
	if err != nil {
		return nil, fmt.Errorf(errorProcessing, paymentConst)
	}
	p.logger.Info("payment processed successfully")
	return updatedPayment, nil
}

// ProcessRefund payment can only be accessed by a Merchant, to execute the devolution for client money
func (p *paymentUseCase) ProcessRefund(uuid uuid.UUID, merchantId string) error {
	pay, err := p.repository.GetByID(uuid)
	if err != nil {
		p.logger.Error(err.Error())
		return fmt.Errorf(errorProcessing, refundConst)
	}

	if strconv.Itoa(int(pay.MerchantID)) != merchantId {
		log.Error("merchants don't match")
		return fmt.Errorf(errorProcessing, refundConst)
	}

	p.logger.Info("payment to be refunded", pay)

	//only refund a successful operation
	statesMap := statesToMap(pay.States)
	_, isRefunded := statesMap[string(entity.Refunded)]
	_, isSucceeded := statesMap[string(entity.Succeeded)]

	if isRefunded || !isSucceeded {
		return errors.New(fmt.Sprintf("%s ", errorInvalidState))
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

// BankSimulator simulates the communication with the Bank
// It fails if it isn't enough found for an operation
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
	p.logger.Info("processing operation in mock Bank")
	err = p.repository.UpdateCardAndMerchant(card, merch)
	if err != nil {
		p.logger.Error(err.Error())
		return errors.New("error from mock bank api")
	}

	return nil
}

func statesToMap(states []entity.State) map[string]uint {
	statesMap := make(map[string]uint)
	for _, state := range states {
		statesMap[state.Name] = state.ID
	}
	return statesMap
}
