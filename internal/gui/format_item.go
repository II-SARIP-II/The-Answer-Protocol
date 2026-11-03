package gui

import (
	"encoding/json"
	"os"
	"strings"
)

const configFilePath = "./data/world.json"

type ItemConfig struct {
	Name   string `json:"name"`
	Damage int    `json:"damage"`
}

type GameData struct {
	Items map[string]ItemConfig `json:"items"`
}

func loadGameData() (GameData, error) {
	var data GameData
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(bytes, &data)
	return data, err
}

func getName(itemID string) string {
	itemID = strings.TrimSpace(itemID)
	gameData, err := loadGameData()
	if err != nil {
		return itemID
	}

	if item, exists := gameData.Items[itemID]; exists && item.Name != "" {
		return item.Name
	}
	return itemID
}
