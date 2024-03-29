package connection

import (
	"fmt"
	"github.com/alvarezcarlos/payment/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

const (
	maxRetries = 5
)

type PostgresRepository interface {
	GetConnection() *gorm.DB
}

type PostgresConnection struct {
	url        string
	config     *gorm.Config
	connection *gorm.DB
	logger     *slog.Logger
}

func NewPostgresConnection(conf *gorm.Config, logger *slog.Logger) *PostgresConnection {
	var dbConfig = config.Config().Database
	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.Port)
	return &PostgresConnection{
		url:        url,
		connection: nil,
		config:     conf,
		logger:     logger,
	}
}
func (pg *PostgresConnection) GetConnection() *gorm.DB {
	var db *gorm.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(pg.url), pg.config)
		if err != nil {
			pg.logger.Error(fmt.Errorf("%w, attempt %d of %d", err, i, maxRetries).Error())
			time.Sleep(3 * time.Second)
			continue
		}
		pg.logger.Debug(" =======> db connected")
		return db
	}
	panic(err)
}
