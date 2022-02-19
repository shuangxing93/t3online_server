package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Searching string `json:"Searching"`
	IsNoob    string `json:"IsNoob"`
}

// var Channel1test chan bool
var players []Player
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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
func findGameWrap(logger chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//sends message to channel logger when this endpoint is hit
		logger <- "Received request from /findgame"
	}
}

func loginWrap() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		dummyResponse := `{"userId":"1"}`
		io.WriteString(w, dummyResponse)
	}
}

// wrapper function to pass in channel, to allow communication between http and ws
func webSocketWrap(logger chan string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//here we add a check to determine whether an incoming request is allowed to connect or not
		//return true > always allow for now
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err.Error())
		}
		// Once client connected, send hello
		log.Println("Client Connection Received")
		err = ws.WriteMessage(1, []byte("Hello Client! Socket Connection Success!"))
		if err != nil {
			log.Println(err)
			return
		}
		// Sends 1 JSON message upon successful connection
		players = createDummyJSON()
		jerr := ws.WriteJSON(players)
		if jerr != nil {
			log.Println(jerr)
			return
		}
		// Waits for message on logger channel. When message is available send it over to client
		// logger will receive message when http endpoint /findgame is called
		for msg := range logger {
			err = ws.WriteMessage(1, []byte(msg))
			if err != nil {
				log.Println(err)
				return
			}
		}

		readWs(ws)
	}
}
func readWs(w *websocket.Conn) {
	for {
		_, p, err := w.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("message received from client:", string(p))
	}

}

func createDummyJSON() []Player {
	players = []Player{
		{Id: "1", Name: "Adam", Searching: "No", IsNoob: "Yes"},
		{Id: "2", Name: "Steve", Searching: "Yes", IsNoob: "Yes"},
	}
	return players
}

func main() {

	fmt.Println("t3online_server v 0.1")

	findGameChannel := make(chan string)
	findGameRequest := findGameWrap(findGameChannel)
	webSocketConn := webSocketWrap(findGameChannel)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/socket", webSocketConn)
	go myRouter.HandleFunc("/findGame", findGameRequest)
	myRouter.HandleFunc("/login", loginWrap())

	// myRouter := makeRoutes(handlewrap)
	//Fatal is equivalent to Print() followed by a call to os.Exit(1).
	log.Fatal(http.ListenAndServe("192.168.1.216:8080", myRouter))

}

// func testwrapper(logger chan string) {
// 	for item := range logger {
// 		fmt.Println("1. Item", item)
// 	}
// }

//create upgrader to upgrade incoming connection from standard HTTP to long lasting Websocket
