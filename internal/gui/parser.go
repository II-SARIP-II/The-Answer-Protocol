package gui

import (
	"encoding/json"

	"tap/internal/server"
)

func ParsePlayer(str_player string) server.Player {
	str := []byte(str_player)

	var player server.Player
	err := json.Unmarshal(str, &player)
	if err != nil {
		panic(err)
	}
	return player
}



//type Player struct {
//	Server      *Server
//	Username    string
//	Conn        net.Conn
//	State       protocol.State
//	CurrentRoom string
//	// TODO RoomID, HP, Inventory...
//}

