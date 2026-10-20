package server

import (
	"fmt"
	"strings"
	"tap/internal/protocol"
	"tap/internal/world"
)

func (p *Player) handleEnemyDeath(enemyID string) {
	currentRoom := p.Server.world.Rooms[p.CurrentRoom]
	enemy := p.Server.world.NPCs[enemyID]

	UpdatedNPCsInRoom := make([]string, 0, len(currentRoom.NPCs))
	for _, id := range currentRoom.NPCs {
		if id != enemyID {
			UpdatedNPCsInRoom = append(UpdatedNPCsInRoom, id)
		}
	}
	currentRoom.NPCs = UpdatedNPCsInRoom

	if enemy.Item != nil && *enemy.Item != "" {
		itemID := p.resolveID(*enemy.Item)
		currentRoom.Items = append(currentRoom.Items, itemID)
	}
	p.Server.world.Rooms[p.CurrentRoom] = currentRoom

	p.broadcastRoom(p.CurrentRoom, "ROOM CHAT", fmt.Sprintf("Player %s has eliminated %s!", p.Username, enemy.Name))
	p.broadcastGroup("GROUP CHAT", fmt.Sprintf("Player %s has eliminated %s!", p.Username, enemy.Name))
}



func (p *Player) handlePlayerDeath(target string) {
	roomWhenDeath := p.CurrentRoom

	p.broadcastRoom(roomWhenDeath, protocol.EvtRoomLeave, p.Username)
	p.broadcastRoom(roomWhenDeath, "ROOM CHAT", fmt.Sprintf("Server Player %s has been defeated by %s!", p.Username, target))
	p.broadcastGroup("GROUP CHAT", fmt.Sprintf("Server Player %s was defeated by %s!", p.Username, target))

	p.CurrentRoom = "room0"
	p.MaxHP /= 2
	if p.MaxHP < 25 {
		p.MaxHP = 25
	}
	p.HP = p.MaxHP
	p.Status = protocol.StatusDefeated

	p.broadcastRoom(p.CurrentRoom, protocol.EvtRoomEnter, p.Username)
}


// getWeaponDamage retourne les dégâts d'une arme (ou les dégâts à mains nues si vide/introuvable)
func (p *Player) getAttackDamage(weaponNameOrID string) int {
	if weaponNameOrID == "" {
		return protocol.FistFightDamage
	}
	for itemID, item := range p.Server.world.Items {
		if p.matchNameOrID(itemID, weaponNameOrID) {
			return item.Damage
		}
	}
	return protocol.FistFightDamage
}


func (p *Player) findAndValidateEnemy(target string) (string, world.NPC, bool) {
	cleanTarget := strings.TrimSpace(target)

	p.Server.mu.Lock()

	var foundID string

	if cleanTarget == "" {
		// Pas de cible précisée (ex: DEFEND tout seul) : prend le premier ennemi hostile de la pièce
		for _, npcID := range p.Server.world.Rooms[p.CurrentRoom].NPCs {
			npc := p.Server.world.NPCs[npcID]
			if npc.Life > 0 && (npc.Role == "enemy" || npc.Role == "target") {
				foundID = npcID
				break
			}
		}
	} else {
		// Cible précisée pour attack ou defense
		for _, targetID := range p.Server.world.Rooms[p.CurrentRoom].NPCs {
			if p.matchNameOrID(targetID, cleanTarget) {
				foundID = targetID
				break
			}
		}
	}

	if foundID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotFound, protocol.MsgNPCNotFound)
		return "", world.NPC{}, false
	}

	enemyData := p.Server.world.NPCs[foundID]
	if enemyData.Life <= 0 || (enemyData.Role != "enemy" && enemyData.Role != "target") {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotHostile, protocol.MsgNPCNotHostile)
		return "", world.NPC{}, false
	}

	return foundID, enemyData, true
}
