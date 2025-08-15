package db

import (
	"fmt"
	"log/slog"

	"github.com/sidiik/moonpay/auth_service/internal/domain"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnections struct {
	Writer *gorm.DB
	Reader *gorm.DB
}

func NewPostgres() (*DBConnections, error) {
	// writer db
	writerDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBWriterHost,
		config.AppConfig.DBWriterUser,
		config.AppConfig.DBWriterPassword,
		config.AppConfig.DBWriterName,
		config.AppConfig.DBWriterPort,
	)

	writerDb, err := gorm.Open(postgres.Open(writerDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect writer DB: %w", err)
	}

	// reader db
	readerDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.AppConfig.DBReaderHost,
		config.AppConfig.DBReaderUser,
		config.AppConfig.DBReaderPassword,
		config.AppConfig.DBReaderName,
		config.AppConfig.DBWriterPort,
	)

	readerDB, err := gorm.Open(postgres.Open(readerDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect reader DB: %w", err)
	}

	// Auto migrate schema
	if err := writerDb.AutoMigrate(&domain.User{}); err != nil {
		return nil, err
	}

	slog.Info("Writer and Reader DB connected successfully")

	return &DBConnections{
		Writer: writerDb,
		Reader: readerDB,
	}, nil
}
