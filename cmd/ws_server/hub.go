package main

import (
	structs "github.com/shuangxing93/tic3online_server/pkg/structs"
)

type PlayerList map[*structs.Player]bool

func NewPlayerList() PlayerList {
	return make(map[*structs.Player]bool)
}
