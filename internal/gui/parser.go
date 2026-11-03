package gui

import (
	"encoding/json"
	"fmt"
	"strings"
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
	HP     int    `json:"hp"`
	MaxHP  int    `json:"max_hp"`
	Status string `json:"status"`
}

type WhoResponse struct {
	All    []string `json:"all"`
	Server int      `json:"server"`
}

type FleeResponse struct {
	Direction string `json:"direction"`
	Room      string `json:"room"`
}

type QuestResponse struct {
	QuestID     string `json:"quest_id"`
	Status      string `json:"status"`
	Progress    string `json:"progress,omitempty"`
	Description string `json:"description,omitempty"`
	Reward      string `json:"reward,omitempty"`
}

type QuestsItemResponse struct {
	QuestID  string `json:"quest_id"`
	Status   string `json:"status"`
	Progress string `json:"progress,omitempty"`
}

type CreateGroupResponse struct {
	Group string
}

type EventChatMessage struct {
	Type    string
	Scope   string
	Channel string
	Sender  string
	Message string
}

type EventStatMessage struct {
	Type    string
	Command string
	Players int
}

type EventInvite struct {
	Type    string
	Command string
	Action  string
	Group   string
}

func parseInvite(raw string) (EventInvite, error) {
	raw = strings.TrimSpace(raw)

	parts := strings.SplitN(raw, " ", 4)
	if len(parts) < 4 {
		return EventInvite{}, fmt.Errorf("Invalid format : %q", raw)
	}

	return EventInvite{
		Type:    parts[0],
		Command: parts[1],
		Action:  parts[2],
		Group:   parts[3],
	}, nil
}

func parseEventChat(raw string) (EventChatMessage, error) {
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

func parseEventStats(raw string) (EventStatMessage, error) {
	raw = strings.TrimSpace(raw)

	parts := strings.SplitN(raw, " ", 3)
	if len(parts) < 3 {
		return EventStatMessage{}, fmt.Errorf("invalid event stats format: %q", raw)
	}

	if parts[0] != "EVT" || parts[1] != "STATS" {
		return EventStatMessage{}, fmt.Errorf("not an event stats command: %q", raw)
	}

	evt := EventStatMessage{
		Type:    parts[0],
		Command: parts[1],
	}

	if strings.HasPrefix(parts[2], "players=") {
		valStr := strings.TrimPrefix(parts[2], "players=")
		var count int
		_, err := fmt.Sscanf(valStr, "%d", &count)
		if err == nil {
			evt.Players = count
		} else {
			return EventStatMessage{}, fmt.Errorf("failed to parse players count from: %q", parts[2])
		}
	}

	return evt, nil
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

func parseAttack(strAttack string) (AttackResponse, error) {
	return parseJSON[AttackResponse](strAttack)
}
func parseFlee(strFlee string) (FleeResponse, error) {
	return parseJSON[FleeResponse](strFlee)
}

func parseDefence(strDefence string) (DefendResponse, error) {
	return parseJSON[DefendResponse](strDefence)
}

func parseStatus(strStatus string) (StatusResponse, error) {
	return parseJSON[StatusResponse](strStatus)
}

func parseQuest(strQuest string) (QuestResponse, error) {
	return parseJSON[QuestResponse](strQuest)
}

func parseQuests(strQuests string) ([]QuestsItemResponse, error) {
	var quests []QuestsItemResponse
	jsonStart := strings.Index(strQuests, "[")
	if jsonStart == -1 {
		return nil, fmt.Errorf("No array found: %q", strQuests)
	}

	err := json.Unmarshal([]byte(strQuests[jsonStart:]), &quests)
	if err != nil {
		return nil, fmt.Errorf("erreur unmarshal Quests: %w", err)
	}

	return quests, nil
}

func parseWho(strWho string) (WhoResponse, error) {
	return parseJSON[WhoResponse](strWho)
}

func parseCreateGroup(raw string) (CreateGroupResponse, error) {
	raw = strings.TrimSpace(raw)

	if strings.Contains(raw, "group=") {
		parts := strings.Split(raw, "group=")
		if len(parts) > 1 {
			groupName := strings.TrimSpace(parts[1])
			return CreateGroupResponse{Group: groupName}, nil
		}
	}

	return CreateGroupResponse{}, fmt.Errorf("impossible d'extraire le nom du groupe depuis: %q", raw)
}

func parseInventory(raw string) ([]string, error) {
	var items []string

	jsonStart := strings.Index(raw, "[")
	if jsonStart == -1 {
		return nil, fmt.Errorf("no inventory JSON array found in: %q", raw)
	}

	err := json.Unmarshal([]byte(raw[jsonStart:]), &items)
	return items, err
}
