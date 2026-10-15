package server

import (
	"fmt"
	"strings"
	"tap/internal/protocol"
)

func (p *Player) handleTake(args string) {
	targetItem := strings.TrimSpace(args)
	if targetItem == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()
	room, exists := p.Server.world.Rooms[p.CurrentRoom]
	if !exists {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
		return
	}

	// L'index -1 est une convention signifiant : « pas encore trouvé ».
	foundIndex := -1
	var takenItem string
	for i, item := range room.Items {
		if strings.EqualFold(item, targetItem) {
			foundIndex = i
			takenItem = item
			break
		}
	}

	if foundIndex == -1 {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeItemNotFound, protocol.MsgItemNotFound)
		return
	}

	p.Inventory = append(p.Inventory, takenItem)
	// pour remove en go, on append les deux parties du tableau, avant et apr`s le foundIndex
	room.Items = append(room.Items[:foundIndex], room.Items[foundIndex+1:]...)
	p.Server.world.Rooms[p.CurrentRoom] = room
	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("taken=%s", takenItem))
}

func (p *Player) handleDrop(args string) {
	targetItem := strings.TrimSpace(args)
	if targetItem == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()

	foundIndex := -1
	var droppedItem string
	for i, item := range p.Inventory {
		if strings.EqualFold(targetItem, item) {
			foundIndex = i
			droppedItem = item
			break
		}
	}

	if foundIndex == -1 {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeItemNotInInventory, protocol.MsgItemNotInInventory)
		return
	}

	room := p.Server.world.Rooms[p.CurrentRoom]
	room.Items = append(room.Items, droppedItem)
	p.Server.world.Rooms[p.CurrentRoom] = room

	p.Inventory = append(p.Inventory[:foundIndex], p.Inventory[foundIndex+1:]...)
	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("droppen=%s", droppedItem))
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
