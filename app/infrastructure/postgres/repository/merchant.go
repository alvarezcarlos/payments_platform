package repository

import (
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"gorm.io/gorm"
)

type merchantRepo struct {
	conn *gorm.DB
}

func NewMerchantRepository(conn *gorm.DB) repository.MerchantRepository {
	return &merchantRepo{conn: conn}
}

func (m *merchantRepo) Create(merchant *entity.Merchant) (*entity.Merchant, error) {
	if tx := m.conn.Create(merchant); tx.Error != nil {
		return nil, tx.Error
	}
	return merchant, nil
}

func (p *merchantRepo) GetByName(name string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	if err := p.conn.Where("name = ?", name).First(&merchant).Error; err != nil {
		return nil, err
	}
	return &merchant, nil
}
