package gui

import (
	"fmt"
	"strings"
	"tap/internal/gui/components"
)

func handleServerError(chat *components.ChatWidget, cmd, msg string) {
	if chat == nil {
		return
	}

	switch {
	case strings.HasPrefix(msg, "ERR 301"):
		chat.AddLine(cmd + " > No exit in this direction.")
	case strings.HasPrefix(msg, "ERR 401"):
		chat.AddLine(cmd + " > You are not in a group.")
	case strings.HasPrefix(msg, "ERR 402"):
		chat.AddLine(cmd + " > Already in a group.")
	case strings.HasPrefix(msg, "ERR 403"):
		chat.AddLine(cmd + " > Player not found.")
	case strings.HasPrefix(msg, "ERR 404 ITEM_NOT_FOUND"):
		chat.AddLine(cmd + " > Item not found on the ground.")
	case strings.HasPrefix(msg, "ERR 404 ITEM_NOT_IN_INVENTORY"):
		chat.AddLine(cmd + " > Item not found in your inventory.")
	case strings.HasPrefix(msg, "ERR 404 NPC_NOT_FOUND"):
		chat.AddLine(cmd + " > NPC not found.")
	case strings.HasPrefix(msg, "ERR 406 NPC_NOT_HOSTILE"):
		chat.AddLine(cmd + " > This NPC is not hostile.")
	case strings.HasPrefix(msg, "ERR 406 NO_QUEST_AVAILABLE"):
		chat.AddLine(cmd + " > This NPC has no available quests.")
	default:
		chat.AddLine("Server > " + msg)
	}
}

