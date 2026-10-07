package server

import (
	"tap/internal/protocol"
)

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

func (p *Player) handleQuit(args string) {
	// TODO
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
		}
	p.State = protocol.StateTerminated
	p.SendOK("Received QUIT.")
}

func (p *Player) handleChat(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received CHAT. (TODO)")
}

func (p *Player) handleWho(args string) {
	// TODO
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received WHO. (TODO)")
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
