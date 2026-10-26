package gui

import (
	"fmt"
	"strings"
	"tap/internal/gui/components"
)

func (g *Game) serverResponse() error {
	for {
		select {
		case msg := <-g.messagesChan:
			chatStruct, errChat := ParseEventChat(msg)
			fmt.Print(msg)
			if errChat == nil && chatStruct.Type == "EVT" {
				if chat := g.bottomCenter.GetChat(); chat != nil {
					chat.AddLine("[" + chatStruct.Scope + "] " + chatStruct.Sender + " > " + chatStruct.Message)
				}
				continue
			}
			inviteStruct, errstruct := ParseInvite(msg)
			if errstruct == nil && inviteStruct.Type == "EVT" && inviteStruct.Cmd == "GROUP" {
				if chat := g.bottomCenter.GetChat(); chat != nil {
					chat.AddLine("INVITATION> Invitation Received, group: '" + inviteStruct.Group + "'")
				}
				continue
			}

			if chatAction := g.bottomRight.GetChat(); chatAction != nil {
				jsonStart := strings.Index(msg, "{")
				response := msg
				if jsonStart != -1 {
					response = msg[jsonStart:]
				}
				stuff := g.bottomLeft.GetStuff()

				switch g.curr_cmd {
				case "LOOK":
					lookResp, err := ParseLookResponse(msg)
					if err == nil {
						room := lookResp.ToWorldRoom()
						g.centerMid.AddRoom(
							room.Name,
							g.location,
							room.Exits,
							room.Items,
							func(itemName string) {
								if g.status == "AWAITING_TAKE_ITEM" {
									cmd := "TAKE " + itemName
									g.curr_cmd = "TAKE"
									SendMessage(g.conn, cmd)
									g.status = ""
								}
							},
							room.NPCs,
							func(NPCsName string) {
								if g.status == "AWAITING_TALK_NPC" {
									cmd := "TALK " + NPCsName
									g.curr_cmd = "TALK"
									SendMessage(g.conn, cmd)
									g.status = ""
								} else if g.status == "AWAITING_QUEST_NPC" {
									cmd := "QUEST " + NPCsName
									g.curr_cmd = "QUEST"
									SendMessage(g.conn, cmd)
									g.status = ""
								} else if g.status == "AWAITING_ATTACK_NPC" {
									cmd := "ATTACK " + NPCsName
									g.curr_cmd = "ATTACK"
									SendMessage(g.conn, cmd)
									g.status = ""
								}
							},
						)
					} else {
						if chatAction != nil {
							chatAction.AddLine("Server > " + msg)
						}
					}

				case "MOVE":
					if strings.Contains(msg, "room=") {
						parts := strings.Split(msg, "room=")
						if len(parts) > 1 {
							room := strings.TrimSpace(parts[1])
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
							if chatAction != nil {
								chatAction.AddLine(g.curr_cmd + "> You moved into " + room)
							}
						}
					} else if msg == "ERR 301 NO_EXIT" {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + "> No exit in this direction")
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + "> Cannot move there")
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
											cmd := "DROP " + itemName
											g.curr_cmd = "DROP"
											stuff.RemoveItem(itemName)
											SendMessage(g.conn, cmd)
											g.status = ""
										} else if g.status == "AWAITING_EQUIP_ITEM" {
											cmd := "EQUIP " + itemName
											g.curr_cmd = "EQUIP"
											stuff.MoveToTop(itemName)
											SendMessage(g.conn, cmd)
											g.status = ""
										}
									},
								)
							}
							if chatAction != nil {
								chatAction.AddLine(g.curr_cmd + " > You took " + itemName)
							}
							g.curr_cmd = "LOOK"
							SendMessage(g.conn, "LOOK")
						}
					} else {
						chatAction.AddLine("Server > " + msg)
					}

				case "DROP":
					if strings.Contains(msg, "dropped=") {
						parts := strings.Split(msg, "dropped=")
						if len(parts) > 1 {
							itemName := strings.TrimSpace(parts[1])
							if chatAction != nil {
								chatAction.AddLine(g.curr_cmd + " > You dropped " + itemName)
							}
						}
					}
					g.status = ""
					g.curr_cmd = "LOOK"
					SendMessage(g.conn, "LOOK")

				case "ATTACK":
					if strings.Contains(msg, "Defeat:") || strings.Contains(msg, "defeated") {
						if chatAction != nil {
							chatAction.AddLine("COMBAT > YOU DIED! Respawned in room0.")
						}
						g.centerMid.RemoveRooms()
						g.location = [2]int{1, 1}
						g.status = ""
						g.curr_cmd = "LOOK"
						SendMessage(g.conn, "LOOK")
						break
					}
					if response != "ERR 406 NPC_NOT_HOSTILE" && strings.Contains(msg, "{") {
						attackJSON, err := ParseAttack(response)
						if err == nil {
							if attackJSON.Status == "victory" {
								if chatAction != nil {
									chatAction.AddLine(g.curr_cmd + " > Enemy defeated!")
								}
								g.hp = attackJSON.AttackerHP
								g.status = ""
								g.curr_cmd = "LOOK"
								SendMessage(g.conn, "LOOK")
							} else if attackJSON.Status == "defeated" {
								if chatAction != nil {
									chatAction.AddLine(g.curr_cmd + " > YOU DIED")
								}
								g.centerMid.RemoveRooms()
								g.location = [2]int{1, 1}
								g.status = ""
								g.curr_cmd = "LOOK"
								g.hp = attackJSON.AttackerHP
								SendMessage(g.conn, "LOOK")
							} else {
								attackMsg := fmt.Sprintf("Enemy took %d damage (%d HP left)", attackJSON.Damage, attackJSON.TargetHP)
								g.hp = attackJSON.AttackerHP
								if chatAction != nil {
									chatAction.AddLine(g.curr_cmd + " > " + attackMsg)
								}
							}
						} else {
							if chatAction != nil {
								chatAction.AddLine("Server > " + msg)
							}
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + " > Target is not hostile")
						}
					}
					g.status = ""

				case "DEFEND":
					if strings.Contains(msg, "Defeat:") || strings.Contains(msg, "defeated") {
						chatAction.AddLine("COMBAT > YOU DIED! Respawned in room0.")
						g.centerMid.RemoveRooms()
						g.location = [2]int{1, 1}
						g.status = ""
						g.curr_cmd = "LOOK"
						SendMessage(g.conn, "LOOK")
						break
					}
					if strings.Contains(msg, "{") {
						defenceJSON, err := ParseDefence(response)
						if err == nil {
							if defenceJSON.Status == "defeated" {
								chatAction.AddLine("DEFEND > YOU DIED")
								g.centerMid.RemoveRooms()
								g.location = [2]int{1, 1}
								g.curr_cmd = "LOOK"
								SendMessage(g.conn, "LOOK")
							} else if defenceJSON.Status == "perfect_defense" {
								if chatAction != nil {
									chatAction.AddLine("DEFEND > Protected! No damage taken")
								}
							} else if defenceJSON.Status == "weakened_defense" || defenceJSON.Status == "broken_defense" {
								defenceMsg := fmt.Sprintf("Defense broken! %d damage taken", defenceJSON.Damage)
								if chatAction != nil {
									chatAction.AddLine("DEFEND > " + defenceMsg)
								}
							}
							g.hp = defenceJSON.DefenderHP
						} else {
							if chatAction != nil {
								chatAction.AddLine("Server > " + msg)
							}
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine("Server > " + msg)
						}
					}
					g.status = ""

				case "TALK":
					talkResponse := response
					if len(response) >= 3 {
						talkResponse = response[3:]
					}
					if chatAction != nil {
						chatAction.AddLine(g.curr_cmd + " > " + talkResponse)
					}

				case "EQUIP":
					if strings.Contains(msg, "equipped=") {
						parts := strings.Split(msg, "equipped=")
						if len(parts) > 1 {
							itemName := strings.TrimSpace(parts[1])
							if chatAction != nil {
								chatAction.AddLine(g.curr_cmd + " > Equipped " + itemName)
							}
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + " > Equip failed")
						}
					}
					g.status = ""

				case "WHO":
					who, err := ParseWho(response)
					if err == nil {
						names := ""

						invitesPanel := g.centerLeft.GetSubPanel("INVITE PLAYER")
        				if invitesPanel != nil {
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
        				            SendMessage(g.conn, "GROUP INVITE "+playerName)
        				        }
        				    }
        				}

						for _, name := range who.Room {
							names += name + " "
						}
						whoMsg := fmt.Sprintf("Players: %s| Server: %d", names, who.Server)
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + " > " + whoMsg)
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine("Server> " + msg)
						}
					}

				case "STATUS":
					status, err := ParseStatus(response)
					if err == nil {
						statusMsg := fmt.Sprintf("%d/%d HP | State: %s", status.HP, status.MaxHP, status.Status)
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + " > " + statusMsg)
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine("Server > " + msg)
						}
					}

				case "QUEST":
					quest, err := ParseQuest(response)
					if err == nil {
						statusMsg := fmt.Sprintf("%s | reward: %s", quest.Description, quest.Reward)
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + " > " + statusMsg)
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine("Server > " + msg)
						}
					}

				case "QUESTS":
					quests, err := ParseQuests(msg)
					if err == nil {
						if len(quests) == 0 {
							if chatAction != nil {
								chatAction.AddLine(g.curr_cmd + " > No quest has been started.")
							}
						} else {
							for _, quest := range quests {
								statusMsg := fmt.Sprintf("[%s] status: %s, progess: %s", quest.QuestID, quest.Status, quest.Progress)

								if chatAction != nil {
									chatAction.AddLine(g.curr_cmd + " > " + statusMsg)
								}
							}
						}
					} else {
						if chatAction != nil {
							chatAction.AddLine("Server > " + msg)
						}
					}
					g.status = ""

				case "QUIT":
					return ErrQuit

				case "CREATE GROUP":
					if strings.HasPrefix(msg, "ERR") {
						if msg == "ERR 402 ALREADY_IN_GROUP" {
							if chatAction != nil {
								chatAction.AddLine("GROUP> Cannot create group: you are already in a group")
							}
						} else {
							if chatAction != nil {
								chatAction.AddLine("GROUP> " + msg)
							}
						}
					} else {
						group, err := ParseCreateGroup(msg)
						if err == nil {
							if chatAction != nil {
								chatAction.AddLine("GROUP> You created " + group.Group)
							}
						} else {
							if chatAction != nil {
								chatAction.AddLine("GROUP> " + msg)
							}
						}
					}
					g.status = ""

				case "LEAVE GROUP":
					if msg == "ERR 401 NOT_IN_GROUP" {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + "> You can not leave because you are not in a group")
						}
					} else {
						if chatAction != nil {
							msg := fmt.Sprintf("You leaved the group")
							chatAction.AddLine("GROUP> " + msg)
						}
					}

				case "INVITE":
					fmt.Print(msg)
					if msg == "ERR 403 UNKNOWN_PLAYER" {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + "> Player not found")
						}
					} else if msg == "ERR 401 NOT_IN_GROUP" {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + "> You can not invite, because you are not in a group")
						}
					} else if msg == "ERR 402 ALREADY_IN_GROUP" {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + "> Already in a group")
						}
					} else {
						if chatAction != nil {
							msg := fmt.Sprintf("You invitation is pending")
							chatAction.AddLine("GROUP> " + msg)
						}
					}

				default:
					if msg == "ERR 403 UNKNOWN_PLAYER" {
						if chatAction != nil {
							chatAction.AddLine(g.curr_cmd + "> Player not found")
						}
					} else if msg == "ERR 402 ALREADY_IN_GROUP" {
						if chatAction != nil {
							chatAction.AddLine("GROUP> Already in a group")
						}
					} else if msg == "ERR 401 NOT_IN_GROUP" {
						if chatAction != nil {
							chatAction.AddLine("GROUP> You are not in a group")
						}
					} else if msg == "ERR 406 NPC_NOT_HOSTILE" {
						if chatAction != nil {
							chatAction.AddLine("ATTACK> You are not in a group")
						}
					} else if msg == "ERR 301 NO_EXIT" {
						if chatAction != nil {
							chatAction.AddLine("MOVE> No exit in this direction")
						}
					} else if msg == "ERR 404 ITEM_NOT_FOUND" {
						if chatAction != nil {
							chatAction.AddLine("TAKE> Item not found")
						}
					} else if msg == "ERR 404 ITEM_NOT_IN_INVENTORY" {
						if chatAction != nil {
							chatAction.AddLine("ITEM> Item not found in inventory")
						}
					} else if msg == "ERR 404 NPC_NOT_FOUND" {
						if chatAction != nil {
							chatAction.AddLine("NPC> NPC not found")
						}
					} else if msg == "ERR 406 NO_QUEST_AVAILABLE" {
						if chatAction != nil {
							chatAction.AddLine("QUEST> The NPC has no quest to give")
						}
					}
				}
			}

		default:
			goto Done
		}
	}
Done:
	return nil
}
