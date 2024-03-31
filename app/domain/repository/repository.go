package repository

import (
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/google/uuid"
)

type MerchantRepository interface {
	Create(merchant *entity.Merchant) (*entity.Merchant, error)
	GetByName(name string) (*entity.Merchant, error)
}

type PaymentRepository interface {
	Create(payment *entity.Payment) error
	Update(payment *entity.Payment) (*entity.Payment, error)
	GetByID(id uuid.UUID) (*entity.Payment, error)
	CreateCard(card *entity.Card) error
	GetCardByNumber(number string) (*entity.Card, error)
	GetMerchantByID(id uint) (*entity.Merchant, error)
	UpdateCardAndMerchant(card *entity.Card, merchant *entity.Merchant) error
}
