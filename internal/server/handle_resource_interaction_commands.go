package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"tap/internal/protocol"
	"tap/internal/world"
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
	room.Items = append(room.Items[:foundIndex], room.Items[foundIndex+1:]...)
	p.Server.world.Rooms[p.CurrentRoom] = room

	var provokedPassiveEnemy *world.NPC
	var damage int
	for _, npcID := range room.NPCs {
		npc := p.Server.world.NPCs[npcID]
		if npc.Life > 0 && npc.Role == "passive enemy" {
			npc.Role = "enemy"
			p.Server.world.NPCs[npcID] = npc

			damage = protocol.FistFightDamage
			if npc.Item != nil {
				damage = p.getAttackDamage(*npc.Item)
			}

			p.HP -= damage
			if p.HP < 0 {
				p.HP = 0
			}
			p.Status = protocol.StatusCombat
			p.DefendCount = 0

			slog.Info("Passive enemy provoked by TAKE",
				"npc", npc.Name,
				"player", p.Username,
				"item", takenItem,
				"damage", damage,
				"player_hp", p.HP,
			)

			provokedPassiveEnemy = &npc
			break
		}
	}

	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("taken=%s", takenItem))

	if provokedPassiveEnemy != nil {
		p.SendEvt("ROOM CHAT", fmt.Sprintf("%s is enraged: \"No one has the right to touch my knives!\" He attacks you for %d damage!", provokedPassiveEnemy.Name, damage))
		p.broadcastRoom(p.CurrentRoom, "ROOM CHAT", fmt.Sprintf("%s is enraged that %s took %s! He attacks for %d damage!", provokedPassiveEnemy.Name, p.Username, takenItem, damage))

		if p.HP == 0 {
			p.Server.mu.Lock()
			p.handlePlayerDeath(provokedPassiveEnemy.Name)
			p.Server.mu.Unlock()
		}
	}
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

	if droppedItem == p.EquippedWeapon {
		p.EquippedWeapon = ""
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

func (p *Player) handleEquip(args string) {
	targetWeapon := strings.TrimSpace(args)
	if targetWeapon == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()
	foundWeapon := ""

	for _, weapon := range p.Inventory {
		if p.matchNameOrID(weapon, targetWeapon) {
			foundWeapon = weapon
			p.EquippedWeapon = foundWeapon
			break
		}
	}

	if foundWeapon == "" {
		p.SendErr(protocol.CodeItemNotInInventory, protocol.MsgItemNotInInventory)
		p.Server.mu.Unlock()
		return
	}
	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("equipped=%s", foundWeapon))
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

type StatusResponse struct {
	HP     int    `json:"hp"`
	MaxHP  int    `json:"max_hp"`
	Status string `json:"status"`
}

func (p *Player) handleStatus(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	response := StatusResponse{
		HP:     p.HP,
		MaxHP:  p.MaxHP,
		Status: p.Status,
	}

	data, err := json.Marshal(response)
	if err != nil {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
	}

	p.SendOK(string(data))
}

func (p *Player) handleQuest(args string) {
	targetNpc := strings.TrimSpace(args)

	p.Server.mu.Lock()

	room := p.Server.world.Rooms[p.CurrentRoom]
	var questGiverID string

	if targetNpc == "" {
		for _, npcID := range room.NPCs {
			npc := p.Server.world.NPCs[npcID]
			if npc.Role == "quest-giver" {
				questGiverID = npcID
				break
			}
		}
	} else {
		for _, npcID := range room.NPCs {
			if p.matchNameOrID(npcID, targetNpc) {
				questGiverID = npcID
				break
			}
		}
	}

	if questGiverID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotFound, protocol.MsgNPCNotFound)
		return
	}

	questGiverData := p.Server.world.NPCs[questGiverID]
	if questGiverData.Role != "quest-giver" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNoQuestAvailable, protocol.MsgNoQuestAvailable)
		return
	}

	if p.Server.world.Quests == nil {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNoQuestAvailable, protocol.MsgNoQuestAvailable)
		return
	}


	// // si prerequisite: recuperer le statut du player pour cette quête et verifiee qu elle est
	// // accomplie, sinon continue
	// // si rien à la fin, message pas de quête disponible
	// playerNextQuest := Quest {}
	// for quest := range p.Server.world.Quests {
	// 	if quest.prerequisite == "" {
	// 		playerNextQuest = quest
	// 		break
	// 	} else {
	// 		completedQuest := false
	// 		for playerQuest := range p.Quests {
	// 			if p.Quests[quest].Status == "completed" && p.Quests[quest].ID == quest {
	// 				completedQuest := true
	// 				break
	// 			}
	// 		if !completedQuest {
	// 			p.Server.mu.Unlock()
	// 			p.SendErr(protocol.CodeNoQuestAvailable, protocol.MsgNoQuestAvailable)
	// 			return
	// 			}
	// 		playerNextQuest = quest
	// 		}
	// 	}


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

// il fait QUEST pour récupérer sa quête (status = "active", progerssion = (0/2))
// ( qui est intégrée à la structure du player sous forme dúne liste de quêtes,
// visible avec QUESTS. La quête ne peut être attribuée au player que s'il a les prérequis
//  On met en place en compteur pour les quêtes en plusieurs étapes.
// Une fois la quête accomplie, il retourne auprès du donneur d'ordre
// pour lui siginifier l'accoplissement, progression == count/count
// ca fait status à "completed"
// Il recoit alors sa récompense. PLusieurs joueurs peuvent chercher à accomplir une même quête.
// Une fois que la quête ultime est réalisée, message de Victoire totale donnée par le donneur de quête.
//  TT est broadcasté et loggé.
