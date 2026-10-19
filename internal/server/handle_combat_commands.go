package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"tap/internal/protocol"
)

type AttackResponse struct {
	AttackerHP int    `json:"attacker_hp"`
	TargetHP   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}

func (p *Player) handleAttack(args string) {
	target := strings.TrimSpace(args)
	if target == "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
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
	if enemyData.Life <= 0 || (enemyData.Role != "enemy" && enemyData.Role != "target") {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotHostile, protocol.MsgNPCNotHostile)
		return
	}

	p.Status = protocol.StatusCombat

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
	p.Server.world.NPCs[foundTargetID] = enemyData

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
		p.handleEnemyDeath(foundTargetID)
	} else if p.HP == 0 {
		combatStatus = "defeated"
		p.handlePlayerDeath()
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
		"target", foundTargetID,
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
			otherPlayer.SendEvt(protocol.ChatRoom + protocol.CmdChat, eventData)
		}
	}

	if p.GroupID != ""{
		for _, otherPlayer := range p.Server.players {
			if otherPlayer.Username != p.Username && otherPlayer.GroupID == p.GroupID && otherPlayer.CurrentRoom != p.CurrentRoom {
				eventData := fmt.Sprintf("Player %s has eliminated the enemy %s!", p.Username, enemy.Name)
				otherPlayer.SendEvt(protocol.ChatGroup + protocol.CmdChat, eventData)
			}
		}
	}
}

func (p *Player) handlePlayerDeath(){
	roomWhenDeath := p.CurrentRoom

	for _, otherPlayer := range p.Server.players {
		if otherPlayer.Username != p.Username && strings.EqualFold(otherPlayer.CurrentRoom, roomWhenDeath) {
			otherPlayer.SendEvt(protocol.EvtRoomLeave, p.Username)
		}
	}

	p.CurrentRoom = "room0"
	p.MaxHP /= 2
	if p.MaxHP < 25 {
		p.MaxHP = 25
	}
	p.HP = p.MaxHP
	p.Status = protocol.StatusHealthy

	// TODO ajouter sendEVT aux joueurs de la pièce
}


// ◦ Combat results must be logged and broadcast to relevant players TODO
// ◦ Design Choice: Your group must design and implement turn-based combat
// mechanics, including damage formulas, initiative order, and additional combat
// commands (DEFEND, FLEE, etc.). Document and justify your choices in the
// README.

func (p *Player) handleDefend(args string) {
	target := strings.TrimSpace(args)
	if target == "" {
		if args == "" {
			p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
			return
		}
	}
	p.SendOK("Received DEFEND. (TODO)")
	// si player envoie DEFEND, 0 damage recu/donné
	// au 2ème DEFEND, 1/2 damage recu/100% donné
	// si MOVE avant 2eme annulation du mode defense
}

func (p *Player) handleFlee(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	p.SendOK("Received FLEE. (TODO)")
	// si fuite au lieu de ATTACK, pas de damage mais Move par une exit au hasard
}
