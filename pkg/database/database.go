package database

import (
	"database/sql"
	"fmt"
	"log"

	config "github.com/shuangxing93/tic3online_server/pkg/config"
)

func ConnectToDatabase() *sql.DB {
	config, err := config.LoadConfig("../../")
	fmt.Println("connecting to MySQL database...")
	if err != nil {
		log.Fatal("cannot load config file", err)
	}
	dbpath := fmt.Sprintf("%s:%s@%s/%s", config.Database.DBUser, config.Database.DBPassword, config.Database.ServerAddress, config.Database.DBTable)
	dbdriver := config.Database.DBDriver
	db, err := sql.Open(dbdriver, dbpath)

	if err != nil {
		log.Fatal("cannot connect to database", err)
	}
	return db
}
