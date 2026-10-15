package server

import (
	"encoding/json"
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
		if p.matchNameOrID(item, targetItem) {
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
		if p.matchNameOrID(item, targetItem) {
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

	p.SendOK(fmt.Sprintf("dropped=%s", droppedItem))
}

func (p *Player) handleInventory(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()
	inventory := p.Inventory
	if inventory == nil {
		inventory = []string{}
	}
	p.Server.mu.Unlock()

	data, err := json.Marshal(inventory)
	if err != nil {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
		return
	}
	p.SendOK(string(data))
}

func (p *Player) handleTalk(args string) {
	targetNpc := strings.TrimSpace(args)
	if targetNpc == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()

	var foundNPCID string
	room := p.Server.world.Rooms[p.CurrentRoom]
	for _, npcID := range room.NPCs {
		if p.matchNameOrID(npcID, targetNpc) {
			foundNPCID = npcID
			break
		}
	}

	if foundNPCID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotFound, protocol.MsgNPCNotFound)
		return
	}

	npcData, exists := p.Server.world.NPCs[foundNPCID]
	if !exists || len(npcData.Dialogue) == 0 {
		p.Server.mu.Unlock()
		p.SendOK(fmt.Sprintf("%s has nothing to say.", npcData.Name))
		return
	}

	npcReply := npcData.Dialogue[0]
	npcData.Dialogue = append(npcData.Dialogue[1:], npcReply)
	p.Server.world.NPCs[foundNPCID] = npcData

	p.Server.mu.Unlock()

	p.SendOK(npcReply)
}

func (p *Player) handleAttack(args string) {
	target := strings.TrimSpace(args)
	if target == "" {
		if args == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
		}
	}

	p.Server.mu.Lock()

	var foundTargetID string
	for _, targetID := range p.Server.world.Rooms[p.CurrentRoom].NPCs {
		if p.matchNameOrID(targetID, target) {
			foundTargetID = targetID
			break
		}
	}

	if foundTargetID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotFound, protocol.MsgNPCNotFound)
		return
		}

	enemyData := p.Server.world.NPCs[foundTargetID]
	// on peut bien attaquer les "passive enemies"??
	if strings.ToLower(enemyData.Role) != "enemy" && strings.ToLower(enemyData.Role) != "passive enemy" {
		p.SendErr(protocol.CodeNPCNotHostile, protocol.MsgNPCNotHostile)
		return
	}

	// Calcul des dégâts infligés par le joueur :
	//    Basé sur l'arme dans l'inventaire du joueur (ex: Colt 1911 -> 4 dégâts) ou poings (ex: 2) si pas d'arme?
	//    ajouter une commande WEAPON pour choisir l árme dans l'inventory?
	// Application des dégâts au PNJ : npc.Life -= degats
	//              * Le PNJ survit (npc.Life > 0) :
	//              		Le PNJ contre-attaque avec son arme (ex: M16 -> 7 dégâts infligés au joueur)
	//              		status = "combat"
	//              * Le PNJ meurt (npc.Life <= 0) :
	//              		status = "victory" (ou "defeated")
	//              		L'objet du PNJ tombe au sol dans la pièce (Loot !)
	//              		Le PNJ est retiré de room.NPCs
	// Réponse JSON au client :
	//    OK {"attacker_hp": 93, "target_hp": 6, "damage": 4, "status": "combat"}// TODO
	// si le joueur meurt: (cf sujet)
	// 			Players with 0 HP must respawn at a safe location with reduced health
	// ◦ Combat results must be logged and broadcast to relevant players
	// ◦ Design Choice: Your group must design and implement turn-based combat
	// mechanics, including damage formulas, initiative order, and additional combat
	// commands (DEFEND, FLEE, etc.). Document and justify your choices in the
	// README.

	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("Received ATTACK from %s. (TODO)", foundTargetID))
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
