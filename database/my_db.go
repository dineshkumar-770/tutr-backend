package database

import (
	"database/sql"
	"fmt"
	"log"
	"main/utils"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func Initialize() *sql.DB {
	envs, err := utils.GetEnvVars()
	if err != nil {
		fmt.Println(err)
	}
	once.Do(func() {
		var err error
		dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", envs.DatabaseRole, envs.DatabasePassword, envs.DatabaseIPAddress, envs.DatabasePORT, envs.DatabaseName)
		dbInstance, err = sql.Open("mysql", dbUrl)
		if err != nil {
			log.Fatalf("Failed to connect to MySQL: %v", err)
		}

		if err = dbInstance.Ping(); err != nil {
			log.Fatalf("Failed to ping MySQL: %v", err)
		}

		log.Println("Connected to MySQL successfully")
	})

	return dbInstance
}

func GetDBInstance() *sql.DB {
	if dbInstance == nil {
		log.Fatal("Database not initialized. Call Initialize first.")
	}
	return dbInstance
}
