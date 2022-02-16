package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func makeRoutes() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/socket", wsEndPoint)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	// Fprintf prints to an io writer instead of to the terminal
	fmt.Fprintf(w, "hello home page")
}
func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	//here we add a check to determine whether an incoming request is allowed to connect or not
	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err.Error())
	}
	log.Println("client connected")
	err = ws.WriteMessage(1, []byte("Hello client"))
	if err != nil {
		log.Println(err)
	}

	Readws(ws)
}

func Readws(w *websocket.Conn) {
	for {
		messageType, p, err := w.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))
		err = w.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			return
		}

	}
}

func main() {

	fmt.Println("ws tutorial")

	makeRoutes()
	//Fatal is equivalent to Print() followed by a call to os.Exit(1).
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}

//create upgrader to upgrade incoming connection from standard HTTP to long lasting Websocket
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
