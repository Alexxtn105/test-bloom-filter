package models

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

var bloomFilter *bloom.BloomFilter

func DBInit() error {
	// Настраиваем логгер
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Microsecond, // Slow SQL threshold
			LogLevel:                  logger.Info,      // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound
			Colorful:                  true,             // Disable color
		},
	)

	// Инициализация БД
	var err error
	db, err = gorm.Open(sqlite.Open("storage/storage.db"),
		&gorm.Config{Logger: newLogger})
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	// Автомиграция
	err = db.AutoMigrate(&UserAccess{})
	if err != nil {
		return err
	}
	return nil
}

// BloomFilterSetup Настройка bloom-фильтра
func BloomFilterSetup() {
	// инициализация bloom-filter с оценками (Estimates)
	bloomFilter = bloom.NewWithEstimates(
		10000, // примерное количество элементов, которые мы добавим в фильтр
		0.01,  // желаемый false-positive rate (1 раз в 100 запросов может быть ошибочный positive)
	)
}
