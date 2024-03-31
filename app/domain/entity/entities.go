package entity

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID         uint `gorm:"primaryKey;autoIncremental"`
	HolderID   uint
	HolderName string
	Balance    float64
	Number     string `gorm:"unique"`
	Code       string
	Month      int
	Year       int
}

type Merchant struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"unique"`
	Balance   float64   `json:"balance"`
	Payments  []Payment `json:"payments" gorm:"foreignKey:MerchantID"`
	Password  string    `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Payment struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Amount     float64   `json:"amount"`
	CardNumber string    `json:"-"`
	MerchantID uint      `json:"merchant_id" gorm:"index"`
	States     []State   `json:"states" gorm:"many2many:payment_states;"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type State struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex"`
}

func SetState(name StateEnum) State {
	s := State{ID: 0, Name: string(name)}
	switch name {
	case Pending:
		s.ID = 1
	case Rejected:
		s.ID = 2
	case Succeeded:
		s.ID = 3
	case Refunded:
		s.ID = 4
	}
	return s
}

type StateEnum string

const (
	Pending   StateEnum = "Pending"
	Succeeded StateEnum = "Succeeded"
	Rejected  StateEnum = "Rejected"
	Refunded  StateEnum = "Refunded"
)

func (Merchant) TableName() string {
	return "merchants"
}

func (Payment) TableName() string {
	return "payments"
}
