package server

import (
	"fmt"
	"strings"
	"tap/internal/protocol"
)

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
	//    ajouter une commande EQUIP pour choisir l árme dans l'inventory??
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
