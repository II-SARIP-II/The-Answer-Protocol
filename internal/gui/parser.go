package gui

import (
	"encoding/json"
	"fmt"
	"strings"

	"tap/internal/server"
)
type AttackResponse struct {
	AttackerHP int    `json:"attacker_hp"`
	TargetHP   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}

type DefendResponse struct {
	DefenderHP int    `json:"defender_hp"`
	TargetHP   int    `json:"target_hp"`
	Damage     int    `json:"damage"`
	Status     string `json:"status"`
}

type StatusResponse struct {
	HP int   `json:"hp"`
	MaxHP int `json:"max_hp"`
	Status string `json:"status"`
}

type WhoResponse struct {
	Room   []string `json:"room"`
	Server int      `json:"Server"`
}

type EventChatMessage struct {
	Type    string // "EVT"
	Scope   string // "GLOBAL" ou "GROUP"
	Channel string // "CHAT"
	Sender  string // "Pâris"
	Message string // "ECho ECHO"
}

func ParseEventChat(raw string) (EventChatMessage, error) {
	raw = strings.TrimSpace(raw)

	parts := strings.SplitN(raw, " ", 5)
	if len(parts) < 5 {
		return EventChatMessage{}, fmt.Errorf("Invalid format : %q", raw)
	}

	return EventChatMessage{
		Type:    parts[0],
		Scope:   parts[1],
		Channel: parts[2],
		Sender:  parts[3],
		Message: parts[4],
	}, nil
}

func parseJSON[T any](raw string) (T, error) {
	var target T
	jsonStart := strings.Index(raw, "{")
	if jsonStart == -1 {
		return target, fmt.Errorf("aucun JSON trouvé dans la réponse: %q", raw)
	}

	err := json.Unmarshal([]byte(raw[jsonStart:]), &target)
	if err != nil {
		return target, fmt.Errorf("erreur unmarshal JSON: %w", err)
	}

	return target, nil
}

func ParseAttack(strAttack string) (AttackResponse, error) {
	return parseJSON[AttackResponse](strAttack)
}

func ParseDefence(strDefence string) (DefendResponse, error) {
	return parseJSON[DefendResponse](strDefence)
}

func ParseStatus(strStatus string) (StatusResponse, error) {
	return parseJSON[StatusResponse](strStatus)
}

func ParseWho(strWho string) (WhoResponse, error) {
	return parseJSON[WhoResponse](strWho)
}

func ParsePlayer(strPlayer string) (server.Player, error) {
	return parseJSON[server.Player](strPlayer)
}
