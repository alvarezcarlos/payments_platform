package application

import "github.com/alvarezcarlos/payment/app/domain/entity"

type MerchantUseCaseInterface interface {
	Create(merchant *entity.Merchant) error
}

type PaymentUseCaseInterface interface {
	Create(payment *entity.Payment) error
}
