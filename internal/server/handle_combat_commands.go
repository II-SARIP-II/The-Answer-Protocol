package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"tap/internal/protocol"
	"tap/internal/world"
	"math/rand/v2"
)

type AttackResponse struct {
	AttackerHP int    `json:"attacker_hp"`
	TargetHP   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}


type DefendResponse struct {
	DefenderHP int    `json:"defender_hp"`
	TargetHP   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}

func (p *Player) findAndValidateEnemy(target string) (string, world.NPC, bool) {
	cleanTarget := strings.TrimSpace(target)
	if cleanTarget == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return "", world.NPC{}, false
	}

	p.Server.mu.Lock()

	var foundID string
	for _, targetID := range p.Server.world.Rooms[p.CurrentRoom].NPCs {
		if p.matchNameOrID(targetID, cleanTarget) {
			foundID = targetID
			break
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

func (p *Player) handleAttack(args string) {
	enemyID, enemyData, ok := p.findAndValidateEnemy(args)
	if !ok {
		return
	}

	p.Status = protocol.StatusCombat
	p.DefendCount = 0

	damageDealt := protocol.FistFightDamage
	if len(p.Inventory) == 1 {
		p.EquippedWeapon = p.Inventory[0]
	}
	if p.EquippedWeapon != "" {
		for itemID, item := range p.Server.world.Items {
			if p.matchNameOrID(itemID, p.EquippedWeapon) {
				damageDealt = item.Damage
				break
			}
		}
	}

	enemyData.Life -= damageDealt
	if enemyData.Life < 0 {
		enemyData.Life = 0
	}
	p.Server.world.NPCs[enemyID] = enemyData


	damageTaken := 0
	if enemyData.Life > 0 {
		damageTaken = protocol.FistFightDamage
		if enemyData.Item != nil {
			for itemID, item := range p.Server.world.Items {
				if p.matchNameOrID(itemID, *enemyData.Item) {
					damageTaken = item.Damage
					break
				}
			}
		}

		p.HP -= damageTaken
		if p.HP < 0 {
			p.HP = 0
		}
	}

	combatStatus := protocol.StatusCombat
	if enemyData.Life == 0 {
		combatStatus = "victory"
		p.Status = protocol.StatusHealthy
		p.handleEnemyDeath(enemyID)
	} else if p.HP == 0 {
		combatStatus = "defeated"
		p.handlePlayerDeath(enemyData.Name)
	}

	response := AttackResponse{
		AttackerHP: p.HP,
		TargetHP:   enemyData.Life,
		Damage:     damageDealt,
		Status:     combatStatus,
	}

	data, err := json.Marshal(response)
	if err != nil {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
		return
	}

	slog.Info("Combat turn",
		"player", p.Username,
		"target", enemyID,
		"damage_dealt", damageDealt,
		"damage_taken", damageTaken,
		"player_hp", p.HP,
		"target_hp", enemyData.Life,
		"result", combatStatus,
	)

	p.Server.mu.Unlock()
	p.SendOK(string(data))
}


func (p *Player) handleEnemyDeath(enemyID string) {
	currentRoom := p.Server.world.Rooms[p.CurrentRoom]
	enemy := p.Server.world.NPCs[enemyID]

	p.SendEvt("Victory:", fmt.Sprintf("You have eliminated %s !", enemy.Name))
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

	for _, otherPlayer := range p.Server.players {
		if otherPlayer.Username != p.Username && otherPlayer.CurrentRoom == p.CurrentRoom {
			eventData := fmt.Sprintf("Player %s has eliminated the enemy %s!", p.Username, enemy.Name)
			otherPlayer.SendEvt(protocol.ChatRoom+ " " +protocol.CmdChat, eventData)
		}
	}

	if p.GroupID != "" {
		for _, otherPlayer := range p.Server.players {
			if otherPlayer.Username != p.Username && otherPlayer.GroupID == p.GroupID && otherPlayer.CurrentRoom != p.CurrentRoom {
				eventData := fmt.Sprintf("Player %s has eliminated the enemy %s!", p.Username, enemy.Name)
				otherPlayer.SendEvt(protocol.ChatGroup+ " " +protocol.CmdChat, eventData)
			}
		}
	}
}

func (p *Player) handlePlayerDeath(target string) {
	roomWhenDeath := p.CurrentRoom
	p.SendEvt("Defeat:", fmt.Sprintf("%s killed you, you will spawn in room0 !", target))
	for _, otherPlayer := range p.Server.players {
		if otherPlayer.Username != p.Username && strings.EqualFold(otherPlayer.CurrentRoom, roomWhenDeath) {
			otherPlayer.SendEvt(protocol.EvtRoomLeave, p.Username)
			otherPlayer.SendEvt("ROOM CHAT", fmt.Sprintf("Player %s has been defeated by %s!", p.Username, target))
		}
	}

	if p.GroupID != "" {
		for _, otherPlayer := range p.Server.players {
			if otherPlayer.Username != p.Username && otherPlayer.GroupID == p.GroupID {
				otherPlayer.SendEvt("GROUP CHAT", fmt.Sprintf("Player %s was defeated by %s!", p.Username, target))
			}
		}
	}

	p.CurrentRoom = "room0"
	p.MaxHP /= 2
	if p.MaxHP < 25 {
		p.MaxHP = 25
	}
	p.HP = p.MaxHP
	p.Status = protocol.StatusDefeated

	for _, otherPlayer := range p.Server.players {
		if otherPlayer.Username != p.Username && strings.EqualFold(otherPlayer.CurrentRoom, p.CurrentRoom) {
			otherPlayer.SendEvt(protocol.EvtRoomEnter, p.Username)
		}
	}
}

func (p *Player) handleDefend(args string) {
	_, enemyData, ok := p.findAndValidateEnemy(args)
	if !ok {
		return
	}

	p.Status = protocol.StatusPerfectDefense

	theoricalDamageTaken := protocol.FistFightDamage
	if enemyData.Item != nil {
		for itemID, item := range p.Server.world.Items {
			if p.matchNameOrID(itemID, *enemyData.Item) {
				theoricalDamageTaken = item.Damage
				break
			}
		}
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

	response := DefendResponse{
		DefenderHP: p.HP,
		TargetHP:   enemyData.Life,
		Damage:     damageTaken,
		Status:     p.Status,
	}

	data, err := json.Marshal(response)
	if err != nil {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
		return
	}

	slog.Info("Defense turn",
		"player", p.Username,
		"target", enemyData.Name,
		"damage_taken", damageTaken,
		"player_hp", p.HP,
		"target_hp", enemyData.Life,
		"status", p.Status,
	)

	p.Server.mu.Unlock()
	p.SendOK(string(data))
}

func (p *Player) handleFlee(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()

	room := p.Server.world.Rooms[p.CurrentRoom]
	availableDirections := make([]string, 0, len(room.Exits))
	for dir := range room.Exits {
		availableDirections = append(availableDirections, dir)
	}

	if len(availableDirections) == 0 {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNoExit, protocol.MsgNoExit)
		return
	}

	randomDirection := availableDirections[rand.IntN(len(availableDirections))]
	exit := room.Exits[randomDirection]
	newRoom := exit.Target
	oldRoom := p.CurrentRoom

	for _, otherPlayer := range p.Server.players {
		if otherPlayer.Username != p.Username && strings.EqualFold(otherPlayer.CurrentRoom, oldRoom) {
			otherPlayer.SendEvt(protocol.EvtRoomLeave, p.Username)
			otherPlayer.SendEvt("ROOM CHAT", fmt.Sprintf("Server Player %s fled the combat towards %s!", p.Username, randomDirection))
		}
	}

	p.CurrentRoom = newRoom
	p.Status = protocol.StatusHealthy
	p.DefendCount = 0

	for _, otherPlayer := range p.Server.players {
		if otherPlayer.Username != p.Username && strings.EqualFold(otherPlayer.CurrentRoom, p.CurrentRoom) {
			otherPlayer.SendEvt(protocol.EvtRoomEnter, p.Username)
		}
	}

	slog.Info("Player fled combat",
		"player", p.Username,
		"from_room", oldRoom,
		"to_room", newRoom,
		"direction", randomDirection,
	)

	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("room=%s", newRoom))
}


