package database

import (
	"fmt"
	"hewhew-backend/config"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	*gorm.DB
}

var (
	once                     sync.Once
	postgresDatabaseInstance *postgresDatabase
)

func (db *postgresDatabase) Connect() *gorm.DB {
	return db.DB
}

func NewPostgresDatabase(conf *config.Database) Database {
	once.Do(func() {
		fmt.Printf("Connecting to PostgreSQL database %s at %s:%d...\n", conf.DBName, conf.Host, conf.Port)
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			conf.Host, conf.User, conf.Password, conf.DBName, conf.Port)

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		postgresDatabaseInstance = &postgresDatabase{DB: conn}
		log.Printf("Connected to PostgreSQL database %s", conf.DBName)

	})
	return postgresDatabaseInstance
}
