package server

import (
	"fmt"
	"log/slog"
	"strings"
	"tap/internal/protocol"
	"time"
)

func (s *Server) scheduleNPCRespawn(enemyID string, roomID string) {

	time.AfterFunc(NPCRespawnDelay*time.Second, func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		npc := s.world.NPCs[enemyID]
		npc.Life = npc.MaxLife
		if enemyID == "chef" {
			npc.Role = "passive enemy"
		}
		s.world.NPCs[enemyID] = npc

		room := s.world.Rooms[roomID]
		room.NPCs = append(room.NPCs, enemyID)
		s.world.Rooms[roomID] = room

		slog.Info("NPC respawned", "npc", enemyID, "room", roomID)

		for _, player := range s.players {
			if strings.EqualFold(player.CurrentRoom, roomID) {
				player.SendEvt(protocol.EvtChatRoom, fmt.Sprintf("Server %s has reappeared!", npc.Name))
			}
		}
	})
}

func (p *Player) processTargetsRespawn(targets []string) {

	for _, target := range targets {
		foundIndex := -1
		for index, playerItem := range p.Inventory {
			if strings.EqualFold(target, playerItem) {
				foundIndex = index
				break
			}
		}
		if foundIndex == -1 {
			continue
		}

		if p.Inventory[foundIndex] == p.EquippedWeapon {
			p.EquippedWeapon = ""
		}

		p.Inventory = append(p.Inventory[:foundIndex], p.Inventory[foundIndex+1:]...)

		initialRoomID := p.Server.spawnsQuestsItems[target]
		roomData, exists := p.Server.world.Rooms[initialRoomID]
		if !exists {
			p.Server.mu.Unlock()
			p.SendErr(protocol.CodeNPCNotFound, protocol.MsgNPCNotFound)
			return
		}
		roomData.Items = append(roomData.Items, target)
		p.Server.world.Rooms[initialRoomID] = roomData
		p.Server.resetPassiveEnemy(initialRoomID, target)

		slog.Info("Item respawn", "item", target, "room", initialRoomID)
	}
}

func (s *Server) resetPassiveEnemy(roomID, itemID string) {
	room, exists := s.world.Rooms[roomID]
	if !exists {
		return
	}
	for _, npcID := range room.NPCs {
		npc := s.world.NPCs[npcID]
		if npc.Item != nil && *npc.Item == itemID {
			wasHostile := (npc.Role == "enemy")
			npc.Role = "passive enemy"
			npc.Life = npc.MaxLife
			s.world.NPCs[npcID] = npc
			slog.Info("Passive enemy calmed and reset", "npc", npcID, "item", itemID, "room", roomID)

			if wasHostile {
				for _, player := range s.players {
					if strings.EqualFold(player.CurrentRoom, roomID) {
						player.SendEvt(protocol.EvtChatRoom, fmt.Sprintf("Server %s has calmed down.", npc.Name))
					}
				}
			}
		}
	}
}
