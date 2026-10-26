package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"tap/internal/protocol"
)

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
		p.Server.mu.Lock()
		p.SendOK("")
		for _, dest := range p.Server.players {
			dest.SendEvt(protocol.EvtChatGlobal, eventData)
		}
		p.Server.mu.Unlock()

	case protocol.ChatGroup:
		p.Server.mu.Lock()
		if p.GroupID == "" {
			p.Server.mu.Unlock()
			p.SendErr(protocol.CodeNotInGroup, protocol.MsgNotInGroup)
			return
		}

		p.SendOK("")

		for _, dest := range p.Server.players {
			if dest.GroupID == p.GroupID {
				dest.SendEvt(protocol.EvtChatGroup, eventData)
			}
		}
		p.Server.mu.Unlock()

	case protocol.ChatRoom:
		p.Server.mu.Lock()
		p.SendOK("")
		for _, dest := range p.Server.players {
			if dest.CurrentRoom == p.CurrentRoom {
				dest.SendEvt(protocol.EvtChatRoom, eventData)
			}
		}
		p.Server.mu.Unlock()

	default:
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
	}
}

func (p *Player) handleWho(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.Server.mu.Lock()
	usernamesList := make([]string, 0, len(p.Server.players))
	for username := range p.Server.players {
		usernamesList = append(usernamesList, username)
	}

	playersNumber := len(p.Server.players)
	p.Server.mu.Unlock()

	// utilisation d'une structure anonyme dans la méthode comme elle ne sert que cette fois-là
	response := struct {
		Room   []string `json:"room"`
		Server int      `json:"Server"`
	}{
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
