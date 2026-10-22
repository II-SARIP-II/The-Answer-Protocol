package gui

import (
	"fmt"
	"strings"
)


func (g *Game)serverResponse() error {
	for {
		select {
	case msg := <-g.messagesChan:
		if chat := g.bottomRight.GetChat(); chat != nil {
			print("\n\n", msg)
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
					if chat != nil {
						chat.AddLine("Server > " + msg)
					}
				}

			case "MOVE":
				if strings.Contains(msg, "room=") {
					parts := strings.Split(msg, "room=")
					if len(parts) > 1 {
						room := strings.TrimSpace(parts[1])
						if chat != nil {
							chat.AddLine(g.curr_cmd + "> " + "You moved in " + room)
						}
					}
				} else {
					if msg == "ERR 301 NO_EXIT" {
						if chat != nil {
							chat.AddLine(g.curr_cmd + "> " + "You cannot move in this direction, No exit")
						}
					} else {
						if chat != nil {
							chat.AddLine(g.curr_cmd + "> " + "You cannot move there")
						}
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
						if chat != nil {
							chat.AddLine(g.curr_cmd + "> " + "You took a " + itemName)
						}
						g.curr_cmd = "LOOK"
						SendMessage(g.conn, "LOOK")
					}
				} else {
					if chat != nil {
						chat.AddLine("Server > " + msg)
					}
				}

			case "DROP":
				if strings.Contains(msg, "dropped=") {
					parts := strings.Split(msg, "dropped=")
					if len(parts) > 1 {
						itemName := strings.TrimSpace(parts[1])
						if chat != nil {
							chat.AddLine(g.curr_cmd + "> " + "You dropped a " + itemName)
						}
					}
				}
				g.status = ""
				g.curr_cmd = "LOOK"
				SendMessage(g.conn, "LOOK")

			case "ATTACK":
				if strings.Contains(msg, "Defeat:") || strings.Contains(msg, "defeated") {
    			    if chat != nil {
    			        chat.AddLine("COMBAT > YOU DIED! Respawn in room0.")
    			    }
    			    g.status = ""
    			    g.curr_cmd = "LOOK"
    			    SendMessage(g.conn, "LOOK")
    			    break
    			}
				if response != "ERR 406 NPC_NOT_HOSTILE" && strings.Contains(msg, "{") {
					attack_json, err := ParseAttack(response)
					if err == nil {
						if attack_json.Status == "victory" {
							if chat != nil {
								chat.AddLine(g.curr_cmd + " > " + "The enemy died")
							}
							g.hp = attack_json.AttackerHP
							g.status = ""
							g.curr_cmd = "LOOK"
							SendMessage(g.conn, "LOOK")
						} else if attack_json.Status == "defeated" {
							if chat != nil {
								chat.AddLine(g.curr_cmd + " > " + "YOU DIED")
							}
							g.status = ""
							g.curr_cmd = "LOOK"
							g.hp = attack_json.AttackerHP
							SendMessage(g.conn, "LOOK")
						} else {
							attack_msg := fmt.Sprintf("The enemy took %d damage, enemy: %dhp", attack_json.Damage, attack_json.TargetHP)
							g.hp = attack_json.AttackerHP
							if chat != nil {
								chat.AddLine(g.curr_cmd + " > " + attack_msg)
							}
						}
					} else {
						if chat != nil {
							chat.AddLine("Server > " + msg)
						}
					}
				} else {
					if chat != nil {
						chat.AddLine(g.curr_cmd + "> " + "Cannot attack, he is not hostile")
					}
				}
				g.status = ""

			case "DEFEND":
				if strings.Contains(msg, "Defeat:") || strings.Contains(msg, "defeated") {
    			    if chat != nil {
    			        chat.AddLine("COMBAT > YOU DIED! Respawn in room0.")
    			    }
    			    g.status = ""
    			    g.curr_cmd = "LOOK"
    			    SendMessage(g.conn, "LOOK")
    			    break
    			}

    			if strings.Contains(msg, "{") {
    			    defenceJSON, err := ParseDefence(response)
					if err == nil {
    			    	if defenceJSON.Status == "defeated" {
    			    	    if chat != nil {
    			    	        chat.AddLine("DEFEND > YOU DIED")
    			    	    }
    			    	    g.curr_cmd = "LOOK"
    			    	    SendMessage(g.conn, "LOOK")
    			    	} else if defenceJSON.Status == "perfect_defense" {
    			    	    if chat != nil {
    			    	        chat.AddLine("DEFEND > You protect yourself, no damage taken")
    			    	    }
    			    	} else if defenceJSON.Status == "weakened_defense" || defenceJSON.Status == "broken_defense" {
    			    	    defenceMsg := fmt.Sprintf("The enemy broke your defence! %d damage taken", defenceJSON.Damage)
    			    	    if chat != nil {
    			    	        chat.AddLine("DEFEND > " + defenceMsg)
    			    	    }
    			    	}
    			    	g.hp = defenceJSON.DefenderHP
					} else {
						if chat != nil {
							chat.AddLine("Server > " + msg)
						}
					}

    			} else {
    			    if chat != nil {
    			        chat.AddLine("Server > " + msg)
    			    }
    			}
    			g.status = ""

			case "TALK":
				talk_response := response[3:]
				if chat != nil {
					chat.AddLine(g.curr_cmd + "> " + talk_response)
				}

			case "EQUIP":
				if strings.Contains(msg, "equipped=") {
					parts := strings.Split(msg, "equipped=")
					if len(parts) > 1 {
						itemName := strings.TrimSpace(parts[1])
						if chat != nil {
							chat.AddLine(g.curr_cmd + "> " + "You equipped " + itemName)
						}
					}
				} else {
					if chat != nil {
						chat.AddLine(g.curr_cmd + "> " + "Equip failed")
					}
				}
				g.status = ""

			case "WHO":
				who, err := ParseWho(response)
				if err == nil {
					names := ""
					for _, name := range who.Room {
						names += name + " "
					}
					who_msg := fmt.Sprintf("Players: %s| Server: %d", names, who.Server)
					if chat != nil {
						chat.AddLine(g.curr_cmd + "> " + who_msg)
					}
				} else {
					if chat != nil {
						chat.AddLine("Server > " + msg)
					}
				}

			case "STATUS":
				status, err := ParseStatus(response)
				if err == nil {
					status_msg := fmt.Sprintf("You have %d/%d HP, You are %s", status.HP, status.MaxHP, status.Status)
					if chat != nil {
						chat.AddLine(g.curr_cmd + "> " + status_msg)
					}
				} else {
					if chat != nil {
						chat.AddLine("Server > " + msg)
					}
				}

//			case "GROUP":
//    		    world.ParseRoom(msg)
			case "QUIT":
				fmt.Printf("QUIT")
				return ErrQuit
			}
		}
		default:
			goto Done
		}
	}
Done:
	return nil
}
