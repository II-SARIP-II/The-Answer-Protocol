package server

import (
		"tap/internal/protocol"
)

func (p *Player) handleGroup(args string) {
	// TODO
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received GROUP. (TODO)")
}
