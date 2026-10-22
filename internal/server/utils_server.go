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
