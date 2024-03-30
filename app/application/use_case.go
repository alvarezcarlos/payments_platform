package application

import (
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/google/uuid"
)

type MerchantUseCaseInterface interface {
	Create(merchant *entity.Merchant) error
}

type PaymentUseCaseInterface interface {
	Create(payment *entity.Payment) (*entity.Payment, error)
	GetByID(uuid uuid.UUID) (*entity.Payment, error)
	ProcessPayment(
		payment *entity.Payment,
		customer *entity.Card) (*entity.Payment, error)
}
