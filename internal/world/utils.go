package world

import (
	"fmt"
	"strings"
)

func (world *GameData) CheckEmptyWorld() error {
	if len(world.Rooms) == 0 {
		return fmt.Errorf("world data must contain at least one room")
	}
	return nil
}

func (world *GameData) ValidateNoDuplicateItem() error {
	seenItems := make(map[string]string)
	for roomID, room := range world.Rooms {
		for _, item := range room.Items {
			normalizedItem := strings.ToLower(strings.TrimSpace(item))
			if normalizedItem == "" {
				continue
			}

			if existingInRoom, exists := seenItems[normalizedItem]; exists {
				if existingInRoom == roomID {
					return fmt.Errorf("duplicated item '%s' found multiple times in room '%s'", item, roomID)
				}
				return fmt.Errorf("duplicated item '%s' found in both '%s' and '%s'", item, existingInRoom, roomID)
			}
			seenItems[normalizedItem] = roomID
		}
	}
	return nil
}

func (world *GameData) CheckItemExists() error {
	for roomID, room := range world.Rooms {
		for _, targetItem := range room.Items {
			targetItem = strings.ToLower(strings.TrimSpace(targetItem))
			if targetItem == "" {
				continue
			}

			targetExists := false
			for itemID, itemDefined := range world.Items {
				if strings.EqualFold(itemID, targetItem) || strings.EqualFold(itemDefined.Name, targetItem) {
					targetExists = true
					break
				}
			}

			if !targetExists {
				return fmt.Errorf("item '%s' in room '%s' does not exist in items table", targetItem, roomID)
			}
		}
	}
	return nil
}

func (world *GameData) ValidateExit() error {
	if _, exists := world.Rooms["hotel_entrance"]; !exists {
		return fmt.Errorf("mandatory starting room 'hotel_entrance' is missing from world")
	}
	for roomID, room := range world.Rooms {
		for dir, exit := range room.Exits {
			cleanDir := strings.ToLower(strings.TrimSpace(dir))
			if cleanDir != "north" && cleanDir != "east" && cleanDir != "south" && cleanDir != "west" {
				return fmt.Errorf("exit '%s' in room '%s' is not a valid direction", dir, roomID)
			}
			if _, exists := world.Rooms[exit.Target]; !exists {
				return fmt.Errorf("exit '%s' in room '%s' targets non-existent room '%s'", dir, roomID, exit.Target)
			}
		}
	}
	return nil
}

func (world *GameData) CheckNPCExists() error {
	for roomID, room := range world.Rooms {
		for _, targetNPC := range room.NPCs {
			targetNPC = strings.ToLower(strings.TrimSpace(targetNPC))
			if targetNPC == "" {
				continue
			}

			targetExists := false
			for npcID, npcDefined := range world.NPCs {
				if strings.EqualFold(npcID, targetNPC) || strings.EqualFold(npcDefined.Name, targetNPC) {
					targetExists = true
					break
				}
			}

			if !targetExists {
				return fmt.Errorf("NPC '%s' in room '%s' does not exist in NPCs table", targetNPC, roomID)
			}
		}
	}

	for _, npc := range world.NPCs {
		if npc.Item != nil && *npc.Item != "" {
			foundItem := false
			for itemID, itemDefined := range world.Items {
				if strings.EqualFold(*npc.Item, itemID) || strings.EqualFold(*npc.Item, itemDefined.Name) {
					foundItem = true
					break
				}
			}
			if !foundItem {
				return fmt.Errorf("NPC '%s' holds non-existent item '%s'", npc.Name, *npc.Item)
			}
		}
	}

	return nil
}

func (world *GameData) ValidateQuests() error {
	for questID, quest := range world.Quests {
		giverExists := false
		for npcID, npc := range world.NPCs {
			if strings.EqualFold(npcID, quest.Giver) || strings.EqualFold(npc.Name, quest.Giver) {
				giverExists = true
				break
			}
		}
		if !giverExists {
			return fmt.Errorf("quest '%s' has non-existent giver NPC '%s'", questID, quest.Giver)
		}

		if quest.Type != "fetch" && quest.Type != "defeat" {
			return fmt.Errorf("quest '%s' has invalid type '%s' (expected 'fetch' or 'defeat')", questID, quest.Type)
		}

		for _, target := range quest.Targets {
			target = strings.TrimSpace(target)
			if target == "" {
				continue
			}

			if quest.Type == "fetch" {
				itemFound := false
				for itemID, item := range world.Items {
					if strings.EqualFold(itemID, target) || strings.EqualFold(item.Name, target) {
						itemFound = true
						break
					}
				}
				if !itemFound {
					return fmt.Errorf("fetch quest '%s' references non-existent target item '%s'", questID, target)
				}
			} else if quest.Type == "defeat" {
				npcFound := false
				for npcID, npc := range world.NPCs {
					if strings.EqualFold(npcID, target) || strings.EqualFold(npc.Name, target) {
						npcFound = true
						break
					}
				}
				if !npcFound {
					return fmt.Errorf("defeat quest '%s' references non-existent target NPC '%s'", questID, target)
				}
			}
		}

		reward := strings.TrimSpace(quest.Reward)
		if reward != "" {
			rewardFound := false
			for itemID, item := range world.Items {
				if strings.EqualFold(itemID, reward) || strings.EqualFold(item.Name, reward) {
					rewardFound = true
					break
				}
			}
			if !rewardFound {
				return fmt.Errorf("quest '%s' has non-existent reward item '%s'", questID, reward)
			}
		}

		prereq := strings.TrimSpace(quest.Prerequisite)
		if prereq != "" {
			if _, exists := world.Quests[prereq]; !exists {
				return fmt.Errorf("quest '%s' references non-existent prerequisite quest '%s'", questID, prereq)
			}
		}
	}
	return nil
}
