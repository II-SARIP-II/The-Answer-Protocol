package server

import (
	"strings"
)

// matchNameOrID checks if a world entity matches user input by canonical ID or display name.
func (p *Player) matchNameOrID(worldEntity, userInput string) bool {
	if strings.EqualFold(worldEntity, userInput) {
		return true
	}
	id1 := p.resolveID(worldEntity)
	id2 := p.resolveID(userInput)
	return strings.EqualFold(id1, id2)
}

// resolveID maps a display name or ID to its canonical identifier.
func (p *Player) resolveID(input string) string {
	for id, itemData := range p.Server.world.Items {
		if strings.EqualFold(input, id) || strings.EqualFold(input, itemData.Name) {
			return id
		}
	}
	for id, npcData := range p.Server.world.NPCs {
		if strings.EqualFold(input, id) || strings.EqualFold(input, npcData.Name) {
			return id
		}
	}
	return input
}

func (p *Player) broadcastRoom(roomID, eventType, message string) {
	for _, other := range p.Server.players {
		if other.Username != p.Username && strings.EqualFold(other.CurrentRoom, roomID) {
			other.SendEvt(eventType, message)
		}
	}
}

func (p *Player) broadcastGroup(eventType, message string) {
	if p.GroupID == "" {
		return
	}
	for _, other := range p.Server.players {
		if other.Username != p.Username && other.GroupID == p.GroupID {
			other.SendEvt(eventType, message)
		}
	}
}

func (p *Player) broadcastGlobal(eventType, message string) {
	for _, other := range p.Server.players {
		other.SendEvt(eventType, message)
	}
}
