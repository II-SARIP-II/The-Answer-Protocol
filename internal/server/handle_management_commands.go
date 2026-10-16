package server

import (
	"fmt"
	"strings"
	"tap/internal/protocol"
)

func (p *Player) handleGroup(args string) {
	if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	parts := strings.SplitN(strings.TrimSpace(args), " ", 2)
	subCommand := strings.ToUpper(parts[0])
	subArgs := ""
	if len(parts) > 1 {
		subArgs = strings.TrimSpace(parts[1])
	}

	switch subCommand {
	case protocol.GroupCreate:
		p.handleGroupCreate(subArgs)

	case protocol.GroupInvite:
		p.handleGroupInvite(subArgs)

	case protocol.GroupJoin:
		p.handleGroupJoin(subArgs)

	case protocol.GroupLeave:
		p.handleGroupLeave(subArgs)

	default:
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
	}
}

func (p *Player) handleGroupCreate(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()
	if p.GroupID != "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeAlreadyInGroup, protocol.MsgAlreadyInGroup)
		return
	}

	groupID := fmt.Sprintf("group_%s", p.Username)
	p.GroupID = groupID
	p.Server.mu.Unlock()
	p.SendOK(fmt.Sprintf("group=%s", groupID))
}

// choix fait: tous les membres d'un groupe peuvent inviter qqun,
// pas seulement le leader
// le message d'EVT contient bien le nom du leader
// pr que le player sache quel groupe l'invite

func (p *Player) handleGroupInvite(args string) {
	targetInvited := strings.TrimSpace(args)
	if targetInvited == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()

	if p.GroupID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNotInGroup, protocol.MsgNotInGroup)
		return
	}

	if strings.EqualFold(targetInvited, p.Username) {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	var targetPlayer *Player
	for _, player := range p.Server.players {
		if strings.EqualFold(player.Username, targetInvited) {
			targetPlayer = player
			break
		}
	}

	if targetPlayer == nil {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeUnknownPlayer, protocol.MsgUnknownPlayer)
		return
	}

	if targetPlayer.GroupID != "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeAlreadyInGroup, protocol.MsgAlreadyInGroup)
		return
	}

	alreadyInvited := false
	for _, leader := range targetPlayer.GroupInvitations {
		if strings.EqualFold(leader, p.Username) {
			alreadyInvited = true
			break
		}
	}
	if !alreadyInvited {
		targetPlayer.GroupInvitations = append(targetPlayer.GroupInvitations, p.Username)
	}

	leader := strings.TrimPrefix(p.GroupID, "group_")
	targetPlayer.SendEvt(protocol.EvtGroupInvite, leader)

	p.Server.mu.Unlock()

	p.SendOK("")
}

func (p *Player) handleGroupJoin(args string) {
	targetLeader := strings.TrimSpace(args)
	if targetLeader == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()

	if p.GroupID != "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeAlreadyInGroup, protocol.MsgAlreadyInGroup)
		return
	}

	hasInvitation := false
	for _, leader := range p.GroupInvitations {
		if strings.EqualFold(leader, targetLeader) {
			hasInvitation = true
			break
		}
	}

	if !hasInvitation {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeBadRequest, protocol.MsgNoInvitation)
		return
	}

	var leaderPlayer *Player
	for _, player := range p.Server.players {
		if strings.EqualFold(player.Username, targetLeader) {
			leaderPlayer = player
			break
		}
	}

	if leaderPlayer == nil {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeUnknownPlayer, protocol.MsgUnknownPlayer)
		return
	}

	if leaderPlayer.GroupID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeBadRequest, protocol.MsgGroupNotFound)
		return
	}

	p.GroupID = leaderPlayer.GroupID

	updatedInvitations := make([]string, 0, len(p.GroupInvitations))
	for _, leader := range p.GroupInvitations {
		if !strings.EqualFold(leader, targetLeader) {
			updatedInvitations = append(updatedInvitations, leader)
		}
	}
	p.GroupInvitations = updatedInvitations

	for _, member := range p.Server.players {
		if member.Username != p.Username && member.GroupID == p.GroupID {
			member.SendEvt(protocol.EvtGroupJoin, p.Username)
		}
	}

	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("group=%s", p.GroupID))
}

// Retire le joueur de son groupe et dissout le groupe si c'est le leader.
func (p *Player) leaveGroupGetModif() {
	if p.GroupID == "" {
		return
	}

	isLeader := (p.GroupID == fmt.Sprintf("group_%s", p.Username))
	oldGroupID := p.GroupID
	p.GroupID = ""

	for _, member := range p.Server.players {
		if member.Username != p.Username && member.GroupID == oldGroupID {
			member.SendEvt(protocol.EvtGroupLeave, p.Username)
			if isLeader {
				member.GroupID = ""
			}
		}
	}
}

func (p *Player) handleGroupLeave(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()
	if p.GroupID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNotInGroup, protocol.MsgNotInGroup)
		return
	}

	p.leaveGroupGetModif()
	p.Server.mu.Unlock()

	p.SendOK("")
}
