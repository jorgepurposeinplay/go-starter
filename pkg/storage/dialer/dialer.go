package dialer

import (
	"log"

	storage2 "github.com/oakeshq/go-starter/internal/storage"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"github.com/oakeshq/go-starter/config"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Connect will connect to that dialer engine
func Connect(config *config.Config) (*gorm.DB, error) {
	var db *gorm.DB

	operation := func() error {
		conn, err := gorm.Open(postgres.Open(config.DB.URL), &gorm.Config{})
		db = conn

		if err != nil {
			return errors.Wrap(err, "opening database connection")
		}

		return nil
	}

	err := backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 5))
	if err != nil {
		return nil, err
	}

	db.Logger = db.Logger.LogMode(logger.Info)

	sqlDB, err := db.DB()

	if err == nil {
		err = sqlDB.Ping()
	}

	if err != nil {
		return nil, errors.Wrap(err, "checking database connection")
	}

	return db, nil
}

// Migrate runs the gorm migration for all models
func Migrate(db *gorm.DB) error {
	allModels := []interface{}{&storage2.User{}}

	if err := db.Migrator().DropTable(allModels...); err != nil {
		log.Printf("Failed to drop table, got error %v\n", err)
		return err
	}

	if err := db.AutoMigrate(allModels...); err != nil {
		return err
	}

	for _, m := range allModels {
		if !db.Migrator().HasTable(m) {
			log.Fatalf("Failed to create table for %#v\n", m)
		}
	}

	return nil
}
