package connection

import (
	"fmt"
	"github.com/alvarezcarlos/payment/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
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
	db, err := gorm.Open(postgres.Open(pg.url), pg.config)
	if err != nil {
		panic("failed to connect database")
	}
	pg.logger.Debug(" =======> db connected")
	return db
}
