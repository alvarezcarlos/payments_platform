package repository

import (
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"gorm.io/gorm"
)

type merchantRepo struct {
	conn *gorm.DB
}

func NewMerchantRepository(conn *gorm.DB) MerchantRepository {
	return &merchantRepo{conn: conn}
}

func (m *merchantRepo) Create(merchant *entity.Merchant) error {
	if tx := m.conn.Create(merchant); tx.Error != nil {
		return tx.Error
	}
	return nil
}
