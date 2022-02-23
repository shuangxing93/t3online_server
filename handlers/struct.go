// package handlers

// import "github.com/gorilla/websocket"

// // Client is a middleman between the websocket connection and the hub.
// type Player struct {
// 	hub                 *Hub
// 	webSocketConnection *websocket.Conn
// 	send                chan SocketEventStruct
// 	username            string
// 	userID              string
// }

// type SocketEventStruct struct {
// 	EventName    string      `json:"eventName"`
// 	EventPayload interface{} `json:"eventPayload"`
// }
