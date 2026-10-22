package server

import (
	"fmt"
	"log/slog"
	"strings"
	"tap/internal/protocol"
)

func (p *Player) handleQuest(args string) {
	targetNpc := strings.TrimSpace(args)

	p.Server.mu.Lock()
	questGiverID := p.findQuestGiver(targetNpc)
	if questGiverID == "" {
		return
	}

	if p.Server.world.Quests == nil {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNoQuestAvailable, protocol.MsgNoQuestAvailable)
		return
	}

	var activeQuestID string
	for questID, questData := range p.Quests {
		if questData.Status == "active" {
			activeQuestID = questID
			break
		}
	}
	if activeQuestID != "" {
		p.resolveActiveQuest(activeQuestID)
		return
	}

	p.processNextQuest(questGiverID)
}

// type QuestResponse struct {

// }

func (p *Player) handleQuests(args string) {
	if args != "" {
		p.SendErr(protocol.CodeBadRequest, protocol.MsgArgsError)
		return
	}
	// response := QuestResponse {

	// }
	p.SendOK("Received QUESTS. (TODO)")
}

func (p *Player) findQuestGiver(targetNpc string) string {
	var questGiverID string

	room := p.Server.world.Rooms[p.CurrentRoom]
	if targetNpc == "" {
		for _, npcID := range room.NPCs {
			npc := p.Server.world.NPCs[npcID]
			if npc.Role == "quest-giver" {
				questGiverID = npcID
				break
			}
		}
	} else {
		for _, npcID := range room.NPCs {
			if p.matchNameOrID(npcID, targetNpc) {
				questGiverID = npcID
				break
			}
		}
	}

	if questGiverID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNPCNotFound, protocol.MsgNPCNotFound)
		return ""
	}

	questGiverData := p.Server.world.NPCs[questGiverID]
	if questGiverData.Role != "quest-giver" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNoQuestAvailable, protocol.MsgNoQuestAvailable)
		return ""
	}
	return questGiverID
}

func (p *Player) resolveActiveQuest(activeQuestID string) {
	playerQuest := p.Quests[activeQuestID]
	worldQuest := p.Server.world.Quests[activeQuestID]
	if playerQuest.CompletedSteps >= playerQuest.Steps {
		playerQuest.Status = "completed"
		p.Quests[activeQuestID] = playerQuest
		playerReward := p.Server.world.Quests[activeQuestID].Reward
		if playerReward != "" {
			p.Inventory = append(p.Inventory, playerReward)
		}

		p.broadcastRoom(p.CurrentRoom, "ROOM CHAT", fmt.Sprintf("Player %s completed the quest: %s!", p.Username, activeQuestID))

		p.resolveFinalVictory()

		p.Server.mu.Unlock()
		p.SendOK(fmt.Sprintf(`{"quest_id":"%s", "description": "%s", "reward":"%s", "status":"completed"}`, activeQuestID, worldQuest.Description, worldQuest.Reward))
		return
	} else {
		p.Server.mu.Unlock()
		p.SendOK(fmt.Sprintf(`{"quest_id":"%s","description":"%s","reward":"%s","status":"active","progress":"%d/%d"}`,
			activeQuestID, worldQuest.Description, worldQuest.Reward, playerQuest.CompletedSteps, playerQuest.Steps))
		return
	}
}

func (p *Player) processNextQuest(questGiverID string) {
	var nextQuestID string

	for questID, worldQ := range p.Server.world.Quests {
		if worldQ.Giver != questGiverID {
			continue
		}

		quest, exists := p.Quests[questID]
		if exists {
			if quest.Status == "completed" {
				continue
			}
		}

		if worldQ.Prerequisite == "" || p.Quests[worldQ.Prerequisite].Status == "completed" {
			nextQuestID = questID
			break
		}
	}
	if nextQuestID == "" {
		p.Server.mu.Unlock()
		p.SendErr(protocol.CodeNoQuestAvailable, protocol.MsgNoQuestAvailable)
		return
	}

	nextQuestData := p.Server.world.Quests[nextQuestID]
	p.Quests[nextQuestID] = Quest{
		ID:             nextQuestID,
		Status:         "active",
		Steps:          nextQuestData.Count,
		CompletedSteps: 0,
	}
	p.Server.mu.Unlock()
	p.SendOK(fmt.Sprintf(`{"quest_id":"%s", "description": "%s", "reward":"%s", "status":"active"}`, nextQuestID, nextQuestData.Description, nextQuestData.Reward))
}

func (p *Player) resolveFinalVictory() {
	for worldQuestID := range p.Server.world.Quests {
		if p.Quests[worldQuestID].Status != "completed" {
			return
		}
	}
	slog.Info("Player achieved final game victory", "player", p.Username)
	p.broadcastGlobal("GLOBAL CHAT", fmt.Sprintf("Player %s has completed all the quests. Congratulations !", p.Username))
}
