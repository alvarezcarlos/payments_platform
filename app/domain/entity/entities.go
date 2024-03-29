package entity

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Name      string
	Card      Card
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Card struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	CustomerID uint
	Number     string
	Code       int
	Month      int
	Year       int
}

type Merchant struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	Name          string
	AccountNumber int32
	Payments      []Payment `gorm:"foreignKey:MerchantID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Payment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Amount     float64
	CustomerID uint    `gorm:"index"`
	MerchantID uint    `gorm:"index"`
	States     []State `gorm:"many2many:payment_states;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type State struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
	//PaymentStates []PaymentState `gorm:"foreignKey:StateID"`
}

func SetState(name StateEnum) State {
	s := State{ID: 0, Name: string(name)}
	switch name {
	case Pending:
		s.ID = 1
	case Succeeded:
		s.ID = 2
	case Rejected:
		s.ID = 3
	case Refunded:
		s.ID = 4
	}
	return s
}

//type PaymentState struct {
//	ID        uint `gorm:"primaryKey"`
//	PaymentID uint `gorm:"index"`
//	StateID   uint `gorm:"index"`
//	SetAt     time.Time
//}

type StateEnum string

const (
	Pending   StateEnum = "Pending"
	Succeeded StateEnum = "Succeeded"
	Rejected  StateEnum = "Rejected"
	Refunded  StateEnum = "Refunded"
)

func (Customer) TableName() string {
	return "customers"
}

func (Merchant) TableName() string {
	return "merchants"
}

func (Payment) TableName() string {
	return "payments"
}

//func (PaymentState) TableName() string {
//	return "payment_states"
//}
