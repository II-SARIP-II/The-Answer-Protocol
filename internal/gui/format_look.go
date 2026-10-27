package gui

import (
	"encoding/json"
	"fmt"
	"strings"

	"tap/internal/world"
)

type ServerRoom struct {
	ID          string            `json:"id"`
	X           int               `json:"X"`
	Y           int               `json:"Y"`
	W           int               `json:"W"`
	H           int               `json:"H"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
}

type ServerLookResponse struct {
	Room    ServerRoom `json:"room"`
	Players []string   `json:"players"`
	Items   []string   `json:"items"`
	NPCs    []string   `json:"npcs"`
}

func (resp ServerLookResponse) toWorldRoom() world.Room {
	exits := make(map[string]world.Exit)
	for dir, target := range resp.Room.Exits {
		exits[dir] = world.Exit{
			Target:  target,
			DoorPos: 1,
			Wall:    0,
		}
	}

	return world.Room{
		Name:  resp.Room.Name,
		X:     resp.Room.X,
		Y:     resp.Room.Y,
		W:     resp.Room.H,
		H:     resp.Room.H,
		Exits: exits,
		Items: resp.Items,
		NPCs:  resp.NPCs,
	}
}

func parseLookResponse(msg string) (ServerLookResponse, error) {
	startIdx := strings.Index(msg, "{")
	endIdx := strings.LastIndex(msg, "}")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return ServerLookResponse{}, fmt.Errorf("aucun objet JSON valide trouvé")
	}

	cleanJSON := msg[startIdx : endIdx+1]

	var resp ServerLookResponse
	err := json.Unmarshal([]byte(cleanJSON), &resp)
	if err != nil {
		return ServerLookResponse{}, fmt.Errorf("erreur unmarshal: %w", err)
	}

	return resp, nil
}
