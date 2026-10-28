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

type CombatRoundResult struct {
	DamageDealt int
	DamageTaken int
	TargetHP    int
	Status      string
}

func (p *Player) handleAttack(args string) {
	if strings.TrimSpace(args) == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()
	defer p.Server.mu.Unlock()

	enemyID, enemyData, ok := p.findAndValidateEnemy(args)
	if !ok {
		return
	}

	p.Status = protocol.StatusCombat
	p.DefendCount = 0

	combatRoundResult := p.handleDamage(enemyID, enemyData)

	response := AttackResponse{
		AttackerHP: p.HP,
		TargetHP:   combatRoundResult.TargetHP,
		Damage:     combatRoundResult.DamageDealt,
		Status:     combatRoundResult.Status,
	}

	data, err := json.Marshal(response)
	if err != nil {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgBadRequest)
		return
	}

	slog.Info("Combat turn",
		"player", p.Username,
		"target", enemyID,
		"damage_dealt", combatRoundResult.DamageDealt,
		"damage_taken", combatRoundResult.DamageTaken,
		"player_hp", p.HP,
		"target_hp", combatRoundResult.TargetHP,
		"result", combatRoundResult.Status,
	)

	p.SendOK(string(data))
}

func (p *Player) handleDefend(args string) {
	p.Server.mu.Lock()
	defer p.Server.mu.Unlock()

	_, enemyData, ok := p.findAndValidateEnemy(args)
	if !ok {
		return
	}

	damageTaken := p.applyDefense(enemyData)

	response := DefendResponse{
		DefenderHP: p.HP,
		TargetHP:   enemyData.Life,
		Damage:     damageTaken,
		Status:     p.Status,
	}

	data, err := json.Marshal(response)
	if err != nil {
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

	p.SendOK(string(data))
}

func (p *Player) handleFlee(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}

	p.Server.mu.Lock()
	defer p.Server.mu.Unlock()

	room := p.Server.world.Rooms[p.CurrentRoom]
	availableDirections := make([]string, 0, len(room.Exits))
	for dir := range room.Exits {
		availableDirections = append(availableDirections, dir)
	}

	if len(availableDirections) == 0 {
		p.SendErr(protocol.CodeNoExit, protocol.MsgNoExit)
		return
	}

	if !p.applyOpportunityAttack() {
		p.SendOK(fmt.Sprintf("room=%s", p.CurrentRoom))
		return
	}

	randomDirection := availableDirections[rand.IntN(len(availableDirections))]
	exit := room.Exits[randomDirection]
	newRoom := exit.Target
	oldRoom := p.CurrentRoom

	p.broadcastRoom(oldRoom, protocol.EvtRoomLeave, p.Username)
	p.broadcastRoom(oldRoom, protocol.EvtChatRoom, fmt.Sprintf("Server Player %s fled towards %s!", p.Username, randomDirection))
	p.broadcastRoom(newRoom, protocol.EvtRoomEnter, p.Username)

	p.CurrentRoom = newRoom
	p.Status = protocol.StatusHealthy
	p.DefendCount = 0

	slog.Info("Player fled combat",
		"player", p.Username,
		"from_room", oldRoom,
		"to_room", newRoom,
		"direction", randomDirection,
		"hp_remaining", p.HP,
	)

	p.SendOK(fmt.Sprintf("room=%s", newRoom))
}
