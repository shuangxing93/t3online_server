package main

import (
	"encoding/json"
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
		var gameMessage structs.GameMessage
		err = client.Call("FindGameServer.FindGame", username, &gameMessage)
		if err != nil {
			log.Println("error finding game", err)

		}
		json.NewEncoder(w).Encode(gameMessage)
		gameMessageJSON, err := json.Marshal(gameMessage)
		if err != nil {
			log.Println("error mashalling JSON", err)

		}
		for player := range playerList {
			if player.Username == username {
				select {
				case player.SendGameMessage <- gameMessageJSON:
				default:
					close(player.SendGameMessage)
					delete(playerList, player)
				}
			}
		}
	}
}
func queryURLAndReturnUser(r *http.Request) structs.LoginDetails {
	fmt.Println(r.URL.Query())
	params := r.URL.Query()
	username := params["user"][0]
	passwd := params["passwd"][0]

	user := structs.LoginDetails{
		Username: username,
		Password: passwd,
	}
	return user
}

func registerRequest(w http.ResponseWriter, r *http.Request) {
	//this should get the
	// dummyResponse := `{"userId":"1"}`
	// io.WriteString(w, dummyResponse)

	user := queryURLAndReturnUser(r)
	client, err := rpc.Dial("tcp", ":7070")
	if err != nil {
		log.Fatal("dialing error:", err)
	}
	defer client.Close()
	var userid structs.UserID
	err = client.Call("LoginServer.RegisterUsername", user, &userid)

	if err != nil {
		log.Println("error registering", err)

	}
	var registerSuccess bool
	if userid == -1 {
		registerSuccess = false
	} else {
		registerSuccess = true
	}
	json.NewEncoder(w).Encode(struct{ RegisteredSuccessfully bool }{RegisteredSuccessfully: registerSuccess})

}

func loginWrap(playerList PlayerList) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user := queryURLAndReturnUser(r)

		client, err := rpc.Dial("tcp", ":7070")
		if err != nil {
			log.Fatal("dialing error:", err)
		}
		defer client.Close()
		var userid structs.UserID
		err = client.Call("LoginServer.GetIDbyUsername", user, &userid)
		if err != nil {
			log.Println("error logging in", err)
		}
		json.NewEncoder(w).Encode(struct{ Userid structs.UserID }{Userid: userid})
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
			SendGameMessage:     make(chan interface{}),
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
			close(p.SendGameMessage)
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
	for msg := range p.SendGameMessage {
		gamemessageJSON, err := json.Marshal(msg)
		if err != nil {
			log.Println("write failed", err)
			return
		}
		fmt.Println(string(gamemessageJSON))
		err = p.WebSocketConnection.WriteMessage(1, gamemessageJSON)
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
	myRouter.HandleFunc("/register", registerRequest)

	log.Fatal(http.ListenAndServe(":8080", myRouter))

}
