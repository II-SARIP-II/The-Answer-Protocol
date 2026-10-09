package world

import (
	"encoding/json"
	"fmt"
	"os"
)

type Exit struct {
	Target  string `json:"target"`
	DoorPos int    `json:"door_pos"`
	Wall    int    `json:"wall"`
}

type Room struct {
	Name  string          `json:"name"`
	X     int             `json:"x"`
	Y     int             `json:"y"`
	W     int             `json:"w"`
	H     int             `json:"h"`
	Exits map[string]Exit `json:"exits"`
	Items []string        `json:"items"`
	NPCs  []string        `json:"npcs"`
}

type NPC struct {
	Role string  `json:"role"`
	Name string  `json:"name"`
	Life int     `json:"life"`
	Item *string `json:"item"`
}

type Item struct {
	ID     int `json:"id"`
	Damage int `json:"damage"`
}

type GameData struct {
	Rooms  map[string]Room `json:"rooms"`
	NPCs   map[string]NPC  `json:"npcs"`
	Items  map[string]Item `json:"items"`
	Quests []string        `json:"quests"`
}

// func ReadJson() GameData {
// 	jsonFile, err := os.ReadFile("data/world.json")

// 	if err != nil {
// 		log.Fatalf("Erreur lors de la lecture du fichier : %v", err)
// 	}

// 	var data GameData
// 	if err := json.Unmarshal(jsonFile, &data); err != nil {
// 		log.Fatalf("Json parsing error: %v", err)
// 	}
// 	return data
// }

func ReadJson(filepath string) (*GameData, error) {
	jsonFile, err := os.ReadFile(filepath)

	if err != nil {
		return nil, fmt.Errorf("Cannot read world file: %w", err)
	}

	var data GameData
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}
	return &data, nil
}

