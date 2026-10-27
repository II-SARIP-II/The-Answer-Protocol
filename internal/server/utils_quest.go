package server

import (
	"fmt"
	"log/slog"
	"tap/internal/protocol"
)

func (p *Player) processFetchQuest() {
	for questID, questDataPlayer := range p.Quests {
		if questDataPlayer.Status != "active" {
			continue
		}

		worldQuest, exists := p.Server.world.Quests[questID]
		if !exists || worldQuest.Type != "fetch" {
			continue
		}

		newCompletedsteps := p.countHeldTargets(worldQuest.Targets)

		if newCompletedsteps > questDataPlayer.CompletedSteps {
			questDataPlayer.CompletedSteps = newCompletedsteps
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

// compte combien d'items à recuperer lors de la quête sont déjà dans l'inventaire du player ( s íl en a ramassé avant de commencer la quête...)
func (p *Player) countHeldTargets(targets []string) int {
	count := 0
	for _, target := range targets {
		for _, item := range p.Inventory {
			if p.matchNameOrID(target, item) {
				count++
				break
			}
		}
	}
	return count
}
