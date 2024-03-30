package repository

import (
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/google/uuid"
)

type MerchantRepository interface {
	Create(merchant *entity.Merchant) error
}

type PaymentRepository interface {
	Create(payment *entity.Payment) error
	Update(payment *entity.Payment) (*entity.Payment, error)
	GetByID(id uuid.UUID) (*entity.Payment, error)
	CreateCard(card *entity.Card) error
}
