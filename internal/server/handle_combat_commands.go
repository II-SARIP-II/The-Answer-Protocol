package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"tap/internal/protocol"
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

func (p *Player) handleAttack(args string) {
	if strings.TrimSpace(args) == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	enemyID, enemyData, ok := p.findAndValidateEnemy(args)
	if !ok {
		return
	}

	p.Status = protocol.StatusCombat
	p.DefendCount = 0

	damageDealt := p.getAttackDamage(p.EquippedWeapon)

	enemyData.Life -= damageDealt
	if enemyData.Life < 0 {
		enemyData.Life = 0
	}
	p.Server.world.NPCs[enemyID] = enemyData

	damageTaken := 0
	if enemyData.Life > 0 && enemyData.Item != nil {
		damageTaken = p.getAttackDamage(*enemyData.Item)
	}

	p.HP -= damageTaken
	if p.HP < 0 {
		p.HP = 0
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

func (p *Player) handleDefend(args string) {
	_, enemyData, ok := p.findAndValidateEnemy(args)
	if !ok {
		return
	}

	p.Status = protocol.StatusPerfectDefense

	theoricalDamageTaken := protocol.FistFightDamage
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

	p.broadcastRoom(oldRoom, protocol.EvtRoomLeave, p.Username)
	p.broadcastRoom(oldRoom, "ROOM CHAT", fmt.Sprintf("Server Player %s fled towards %s!", p.Username, randomDirection))
	p.broadcastRoom(newRoom, protocol.EvtRoomEnter, p.Username)

	p.CurrentRoom = newRoom
	p.Status = protocol.StatusHealthy
	p.DefendCount = 0

	slog.Info("Player fled combat",
		"player", p.Username,
		"from_room", oldRoom,
		"to_room", newRoom,
		"direction", randomDirection,
	)

	p.Server.mu.Unlock()

	p.SendOK(fmt.Sprintf("room=%s", newRoom))
}
