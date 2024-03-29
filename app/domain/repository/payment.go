package repository

import (
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentRepo struct {
	conn *gorm.DB
}

func NewPaymentRepository(conn *gorm.DB) PaymentRepository {
	return &paymentRepo{conn: conn}
}

func (p *paymentRepo) Create(payment *entity.Payment) error {
	if tx := p.conn.Create(payment); tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (p *paymentRepo) Update(payment *entity.Payment) (entity.Payment, error) {
	return entity.Payment{}, nil
}
func (p *paymentRepo) GetByID(id uuid.UUID) (entity.Payment, error) {
	return entity.Payment{}, nil
}
