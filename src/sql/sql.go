package sql

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(user, pass, host, port, dbName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Printf("Error %s when creating DB connection", err)
		panic(fmt.Sprintf("Error %s when creating DB connection", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Error %s when creating DB connection", err)
		panic(fmt.Sprintf("Error %s when creating DB connection", err))
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	fmt.Print("sqlite connected\n")

	return db
}
