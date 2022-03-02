package structs

import "github.com/gorilla/websocket"

type Player struct {
	WebSocketConnection *websocket.Conn
	Send                chan interface{}
	Username            string
	UserID              int
}

type SocketEventStruct struct {
	EventName    string      `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

type WsMessage struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

//for login rpc
type Username string
type UserID int

//for findgame rpc
//username as above
// type GameState struct {
// 	State    string
// 	Opponent string
// }
type Opponent struct {
	Name string
	ID   int
}

type GameInfo struct {
	GameID int
	Opponent
	Board  [][]string
	IsTurn bool
}

type GameMessage struct {
	State string
	GameInfo
}
