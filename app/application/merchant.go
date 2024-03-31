package application

import (
	"errors"
	"log/slog"
	"time"

	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"github.com/alvarezcarlos/payment/app/utils"
)

type merchantUseCase struct {
	repository repository.MerchantRepository
	logger     *slog.Logger
}

func NewMerchantUseCase(repository repository.MerchantRepository, logger *slog.Logger) MerchantUseCaseInterface {
	return &merchantUseCase{
		repository: repository,
		logger:     logger}
}

func (m *merchantUseCase) Create(merchant *entity.Merchant) (*entity.Merchant, error) {
	merchant.CreatedAt, merchant.UpdatedAt = time.Now(), time.Now()
	merchant.Balance = utils.RandomFloat()
	merch, err := m.repository.Create(merchant)
	if err != nil {
		m.logger.Error(err.Error())
		return nil, errors.New("error creating merchant")
	}
	return merch, nil
}

func (m *merchantUseCase) GetByName(name string) (*entity.Merchant, error) {
	return m.repository.GetByName(name)
}
