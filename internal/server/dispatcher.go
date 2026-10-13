package server

import (
	"log/slog"
	"strings"
	"tap/internal/protocol"
)

func (p *Player) HandleCommand(rawLine string) {
	parts := strings.SplitN(rawLine, " ", 2)
	command := strings.ToUpper(parts[0])
	args := ""
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}

	slog.Info("Command received",
		"player", p.Username,
		"command", command,
		"parameters", args,
	)

	if p.State != protocol.StateAuthenticated {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgNotAuthenticated)
		return
	}

	switch command {

	case protocol.CmdLook:
		p.handleLook(args)
	case protocol.CmdMove:
		p.handleMove(args)
	case protocol.CmdQuit:
		p.handleQuit(args)
	case protocol.CmdChat:
		p.handleChat(args)
	case protocol.CmdWho:
		p.handleWho(args)
	case protocol.CmdGroup:
		p.handleGroup(args)
	case protocol.CmdTake:
		p.handleTake(args)
	case protocol.CmdDrop:
		p.handleDrop(args)
	case protocol.CmdInventory:
		p.handleInventory(args)
	case protocol.CmdTalk:
		p.handleTalk(args)
	case protocol.CmdAttack:
		p.handleAttack(args)
	case protocol.CmdStatus:
		p.handleStatus(args)
	case protocol.CmdQuest:
		p.handleQuest(args)
	case protocol.CmdQuests:
		p.handleQuests(args)
	case protocol.CmdConnect:
		p.SendErr(protocol.CodeBadRequest, protocol.MsgAlreadyAuthenticated)
	default:
		p.SendErr(protocol.CodeBadRequest, protocol.MsgUnknownCommand)
	}
}
