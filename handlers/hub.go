// package handlers

// // Hub maintains the set of active clients and broadcasts messages to the clients.
// type Hub struct {
// 	registerPlayer   map[*Player]bool
// 	unregisterPlayer chan *Player
// 	currPlayers      chan *Player
// }

// func NewHub() *Hub {
// 	return &Hub{
// 		registerPlayer:   make(map[*Player]bool),
// 		unregisterPlayer: make(chan *Player),
// 		currPlayers:      make(chan *Player),
// 	}
// }
