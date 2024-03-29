package models

type Merchant struct {
	Name    string `json:"name"`
	Account int32  `json:"account"`
}

type PaymentCreateReq struct {
	MerchantID uint    `json:"merchant_id"`
	Amount     float64 `json:"amount"`
}
