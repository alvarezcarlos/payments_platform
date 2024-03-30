package entity

import (
	"github.com/google/uuid"
	"time"
)

type Card struct {
	ID         uint `gorm:"primaryKey;autoIncremental"`
	HolderID   uint
	HolderName string
	Number     string `gorm:"unique"`
	Code       string
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
	CardNumber string
	MerchantID uint    `gorm:"index"`
	States     []State `gorm:"many2many:payment_states;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type State struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
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
