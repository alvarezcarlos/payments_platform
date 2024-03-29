package application

import (
	"errors"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/domain/repository"
	"log/slog"
	"time"
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

func (m *merchantUseCase) Create(merchant *entity.Merchant) error {
	merchant.CreatedAt, merchant.UpdatedAt = time.Now(), time.Now()
	if err := m.repository.Create(merchant); err != nil {
		m.logger.Error(err.Error())
		return errors.New("error creating merchant")
	}
	return nil
}
