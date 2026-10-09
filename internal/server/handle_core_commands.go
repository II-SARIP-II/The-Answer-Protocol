package server

import (
	"tap/internal/protocol"
)

func (p *Player) handleLook(args string) {
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
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.State = protocol.StateTerminated
	p.SendOK("Received QUIT.")
}