func (g *Game) serverResponse() error {
	for {
		select {
		case msg := <-g.messagesChan:
			fmt.Print(msg)
			if chatStruct, errChat := parseEventChat(msg); errChat == nil && chatStruct.Type == "EVT" {
				if chat := g.bottomCenter.GetChat(); chat != nil {
					chat.AddLine("[" + chatStruct.Scope + "] " + chatStruct.Sender + " > " + chatStruct.Message)
				}
				continue
			}

			if inviteStruct, errStruct := parseInvite(msg); errStruct == nil && inviteStruct.Type == "EVT" {
				if inviteStruct.Action == "INVITE" {
					if invitationPanel := g.centerLeft.GetSubPanel("INVITATION"); invitationPanel != nil {
						nextIndex := len(invitationPanel.AcceptInvites)
						newInvitation := components.NewPlayerInvite(
							inviteStruct.Group,
							nextIndex,
							invitationPanel.X,
							invitationPanel.Y,
							invitationPanel.W,
						)

						groupName := inviteStruct.Group
						newInvitation.Btn.OnClick = func() {
							g.curr_cmd = "GROUP"
							sendMessage(g.conn, "GROUP JOIN "+groupName)
						}

						invitationPanel.AcceptInvites = append(invitationPanel.AcceptInvites, newInvitation)
					}

					if chat := g.bottomCenter.GetChat(); chat != nil {
						chat.AddLine("[INVITATION] > Received invite from group: '" + inviteStruct.Group + "'")
					}
					continue
				} else if inviteStruct.Action == "JOIN" {
					if chat := g.bottomCenter.GetChat(); chat != nil {
						chat.AddLine("[INVITATION] > '" + inviteStruct.Group + "' joined the group")
					}
					continue
				}
			}

			if chatAction := g.bottomRight.GetChat(); chatAction != nil {
				if strings.HasPrefix(msg, "ERR") {
					handleServerError(chatAction, g.curr_cmd, msg)
					g.status = ""
					continue
				}

				jsonStart := strings.Index(msg, "{")
				response := msg
				if jsonStart != -1 {
					response = msg[jsonStart:]
				}
				stuff := g.bottomLeft.GetStuff()

				fields := strings.Fields(g.curr_cmd)
				if len(fields) == 0 {
					continue
				}
				verb := strings.ToUpper(fields[0])

				switch verb {
				case "LOOK":
					if lookResp, err := parseLookResponse(msg); err == nil {
						room := lookResp.toWorldRoom()
						g.centerMid.AddRoom(
							room.Name,
							g.location,
							room.Exits,
							room.Items,
							func(itemName string) {
								if g.status == "AWAITING_TAKE_ITEM" {
									g.curr_cmd = "TAKE"
									sendMessage(g.conn, "TAKE "+itemName)
									g.status = ""
								}
							},
							room.NPCs,
							func(NPCsName string) {
								switch g.status {
								case "AWAITING_TALK_NPC":
									g.curr_cmd = "TALK"
									sendMessage(g.conn, "TALK "+NPCsName)
								case "AWAITING_QUEST_NPC":
									g.curr_cmd = "QUEST"
									sendMessage(g.conn, "QUEST "+NPCsName)
								case "AWAITING_ATTACK_NPC":
									g.curr_cmd = "ATTACK"
									sendMessage(g.conn, "ATTACK "+NPCsName)
								}
								g.status = ""
							},
						)
						chatAction.AddLine("LOOK > " + lookResp.Room.Name + ": " + lookResp.Room.Description)
					}

				case "MOVE":
					if strings.Contains(msg, "room=") {
						parts := strings.Split(msg, "room=")
						if len(parts) > 1 {
							room := strings.TrimSpace(parts[1])

							if len(fields) > 1 {
								direction := strings.ToUpper(fields[1])
								switch direction {
								case "NORTH":
									g.location[1] -= 1
									g.direction = 0
								case "EAST":
									g.location[0] += 1
									g.direction = 1
								case "SOUTH":
									g.location[1] += 1
									g.direction = 2
								case "WEST":
									g.location[0] -= 1
									g.direction = 3
								}
							} else {
								switch g.direction {
								case 0:
									g.location[1] -= 1
								case 1:
									g.location[0] += 1
								case 2:
									g.location[1] += 1
								case 3:
									g.location[0] -= 1
								}
							}

							chatAction.AddLine("MOVE > Moved into " + room)
						}
					}

				case "TAKE":
					if strings.Contains(msg, "taken=") {
						parts := strings.Split(msg, "taken=")
						if len(parts) > 1 {
							itemName := strings.TrimSpace(parts[1])
							if stuff != nil {
								stuff.AddItem(
									itemName,
									func(itemName string) {
										if g.status == "AWAITING_DROP_ITEM" {
											g.curr_cmd = "DROP"
											stuff.RemoveItem(itemName)
											sendMessage(g.conn, "DROP "+itemName)
											g.status = ""
										} else if g.status == "AWAITING_EQUIP_ITEM" {
											g.curr_cmd = "EQUIP"
											stuff.MoveToTop(itemName)
											sendMessage(g.conn, "EQUIP "+itemName)
											g.status = ""
										}
									},
								)
							}
							chatAction.AddLine("TAKE > Item taken: " + itemName)
							g.curr_cmd = "LOOK"
							sendMessage(g.conn, "LOOK")
						}
					}

				case "DROP":
					if strings.Contains(msg, "dropped=") {
						parts := strings.Split(msg, "dropped=")
						if len(parts) > 1 {
							chatAction.AddLine("DROP > Item dropped: " + strings.TrimSpace(parts[1]))
						}
					}
					g.status = ""
					g.curr_cmd = "LOOK"
					sendMessage(g.conn, "LOOK")

				case "ATTACK":
					if strings.Contains(msg, "Defeat:") || strings.Contains(msg, "defeated") {
						chatAction.AddLine("COMBAT > YOU DIED! Respawned in hotel_entrance.")
						g.centerMid.RemoveRooms()
						g.location = [2]int{1, 1}
						g.status = ""
						g.curr_cmd = "LOOK"
						sendMessage(g.conn, "LOOK")
						break
					}
					if attackJSON, err := parseAttack(response); err == nil {
						if attackJSON.Status == "victory" {
							chatAction.AddLine("ATTACK > Enemy defeated!")
							g.hp = attackJSON.AttackerHP
							g.status = ""
							g.curr_cmd = "LOOK"
							sendMessage(g.conn, "LOOK")
						} else if attackJSON.Status == "defeated" {
							chatAction.AddLine("ATTACK > YOU DIED")
							g.centerMid.RemoveRooms()
							g.location = [2]int{1, 1}
							g.status = ""
							g.curr_cmd = "LOOK"
							g.hp = attackJSON.AttackerHP
							sendMessage(g.conn, "LOOK")
						} else {
							chatAction.AddLine(fmt.Sprintf("ATTACK > Enemy -%d HP (%d HP remaining)", attackJSON.Damage, attackJSON.TargetHP))
							g.hp = attackJSON.AttackerHP
						}
					}
					g.status = ""

				case "DEFEND":
					if defenceJSON, err := parseDefence(response); err == nil {
						if defenceJSON.Status == "perfect_defense" {
							chatAction.AddLine("DEFEND > Perfect defense! No damage taken.")
						} else if defenceJSON.Status == "weakened_defense" || defenceJSON.Status == "broken_defense" {
							chatAction.AddLine(fmt.Sprintf("DEFEND > Defense broken! %d damage taken.", defenceJSON.Damage))
						}
						g.hp = defenceJSON.DefenderHP
					}
					g.status = ""

				case "TALK":
					talkResponse := response
					if len(response) >= 3 {
						talkResponse = response[3:]
					}
					chatAction.AddLine("TALK > " + talkResponse)

				case "EQUIP":
					if strings.Contains(msg, "equipped=") {
						parts := strings.Split(msg, "equipped=")
						if len(parts) > 1 {
							chatAction.AddLine("EQUIP > Equipped: " + strings.TrimSpace(parts[1]))
						}
					}
					g.status = ""

				case "WHO":
					if who, err := parseWho(response); err == nil {
						if invitesPanel := g.centerLeft.GetSubPanel("INVITE PLAYER"); invitesPanel != nil {
							invitesPanel.PlayerInvites = components.CreateList(
								g.name,
								who.Room,
								invitesPanel.X,
								invitesPanel.Y,
								invitesPanel.W,
							)
							for _, invite := range invitesPanel.PlayerInvites {
								playerName := invite.Name
								invite.Btn.OnClick = func() {
									g.curr_cmd = "INVITE"
									sendMessage(g.conn, "GROUP INVITE "+playerName)
								}
							}
						}
						chatAction.AddLine(fmt.Sprintf("WHO > Players in room: %v | Server: %d", who.Room, who.Server))
					}

				case "STATUS":
					if status, err := parseStatus(response); err == nil {
						chatAction.AddLine(fmt.Sprintf("STATUS > %d/%d HP | State: %s", status.HP, status.MaxHP, status.Status))
					}

				case "QUEST":
					if quest, err := parseQuest(response); err == nil {
						chatAction.AddLine(fmt.Sprintf("QUEST > %s | Reward: %s", quest.Description, quest.Reward))
						g.curr_cmd = "INVENTORY"
						sendMessage(g.conn, "INVENTORY")
					}

				case "QUESTS":
					if quests, err := parseQuests(msg); err == nil {
						if len(quests) == 0 {
							chatAction.AddLine("QUESTS > No active quests.")
						} else {
							for _, quest := range quests {
								chatAction.AddLine(fmt.Sprintf("QUESTS > [%s] Status: %s | Progress: %s", quest.QuestID, quest.Status, quest.Progress))
							}
						}
					}
					g.status = ""

				case "INVENTORY":
					if inventory, err := parseInventory(msg); err == nil {
						if stuff := g.bottomLeft.GetStuff(); stuff != nil {
							stuff.Clear()

							for _, item := range inventory {
								itemName := item
								stuff.AddItem(
									itemName,
									func(selectedItem string) {
										if g.status == "AWAITING_DROP_ITEM" {
											g.curr_cmd = "DROP"
											stuff.RemoveItem(selectedItem)
											sendMessage(g.conn, "DROP "+selectedItem)
											g.status = ""
										} else if g.status == "AWAITING_EQUIP_ITEM" {
											g.curr_cmd = "EQUIP"
											stuff.MoveToTop(selectedItem)
											sendMessage(g.conn, "EQUIP "+selectedItem)
											g.status = ""
										}
									},
								)
							}
						}
					}
					g.status = ""

				case "CREATE", "GROUP":
					if group, err := parseCreateGroup(msg); err == nil {
						if chat := g.bottomCenter.GetChat(); chat != nil {
							chat.AddLine("[INVITATION] > You are now in: " + group.Group)
						}
					} else {
						if chat := g.bottomCenter.GetChat(); chat != nil {
							chat.AddLine("[INVITATION] > " + msg)
						}
					}
					g.status = ""

				case "LEAVE":
					chatAction.AddLine("GROUP > You left the group.")
					g.status = ""

				case "QUIT":
					return ErrQuit

				default:
					chatAction.AddLine("Server > " + msg)
				}
			}

		default:
			goto Done
		}
	}
Done:
	return nil
}
