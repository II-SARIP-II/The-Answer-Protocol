package server

import (
	"fmt"
	"log/slog"
	"strings"
	"tap/internal/protocol"
	"tap/internal/world"
)

const (
	FistFightDamage = 2
	NPCRespawnDelay = 50
)

func (p *Player) handleDamage(enemyID string, enemyData world.NPC) CombatRoundResult {
	damageDealt := p.getAttackDamage(p.EquippedWeapon)

	enemyData.Life -= damageDealt
	if enemyData.Life < 0 {
		enemyData.Life = 0
	}
	p.Server.world.NPCs[enemyID] = enemyData

	damageTaken := FistFightDamage
	if enemyData.Item != nil {
		damageTaken = p.getAttackDamage(*enemyData.Item)
	}

	p.HP -= damageTaken
	if p.HP < 0 {
		p.HP = 0
	}

	combatStatus := protocol.StatusCombat
	if enemyData.Life == 0 {
		combatStatus = protocol.StatusVictory
		p.Status = protocol.StatusVictory
		p.handleEnemyDeath(enemyID)
	} else if p.HP == 0 {
		combatStatus = protocol.StatusDefeated
		p.handlePlayerDeath(enemyData.Name)
	}

	combatRoundResult := CombatRoundResult{
		DamageDealt: damageDealt,
		DamageTaken: damageTaken,
		TargetHP:    enemyData.Life,
		Status:      combatStatus,
	}

	return combatRoundResult
}

func (p *Player) applyDefense(enemyData world.NPC) int {
	theoricalDamageTaken := FistFightDamage
	if enemyData.Item != nil {
		theoricalDamageTaken = p.getAttackDamage(*enemyData.Item)
	}

	damageTaken := 0
	switch p.DefendCount {
	case 0:
		damageTaken = 0
		p.DefendCount++
		p.Status = protocol.StatusPerfectDefense
	case 1:
		damageTaken = theoricalDamageTaken / 2
		p.DefendCount++
		p.Status = protocol.StatusWeakenedDefense
	default:
		damageTaken = theoricalDamageTaken
		p.Status = protocol.StatusBrokenDefense
	}

	p.HP -= damageTaken
	if p.HP < 0 {
		p.HP = 0
	}

	if p.HP == 0 {
		p.Status = "defeated"
		p.handlePlayerDeath(enemyData.Name)
	}
	return damageTaken
}

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

	p.broadcastRoom(p.CurrentRoom, protocol.EvtChatRoom, fmt.Sprintf("Server Player %s has eliminated %s!", p.Username, enemy.Name))
	p.broadcastGroup(protocol.EvtChatGroup, fmt.Sprintf("Server Player %s has eliminated %s!", p.Username, enemy.Name))

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
			damage := FistFightDamage
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

			p.SendEvt(protocol.EvtChatRoom, fmt.Sprintf("Server %s ambushes you! You take %d damage!", npc.Name, damage))
			p.broadcastRoom(roomID, protocol.EvtChatRoom, fmt.Sprintf("Server %s ambushes %s for %d damage!", npc.Name, p.Username, damage))

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
	p.broadcastRoom(roomWhenDeath, protocol.EvtChatRoom, fmt.Sprintf("Server Player %s has been defeated by %s!", p.Username, target))
	p.broadcastGroup(protocol.EvtChatGroup, fmt.Sprintf("Server Player %s was defeated by %s!", p.Username, target))

	p.CurrentRoom = "hotel_entrance"
	p.MaxHP /= 2
	if p.MaxHP < 25 {
		p.MaxHP = 25
	}
	p.HP = p.MaxHP
	p.Status = protocol.StatusDefeated
	data := fmt.Sprintf("Combat Defeat, defeated_player=%s, victorious_enemy=%s", p.Username, target)
	slog.Info(data)

	p.broadcastRoom(p.CurrentRoom, protocol.EvtRoomEnter, p.Username)
	p.SendEvt(protocol.EvtChatRoom, "Server Defeat: you were defeated and respawned in hotel_entrance!")
}

func (p *Player) getAttackDamage(weaponNameOrID string) int {
	if weaponNameOrID == "" {
		return FistFightDamage
	}
	for itemID, item := range p.Server.world.Items {
		if p.matchNameOrID(itemID, weaponNameOrID) {
			return item.Damage
		}
	}
	return FistFightDamage
}

func (p *Player) findAndValidateEnemy(target string) (string, world.NPC, bool) {
	cleanTarget := strings.TrimSpace(target)

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
		p.SendErr(protocol.CodeNPCNotFound, protocol.MsgNPCNotFound)
		return "", world.NPC{}, false
	}

	enemyData := p.Server.world.NPCs[foundID]
	if enemyData.Life <= 0 || (enemyData.Role != "enemy" && enemyData.Role != "target" && enemyData.Role != "passive enemy") {
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
				p.SendEvt(protocol.EvtChatRoom, fmt.Sprintf("Server [QUEST] %s progress: %d/%d",
					questID, questDataPlayer.CompletedSteps, questDataPlayer.Steps))
				if questDataPlayer.CompletedSteps >= questDataPlayer.Steps {
					giverName := "the quest giver"
					if giver, ok := p.Server.world.NPCs[worldQuest.Giver]; ok && giver.Name != "" {
						giverName = giver.Name
					}
					p.SendEvt(protocol.EvtChatRoom, fmt.Sprintf("Server [QUEST] %s: Objectives completed! Return to %s to claim your reward, saying \"quest\".", questID, giverName))
				}
				break
			}
		}
	}
}

func (p *Player) applyOpportunityAttack() bool {
	room := p.Server.world.Rooms[p.CurrentRoom]
	var attackingNPCName string
	fleeDamage := 0
	for _, npcID := range room.NPCs {
		npc := p.Server.world.NPCs[npcID]
		if npc.Life > 0 && (npc.Role == "enemy" || npc.Role == "target") {
			damage := FistFightDamage
			if npc.Item != nil {
				damage = p.getAttackDamage(*npc.Item)
			}
			fleeDamage = damage
			attackingNPCName = npc.Name
			break
		}
	}

	if attackingNPCName == "" {
		return true
	}

	p.HP -= fleeDamage
	if p.HP < 0 {
		p.HP = 0
	}

	p.SendEvt(protocol.EvtChatRoom, fmt.Sprintf("Server %s strikes you in the back as you flee! You take %d damage (%d HP left).", attackingNPCName, fleeDamage, p.HP))
	p.broadcastRoom(p.CurrentRoom, protocol.EvtChatRoom, fmt.Sprintf("Server %s strikes %s in the back as they flee for %d damage!", attackingNPCName, p.Username, fleeDamage))

	if p.HP == 0 {
		p.handlePlayerDeath(attackingNPCName)
		return false
	}

	return true
}
