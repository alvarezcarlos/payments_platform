package models

type Merchant struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PaymentCreateReq struct {
	MerchantID uint    `json:"merchant_id" validate:"required,gt=0"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

type ProcessPaymentReq struct {
	PaymentID string   `json:"payment_id" validate:"required,uuid"`
	Card      Card     `json:"card" validate:"required"`
	Customer  Customer `json:"customer" validate:"required"`
}

type Card struct {
	Number string `json:"number" validate:"required,min=16,max=16"`
	Code   string `json:"code" validate:"required,min=3,max=3"`
	Month  int    `json:"month" validate:"required,min=1,max=12"`
	Year   int    `json:"year" validate:"required,min=0,max=99"`
}
type Customer struct {
	PersonalID uint   `json:"personal_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

type RefundPaymentReq struct {
	PaymentID string `json:"payment_id" validate:"required,uuid"`
}
