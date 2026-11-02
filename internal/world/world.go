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
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	X           int             `json:"x"`
	Y           int             `json:"y"`
	W           int             `json:"w"`
	H           int             `json:"h"`
	Exits       map[string]Exit `json:"exits"`
	Items       []string        `json:"items"`
	NPCs        []string        `json:"npcs"`
}

type NPC struct {
	Role        string   `json:"role"`
	Name        string   `json:"name"`
	Life        int      `json:"life"`
	MaxLife     int      `json:"maxlife"`
	Item        *string  `json:"item"`
	Dialogue    []string `json:"dialogue"`
	InitialRole string   `json:"-"`
}

type Item struct {
	Name   string `json:"name"`
	Damage int    `json:"damage"`
}

type Quest struct {
	Giver        string   `json:"giver"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	Targets      []string `json:"targets"`
	Count        int      `json:"count"`
	Reward       string   `json:"reward"`
	Prerequisite string   `json:"prerequisite"`
}

type GameData struct {
	Rooms  map[string]Room  `json:"rooms"`
	NPCs   map[string]NPC   `json:"npcs"`
	Items  map[string]Item  `json:"items"`
	Quests map[string]Quest `json:"quests"`
}

func ReadJson(filepath string) (*GameData, error) {
	jsonFile, err := os.ReadFile(filepath)

	if err != nil {
		return nil, fmt.Errorf("Cannot read world file: %w", err)
	}

	var data GameData
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}

	if err := data.CheckEmptyWorld(); err != nil {
		return nil, fmt.Errorf("world validation error: %w", err)
	}

	if err := data.ValidateNoDuplicateItem(); err != nil {
		return nil, fmt.Errorf("world validation error: %w", err)
	}

	if err := data.CheckItemExists(); err != nil {
		return nil, fmt.Errorf("world validation error: %w", err)
	}

	if err := data.ValidateExit(); err != nil {
		return nil, fmt.Errorf("world validation error: %w", err)
	}

	if err := data.CheckNPCExists(); err != nil {
		return nil, fmt.Errorf("world validation error: %w", err)
	}

	if err := data.ValidateQuests(); err != nil {
		return nil, fmt.Errorf("world validation error: %w", err)
	}

	return &data, nil
}
