package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ------------------from hub.go----------------------
//
//

type Hub struct {
	// registerPlayer   chan *Player
	// unregisterPlayer chan *Player
	currPlayers map[*Player]bool
}

func NewHub() *Hub {
	return &Hub{
		// registerPlayer:   make(chan *Player),
		// unregisterPlayer: make(chan *Player),
		currPlayers: make(map[*Player]bool),
	}
}

//
//------------------------from struct.go ------------------------------------
//
//

type Player struct {
	hub                 *Hub
	webSocketConnection *websocket.Conn
	send                chan interface{}
	username            string
	userID              string
}

type SocketEventStruct struct {
	EventName    string      `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

//
//
// -----------------------------------------------------------------

// type Player struct {
// 	Id        string `json:"Id"`
// 	Name      string `json:"Name"`
// 	Searching string `json:"Searching"`
// 	IsNoob    string `json:"IsNoob"`
// }

// var Channel1test chan bool
// var players []Player
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsMessage struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

// func makeRoutes(handlewrap func(http.ResponseWriter, *http.Request)) *mux.Router {
// 	// use gorilla mux instead of net/http
// 	// myRouter := mux.NewRouter().StrictSlash(true)

// 	// myRouter.HandleFunc("/", homepage)
// 	// myRouter.HandleFunc("/testget", handlewrap)
// 	// myRouter.HandleFunc("/socket", wsEndPoint)
// 	return myRouter
// }

func homepage(w http.ResponseWriter, r *http.Request) {
	// Fprintf prints to an io writer instead of to the terminal
	fmt.Fprintf(w, "TicTacToe home page")
}

// wrapper function to pass in channel, to allow communication between http and ws
func findGameWrap(hub *Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(r.URL.Query())
		params := r.URL.Query()
		username := params["user"][0]

		for player := range hub.currPlayers {
			if player.username == username {
				select {
				case player.send <- fmt.Sprintf("hello, %v", player.username):
				default:
					close(player.send)
					delete(hub.currPlayers, player)
				}
			}
		}
		//sends message to channel logger when this endpoint is hit
		// logger <- "Received request from /findgame"
	}
}

func loginWrap(hub *Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dummyResponse := `{"userId":"1"}`
		io.WriteString(w, dummyResponse)
	}
}

// // readPump pumps messages from the websocket connection to the hub.
// func (c *Player) readPump() {

// 	defer c.webSocketConnection.Close() //executed last

// 	msgType, body_byte, err := c.webSocketConnection.ReadMessage()

// 	body_str = string(body_byte)

// 	if err != nil {
// 		log.Println("err reading socket msg --  ",err.Error())
// 	}

// 	c.hub.unregisterPlayer <- c

// }

// // writePump pumps messages from the hub to the websocket connection.
// func (c *Player) writePump() {
// }

// wrapper function to pass in channel, to allow communication between http and ws
func webSocketWrap(hub *Hub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//here we add a check to determine whether an incoming request is allowed to connect or not
		//return true > always allow for now
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error,", err)
		}

		// fmt.Println(r.URL.Query())
		// params := r.URL.Query()
		// username := params["user"][0]

		// u, err := url.Parse(s)
		// fmt.Println(username, "has logged in")

		// fmt.Println(username)
		//
		uniqueID := uuid.New()
		player := &Player{
			hub:                 hub,
			webSocketConnection: ws,
			send:                make(chan interface{}),
			username:            "",
			userID:              uniqueID.String(),
		}
		// fmt.Println("player connected: id is", player.userID)

		//register player
		hub.currPlayers[player] = true
		fmt.Println("current hub is", hub)
		// fmt.Println(u.RawQuery)

		// Once client connected, send hello
		log.Println("Client Connection Received")
		err = ws.WriteMessage(1, []byte("Hello Client! Socket Connection Success!"))
		if err != nil {
			log.Println("write failed", err)
			return
		}
		// for {

		// 	err = ws.WriteMessage(1, []byte("sending dummy msg..."))
		// 	if err != nil {
		// 		log.Println(err)
		// 		return
		// 	}
		// 	time.Sleep(5 * time.Second)
		// }
		// Sends 1 JSON message upon successful connection
		// players = createDummyJSON()
		// jerr := ws.WriteJSON(players)
		// if jerr != nil {
		// 	log.Println(jerr)
		// 	return
		// }
		// Waits for message on logger channel. When message is available send it over to client
		// logger will receive message when http endpoint /findgame is called
		// for msg := range logger {
		// 	err = ws.WriteMessage(1, []byte(msg))
		// 	if err != nil {
		// 		log.Println(err)
		// 		return
		// 	}
		// }

		go player.writePump()
		// go dcCleanUp(hub, player, ws)
		go player.getinitmsg()
	}
}

// func dcCleanUp(hub *Hub, p *Player, c *websocket.Conn) {
// 	for {
// 		if _, _, err := c.NextReader(); err != nil {
// 			fmt.Println(p.username, "disconnected.")
// 			delete(hub.currPlayers, p)
// 			c.Close()
// 			break
// 		}
// 	}
// }

func (p *Player) getinitmsg() {
	for {
		msgReceived := wsMessage{}
		err := p.webSocketConnection.ReadJSON(&msgReceived)
		if err != nil {
			log.Println("read error", err)
			// break
		}
		fmt.Println(msgReceived)
		fmt.Println(msgReceived.Username)

		if msgReceived.Id == "1" {
			p.username = msgReceived.Username

			greetingMsg := fmt.Sprintf("hello, %v", p.username)
			err = p.webSocketConnection.WriteMessage(1, []byte(greetingMsg))
			if err != nil {
				log.Println("write failed", err)
				return
			}
		}
	}
}
func (p *Player) writePump() {
	for msg := range p.send {
		fmt.Println(msg)
		err := p.webSocketConnection.WriteMessage(1, []byte(fmt.Sprint(msg)))
		if err != nil {
			log.Println(err)
		}
	}

}

// func readWs(w *websocket.Conn) {
// 	for {
// 		_, p, err := w.ReadMessage()
// 		if err != nil {
// 			hub.unregister(player)
// 			log.Println(err)
// 			return
// 		}

// 		fmt.Println("message received from client:", string(p))
// 	}

// }

// func createDummyJSON() []Player {
// 	players = []Player{
// 		{Id: "1", Name: "Adam", Searching: "No", IsNoob: "Yes"},
// 		{Id: "2", Name: "Steve", Searching: "Yes", IsNoob: "Yes"},
// 	}
// 	return players
// }

func main() {

	fmt.Println("t3online_server v 0.1")
	hub := NewHub()
	// findGameChannel := make(chan string)
	findGameRequest := findGameWrap(hub)
	webSocketConn := webSocketWrap(hub)
	loginRequest := loginWrap(hub)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/socket", webSocketConn)
	myRouter.HandleFunc("/findGame", findGameRequest)
	myRouter.HandleFunc("/login", loginRequest)

	// myRouter := makeRoutes(handlewrap)
	//Fatal is equivalent to Print() followed by a call to os.Exit(1).
	log.Fatal(http.ListenAndServe(":8080", myRouter))

}

// func testwrapper(logger chan string) {
// 	for item := range logger {
// 		fmt.Println("1. Item", item)
// 	}
// }

//create upgrader to upgrade incoming connection from standard HTTP to long lasting Websocket
