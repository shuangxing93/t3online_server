package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/rpc"

	_ "github.com/go-sql-driver/mysql"
	structs "github.com/shuangxing93/tic3online_server/pkg/structs"
)

type LoginServer bool

func (login *LoginServer) GetIDbyUsername(user structs.Username, id *structs.UserID) error {
	fmt.Println("connecting to MySQL database...")

	db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	defer db.Close()
	var userID int
	fmt.Println("successfully connected to mysql database")
	row := db.QueryRow("SELECT userID FROM users where username = ?", string(user))
	err = row.Scan(&userID)
	if err != nil {
		log.Fatal("Username does not exist.", err)
		*id = -1
		return nil
	}
	*id = structs.UserID(userID)
	fmt.Println("userid is", userID)
	return nil
}

func (login *LoginServer) RegisterUsername(user structs.Username, id *structs.UserID) error {
	fmt.Println("connecting to MySQL database...")

	db, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		log.Fatal("cannot connect to MySQL database", err)
	}

	defer db.Close()

	fmt.Println("successfully connected to MySQL database")

	stmt, err := db.Prepare("INSERT INTO users (username) VALUES (?) ")
	if err != nil {
		log.Fatal("unable to prepare statement", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(user)

	if err != nil {
		log.Println("Error - Username already exists.")
		*id = -1
		return nil
	}
	fmt.Println("registration success!")
	var userID int
	row := db.QueryRow("SELECT LAST_INSERT_ID()")
	err = row.Scan(&userID)
	if err != nil {
		log.Fatal("can't get last inserted ID", err)
	}
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
