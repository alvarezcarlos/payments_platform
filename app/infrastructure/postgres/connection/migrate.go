package connection

import (
	"gorm.io/gorm"
	"log/slog"
)

type MigrateInterface interface {
	AutoMigrateAll(tables ...interface{})
}
type migrate struct {
	connection *gorm.DB
	logger     *slog.Logger
}

func NewMigrate(connection *gorm.DB, logger *slog.Logger) MigrateInterface {
	return &migrate{
		connection: connection,
		logger:     logger}
}

func (m *migrate) AutoMigrateAll(tables ...interface{}) {
	err := m.connection.AutoMigrate(tables...)
	if err != nil {
		m.logger.Error("Error migrating tables")
	}
}
