package server

import (
	"fmt"
	"log/slog"
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

	for _, otherPlayer := range p.Server.players {
		if otherPlayer.CurrentRoom == p.CurrentRoom && otherPlayer.Status == "combat" && otherPlayer.Username != p.Username {
			otherPlayer.Status = "healthy"
		}
	}

	p.Server.world.Rooms[p.CurrentRoom] = currentRoom

	p.broadcastRoom(p.CurrentRoom, "ROOM CHAT", fmt.Sprintf("Player %s has eliminated %s!", p.Username, enemy.Name))
	p.broadcastGroup("GROUP CHAT", fmt.Sprintf("Player %s has eliminated %s!", p.Username, enemy.Name))

	p.processDefeatQuest(enemyID)

	p.Server.scheduleNPCRespawn(enemyID, p.CurrentRoom)

	data := fmt.Sprintf("Combat victory, victorious_player=%s, defeated_enemy=%s", p.Username, enemy.Name)
	slog.Info(data)
}

func (p *Player) checkHostileAmbush(roomID string) {
	room := p.Server.world.Rooms[roomID]
	for _, npcID := range room.NPCs {
		npc := p.Server.world.NPCs[npcID]
		if npc.Life > 0 && (npc.Role == "enemy" || npc.Role == "target") {
			damage := protocol.FistFightDamage
			if npc.Item != nil {
				damage = p.getAttackDamage(*npc.Item)
			}

			p.HP -= damage
			if p.HP < 0 {
				p.HP = 0
			}

			p.Status = protocol.StatusCombat
			p.DefendCount = 0

			slog.Info("Hostile NPC ambush",
				"npc", npc.Name,
				"player", p.Username,
				"room", roomID,
				"damage", damage,
				"player_hp", p.HP,
			)

			p.SendEvt("ROOM CHAT", fmt.Sprintf("%s ambushes you! You take %d damage!", npc.Name, damage))
			p.broadcastRoom(roomID, "ROOM CHAT", fmt.Sprintf("%s ambushes %s for %d damage!", npc.Name, p.Username, damage))

			if p.HP == 0 {
				p.handlePlayerDeath(npc.Name)
			}
			break
		}
	}
}

func (p *Player) handlePlayerDeath(target string) {
	roomWhenDeath := p.CurrentRoom

	p.broadcastRoom(roomWhenDeath, protocol.EvtRoomLeave, p.Username)
	p.broadcastRoom(roomWhenDeath, "ROOM CHAT", fmt.Sprintf("Player %s has been defeated by %s!", p.Username, target))
	p.broadcastGroup("GROUP CHAT", fmt.Sprintf("Player %s was defeated by %s!", p.Username, target))

	p.CurrentRoom = "room0"
	p.MaxHP /= 2
	if p.MaxHP < 25 {
		p.MaxHP = 25
	}
	p.HP = p.MaxHP
	p.Status = protocol.StatusDefeated
	data := fmt.Sprintf("Combat Defeat, defeated_player=%s, victorious_enemy=%s", p.Username, target)
	slog.Info(data)

	p.broadcastRoom(p.CurrentRoom, protocol.EvtRoomEnter, p.Username)
	p.SendEvt("ROOM CHAT", "Defeat: you were defeated and respawned in room0!")
}

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
		// Pas de cible précisée : prend le premier ennemi hostile de la pièce
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
	if enemyData.Life <= 0 || (enemyData.Role != "enemy" && enemyData.Role != "target" && enemyData.Role != "passive enemy") {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotHostile, protocol.MsgNPCNotHostile)
		return "", world.NPC{}, false
	}

	if enemyData.Role == "passive enemy" {
		enemyData.Role = "enemy"
		p.Server.world.NPCs[foundID] = enemyData
	}

	return foundID, enemyData, true
}

func (p *Player) processDefeatQuest(enemyID string) {
	for questID, questDataPlayer := range p.Quests {
		if questDataPlayer.Status != "active" {
			continue
		}
		worldQuest, exists := p.Server.world.Quests[questID]
		if !exists || worldQuest.Type != "defeat" {
			continue
		}
		if questDataPlayer.CompletedSteps >= questDataPlayer.Steps {
			continue
		}
		for _, target := range worldQuest.Targets {
			if p.matchNameOrID(target, enemyID) {
				questDataPlayer.CompletedSteps++
				p.Quests[questID] = questDataPlayer
				slog.Info("Quest progress",
					"player", p.Username,
					"quest", questID,
					"progress", fmt.Sprintf("%d/%d", questDataPlayer.CompletedSteps, questDataPlayer.Steps),
				)
				p.SendEvt("ROOM CHAT", fmt.Sprintf("Server [QUEST] %s progress: %d/%d",
					questID, questDataPlayer.CompletedSteps, questDataPlayer.Steps))
				if questDataPlayer.CompletedSteps >= questDataPlayer.Steps {
					giverName := "the quest giver"
					if giver, ok := p.Server.world.NPCs[worldQuest.Giver]; ok && giver.Name != "" {
						giverName = giver.Name
					}
					p.SendEvt("ROOM CHAT", fmt.Sprintf("Server [QUEST] %s: Objectives completed! Return to %s to claim your reward, saying \"quest\".", questID, giverName))
				}
				break
			}
		}
	}
}
