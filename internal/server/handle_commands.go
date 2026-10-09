package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"tap/internal/protocol"
)

func (p *Player) handleQuit(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.State = protocol.StateTerminated
	p.SendOK("Received QUIT.")
}

type whoResponse struct {
	Room   []string `json:"room"`
	Server int      `json:"server"`
}

func (p *Player) handleWho(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.server.mu.Lock()
	usernamesList := make([]string, 0, len(p.server.players))
	for username := range p.server.players {
		usernamesList = append(usernamesList, username)
	}
	// TODO ajouter la vraie liste de la current room
	playersNumber := len(p.server.players)
	p.server.mu.Unlock()

	response := whoResponse{
		Room:   usernamesList,
		Server: playersNumber,
	}

	data, error := json.Marshal(response)
	if error != nil {
		p.SendErr(protocol.CodeSendFailed, protocol.MsgSendFailed)
		return
	}
	p.SendOK(string(data))

}

func (p *Player) handleLook(args string) {
	// TODO
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received LOOK. (TODO)")
}

func (p *Player) handleMove(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received LOOK. (TODO)")
}

func (p *Player) handleChat(args string) {
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	argsParts := strings.SplitN(args, " ", 2)
	if len(argsParts) < 2 {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	scope := strings.ToUpper(argsParts[0])
	chatMsg := strings.TrimSpace(argsParts[1])
	if chatMsg == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	eventData := fmt.Sprintf("%s %s", p.Username, chatMsg)

	switch scope {
	case protocol.ChatGlobal:
		p.SendOK("")
		p.server.mu.Lock()
		for _, dest := range p.server.players {
			dest.SendEvt("GLOBAL CHAT", eventData)
		}
		p.server.mu.Unlock()

	case protocol.ChatGroup:
		// TODO à compléter si player appartient à un groupe
		p.SendErr(protocol.CodeNotInGroup, protocol.MsgNotInGroup)

	case protocol.ChatRoom:
		p.SendOK("")
		p.server.mu.Lock()
		for _, dest := range p.server.players {
			if dest.CurrentRoom == p.CurrentRoom {
				dest.SendEvt("ROOM CHAT", eventData)
			}
		}
		p.server.mu.Unlock()

	default:
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
}

func (p *Player) handleGroup(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received GROUP. (TODO)")
}

func (p *Player) handleTake(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received TAKE. (TODO)")
}

func (p *Player) handleDrop(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received DROP. (TODO)")
}

func (p *Player) handleInventory(args string) {
	// TODO
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received INVENTORY. (TODO)")
}

func (p *Player) handleTalk(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received TALK. (TODO)")
}

func (p *Player) handleAttack(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received ATTACK. (TODO)")
}

func (p *Player) handleStatus(args string) {
	// TODO
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received STATUS. (TODO)")
}

func (p *Player) handleQuest(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received QUEST. (TODO)")
}

func (p *Player) handleQuests(args string) {
	// TODO
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received QUESTS. (TODO)")
}
