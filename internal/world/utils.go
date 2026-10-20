package world

import (
	"fmt"
	"strings"
)

func (g *GameData) ValidateItem() error {
	seenItems := make(map[string]string) // item normalisé -> ID de la pièce

	for roomID, room := range g.Rooms {
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
