package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shuangxing93/tic3online_server/pkg/structs"
)

func main() {
	fmt.Println("Go MySQL")

	db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	defer db.Close()
	fmt.Println("successfully connected to mysql database")
	stmt, err := db.Prepare("UPDATE users SET gameMessage = ? WHERE username = 'eve'")
	if err != nil {
		log.Fatal("unable to prepare statement", err)
	}
	gameMessage := structs.GameMessage{
		State: "Online",
		GameInfo: structs.GameInfo{
			GameID: 1,
			Opponent: structs.Opponent{
				Name: "",
				ID:   -1,
			},
			Board: [][]string{
				{"", "", ""},
				{"", "", ""},
				{"", "", ""},
			},
			IsTurn: false,
		},
	}
	gameMessageJSON, err := json.Marshal(gameMessage)
	if err != nil {
		log.Println("error marshalling JSON", err)
	}
	result, err := stmt.Exec(gameMessageJSON)
	// fmt.Println(test)
	lastinsertedid, _ := result.LastInsertId()
	rowsaffected, _ := result.RowsAffected()
	fmt.Println(lastinsertedid, rowsaffected)

	if err != nil {
		log.Fatal("unable to execute statement", err)
	}

}
