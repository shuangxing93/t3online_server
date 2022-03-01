package main

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	structs "github.com/shuangxing93/tic3online_server/pkg/structs"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TicTacToe home page")
}

// wrapper function to pass in channel, to allow communication between http and ws
func findGameWrap(playerList PlayerList) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		username := params["user"][0]

		client, err := rpc.Dial("tcp", ":7071")
		if err != nil {
			log.Fatal("dialing error:", err)
		}
		var gamestate structs.GameState
		err = client.Call("FindGameServer.FindGame", username, &gamestate)
		if err != nil {
			log.Println("error finding game", err)

		}

		for player := range playerList {
			if player.Username == username {
				select {
				case player.Send <- fmt.Sprintf("hello, %v", player.Username):
				default:
					close(player.Send)
					delete(playerList, player)
				}
			}
		}
	}
}

func loginWrap(playerList PlayerList) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//this should get the
		// dummyResponse := `{"userId":"1"}`
		// io.WriteString(w, dummyResponse)

		client, err := rpc.Dial("tcp", ":7070")
		if err != nil {
			log.Fatal("dialing error:", err)
		}
		var username structs.Username = "Adam"
		var userid structs.UserID
		err = client.Call("LoginServer.GetIDbyUsername", username, &userid)
		if err != nil {
			log.Println("error logging in", err)

		}
	}
}

func webSocketWrap(playerList PlayerList) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//here we add a check to determine whether an incoming request is allowed to connect or not
		//return true > always allow for now
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error,", err)
		}

		player := &structs.Player{
			WebSocketConnection: ws,
			Send:                make(chan interface{}),
			Username:            "",
			UserID:              0,
		}

		playerList[player] = true
		fmt.Println("current playerList is", playerList)

		log.Println("Client Connection Received")
		err = ws.WriteMessage(1, []byte("Hello Client! Socket Connection Success!"))
		if err != nil {
			log.Println("write failed", err)
			return
		}

		go playerList.writeToPlayer(player)
		go playerList.readFromPlayer(player)
	}
}

func (playerList PlayerList) readFromPlayer(p *structs.Player) {
	for {
		msgReceived := structs.WsMessage{}
		err := p.WebSocketConnection.ReadJSON(&msgReceived)
		if err != nil {
			log.Println("Player Disconnected, removing from player list")
			close(p.Send)
			delete(playerList, p)
			break
		}
		fmt.Println(msgReceived)
		fmt.Println(msgReceived.Username)

		if msgReceived.Id == "1" {
			p.Username = msgReceived.Username

			greetingMsg := fmt.Sprintf("hello, %v", p.Username)
			err = p.WebSocketConnection.WriteMessage(1, []byte(greetingMsg))
			if err != nil {
				log.Println("write failed", err)
				return
			}
		}
	}
}
func (playerList PlayerList) writeToPlayer(p *structs.Player) {
	for msg := range p.Send {
		fmt.Println(msg)
		err := p.WebSocketConnection.WriteMessage(1, []byte(fmt.Sprint(msg)))
		if err != nil {
			log.Println(err)
		}
	}

}

func main() {

	fmt.Println("t3online_server v 0.1")
	playerList := NewPlayerList()

	findGameRequest := findGameWrap(playerList)
	webSocketConn := webSocketWrap(playerList)
	loginRequest := loginWrap(playerList)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/socket", webSocketConn)
	myRouter.HandleFunc("/findGame", findGameRequest)
	myRouter.HandleFunc("/login", loginRequest)

	log.Fatal(http.ListenAndServe(":8080", myRouter))

}
