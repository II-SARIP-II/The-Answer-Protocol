package server

import (
	"fmt"
	"log/slog"
	"strings"
	"tap/internal/protocol"
	"time"
)

func (s *Server) scheduleNPCRespawn(enemyID string, roomID string) {

	time.AfterFunc(protocol.NPCRespawnDelay*time.Second, func() {
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
				player.SendEvt("ROOM CHAT", fmt.Sprintf("Server %s has reappeared!", npc.Name))
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

		slog.Info("Item respawn", "item", target, "room", initialRoomID)
	}
}
