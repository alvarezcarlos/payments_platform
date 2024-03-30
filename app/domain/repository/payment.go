package repository

import (
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const invalidDuplicatedKeyErrorCode = "23505"

type paymentRepo struct {
	conn *gorm.DB
}

func NewPaymentRepository(conn *gorm.DB) PaymentRepository {
	return &paymentRepo{conn: conn}
}

func (p *paymentRepo) Create(payment *entity.Payment) error {
	if err := p.conn.Create(payment).Error; err != nil {
		return err
	}
	return nil
}
func (p *paymentRepo) Update(payment *entity.Payment) (*entity.Payment, error) {
	tx := p.conn.Begin()

	if err := tx.Save(payment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	updatedPayment := &entity.Payment{}
	if err := p.conn.Preload("States").First(updatedPayment, payment.ID).Error; err != nil {
		return nil, err
	}

	return updatedPayment, nil
}
func (p *paymentRepo) GetByID(id uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	if err := p.conn.Preload("States").First(&payment, "id = ?", id).Error; err != nil {
		return &entity.Payment{}, err
	}
	return &payment, nil
}

func (p *paymentRepo) CreateCard(card *entity.Card) error {
	if err := p.conn.Create(card).Error; err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if pgErr.Code == invalidDuplicatedKeyErrorCode {
			return nil
		}
		return err
	}
	return nil
}
