package server

import (
	"tap/internal/protocol"
	"strings"
)

func (p *Player) handleGroup(args string) {
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	request := strings.ToLower(strings.TrimSpace(args))
	requestParts := strings.SplitN(request, " ", 3)
	if len(requestParts) != 3 {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	if requestParts[0] != "group" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	// To debug

	// subCommand := requestParts[1]
	// lastPartRequest := requestParts[2]

	// switch subCommand {
	// case protocol.GroupCreate:
	// 	p.handleGroupCreate(lastPartRequest)
	// case protocol.GroupInvite:
	// 	p.handleGroupInvite(lastPartRequest)
	// }
	// TODO .....

	p.SendOK("Received GROUP. (TODO)")
}
