package server

import(
	"strings"
	"tap/internal/protocol"
)

func (p *Player) HandleCommand(rawLine string) {
	parts := strings.SplitN(rawLine, " ", 2)
	command := strings.ToUpper(parts[0])
	args := ""
	if len(parts) > 1 {
		args := strings.TrimSpace(parts[1])
	}

	if p.State != protocol.StateAuthenticated {
		p.sendErr(protocol.CodeNotAuthenticated, protocol.MsgNotAuthenticated)
		return
	}

	switch command {

	case protocol.CmdLook:
		p.handleLook()
	case protocol.CmdMove:
		p.handleMove(args)
	case protocol.CmdQuit:
		p.handleQuit()
	case protocol.CmdChat:
		p.handleChat(args)
	case protocol.CmdWho:
		p.handleWho()
	case protocol.CmdGroup:
		p.handleGroup(args)
	case protocol.CmdTake:
		p.handleTake(args)
	case protocol.CmdDrop:
		p.handleDrop(args)
	case protocol.CmdInventory:
		p.handleInventory()
	case protocol.CmdTalk:
		p.handleTalk(args)
	case protocol.CmdAttack:
		p.handleAttack(args)
	case protocol.CmdStatus:
		p.handleStatus()
	case protocol.CmdQuest:
		p.handleQuest(args)
	case protocol.CmdQuests:
		p.handleQuests()
	case protocol.CmdConnect:
		p.SendErr(protocol.CodeAlreadyAuthenticated, protocol.MsgAlreadyAuthenticated)
	default:
		p.SendErr(protocol.CodeUnknownCommand, protocol.MsgUnknownCommand)
	}
}






















