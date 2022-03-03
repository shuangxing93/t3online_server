package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	_ "github.com/go-sql-driver/mysql"
	database "github.com/shuangxing93/tic3online_server/pkg/database"
	structs "github.com/shuangxing93/tic3online_server/pkg/structs"
	"golang.org/x/crypto/bcrypt"
)

type LoginServer bool

func (logOrReg *LoginServer) GetIDbyUsername(user structs.LoginDetails, id *structs.UserID) error {

	db := database.ConnectToDatabase()

	defer db.Close()
	var userID int
	var passwd string
	fmt.Println("successfully connected to mysql database")

	row := db.QueryRow("SELECT userID, passwd FROM users WHERE username = ?", string(user.Username))
	err := row.Scan(&userID, &passwd)
	if err != nil {
		log.Fatal("Username does not exist.")
		*id = -1
		return nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwd), []byte(user.Password))
	if err != nil {
		log.Fatal("Password do not match.")
		*id = -1
		return nil
	}

	*id = structs.UserID(userID)
	fmt.Println("userid is", userID)
	return nil
}

func (logOrReg *LoginServer) RegisterUsername(user structs.LoginDetails, id *structs.UserID) error {
	fmt.Println("connecting to MySQL database...")

	db := database.ConnectToDatabase()

	defer db.Close()

	fmt.Println("successfully connected to MySQL database")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Fatal("error hashing password", err)
	}

	//use db begin to achieve atomicity
	txn, err := db.Begin()
	if err != nil {
		log.Println("failed to start transaction")
		*id = -1
		return nil
	}

	response, err := txn.Exec("INSERT INTO user_data (username,password) VALUES (?,?) ", user.Username, hashedPassword)
	if err != nil {
		txn.Rollback()
		log.Println("Username already exist.")
		*id = -1
		return nil
	}
	userID, err := response.LastInsertId()
	if err != nil {
		txn.Rollback()
		log.Println("Failed to retrieve ID, please try again")
		*id = -1
		return nil
	}
	_, err = txn.Exec("INSERT INTO game_state (userID) VALUES (?) ", userID)
	if err != nil {
		txn.Rollback()
		log.Println("Error - connection failure. Please try again.")
		*id = -1
		return nil
	}
	err = txn.Commit()

	if err != nil {
		log.Println("Failed to register. Please try again.")
		*id = -1
		return nil
	}
	// defer stmt.Close()
	// var userID int
	// row := db.QueryRow("SELECT LAST_INSERT_ID()")
	// err = row.Scan(&userID)
	// if err != nil {
	// 	log.Fatal("can't get last inserted ID", err)
	// }
	fmt.Println("registration success!")
	// var userID int
	fmt.Println("userid is", userID)
	*id = structs.UserID(userID)
	return nil
}

func main() {
	login := new(LoginServer)
	rpc.Register(login)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":7070")
	if err != nil {
		log.Println("error", err)
	}
	fmt.Println(tcpAddr)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Println("error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error", err)
			continue
		}
		rpc.ServeConn(conn)
	}
}
