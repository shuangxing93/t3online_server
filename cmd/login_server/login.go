package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	structs "github.com/shuangxing93/tic3online_server/pkg/structs"
)

type LoginServer bool

func (login *LoginServer) GetIDbyUsername(user structs.Username, id *structs.UserID) error {
	*id = 1
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
