package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=timescaledb password=timescaledb dbname=timescaledb port=5432 sslmode=disable TimeZone=UTC", // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                                                                                                             // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
