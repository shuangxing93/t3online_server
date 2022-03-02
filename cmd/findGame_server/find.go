package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	structs "github.com/shuangxing93/tic3online_server/pkg/structs"
)

type FindGameServer bool

func (findgame *FindGameServer) FindGame(user structs.Username, state *structs.GameMessage) error {
	*state = structs.GameMessage{
		State: "Finding",
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
	return nil
}

func main() {
	findgame := new(FindGameServer)
	rpc.Register(findgame)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":7071")
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
