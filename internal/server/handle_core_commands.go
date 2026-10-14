package server

import (
	"encoding/json"
	"tap/internal/protocol"
)

type LookInfoRoom struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
}

type LookResponse struct {
	Room    LookInfoRoom `json:"room"`
	Players []string     `json:"players"`
	Items   []string     `json:"items"`
	NPCs    []string     `json:"npcs"`
}

func (p *Player) handleLook(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	room, exists := p.Server.world.Rooms[p.CurrentRoom]
	if !exists {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
		return
	}

	exitsMap := make(map[string]string)
	for dir, exit := range room.Exits {
		exitsMap[dir] = exit.Target
	}

	p.Server.mu.Lock()
	playersInRoom := make([]string, 0)
	for _, otherPlayer := range p.Server.players {
		if otherPlayer.CurrentRoom == p.CurrentRoom {
			playersInRoom = append(playersInRoom, otherPlayer.Username)
		}
	}
	p.Server.mu.Unlock()

	items := room.Items
	if items == nil {
		items = []string{}
	}
	npcs := room.NPCs
	if npcs == nil {
		npcs = []string{}
	}

	response := LookResponse{
		Room: LookInfoRoom{
			Id:          p.CurrentRoom,
			Name:        room.Name,
			Description: "", //TODO ajouter un champ description dans les structures Room
			Exits:       exitsMap,
		},
		Players: playersInRoom,
		Items:   items,
		NPCs:    npcs,
	}

	data, err := json.Marshal(response)
	if err != nil {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
		return
	}

	p.SendOK(string(data))
}

func (p *Player) handleMove(args string) {
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received LOOK. (TODO)")
}

func (p *Player) handleQuit(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.State = protocol.StateTerminated
	p.SendOK("Received QUIT.")
}
