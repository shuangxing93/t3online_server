package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/shuangxing93/tic3online_server/pkg/database"
	structs "github.com/shuangxing93/tic3online_server/pkg/structs"
)

type FindGameServer bool

func (findgame *FindGameServer) FindGame(user structs.UserID, state *structs.GameMessage) error {

	db := database.ConnectToDatabase()
	defer db.Close()

	// TODO: find user that is also currently searching. If nil, change our status to searching

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

func getSearchingOpponent(db *sql.DB, userID structs.UserID) int {
	var opponentID int
	err := db.QueryRow("SELECT FROM game_state userID WHERE current_state = (?)", "Searching").Scan(&opponentID)

	if err != nil {
		fmt.Println("No opponent to be found...")
		return -1
	}
	//Opponent found
	err := db.Prepare("UPDATE ")
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
