package gui

import (
	"net"
	"tap/internal/gui/components"
)

func (g *Game) sendServerDirection(conn net.Conn) {
	if buttons := g.centerMid.GetButtons(); buttons != nil {
		for _, button := range buttons {
			btn := button
			switch btn.Label {
			case "North":
				btn.OnClick = func() {
					if g.status == "AWAITING_MOVE_DIR" {
						msg := "MOVE North"
						SendMessage(conn, msg)
						g.direction = 0
						g.status = ""
					}
				}
			case "East":
				btn.OnClick = func() {
					if g.status == "AWAITING_MOVE_DIR" {
						msg := "MOVE East"
						SendMessage(conn, msg)
						g.direction = 1
						g.status = ""
					}
				}
			case "South":
				btn.OnClick = func() {
					if g.status == "AWAITING_MOVE_DIR" {
						msg := "MOVE South"
						SendMessage(conn, msg)
						g.direction = 2
						g.status = ""
					}
				}
			case "West":
				btn.OnClick = func() {
					if g.status == "AWAITING_MOVE_DIR" {
						msg := "MOVE West"
						SendMessage(conn, msg)
						g.direction = 3
						g.status = ""
					}
				}
			}
		}
	}
}

func (g *Game) sendActionButton(conn net.Conn, chat *components.ChatWidget) {
	if buttons := g.centerRight.GetButtons(); buttons != nil {
		for _, button := range buttons {
			btn := button

			switch btn.Label {
			case "MOVE":
				btn.OnClick = func() {
					g.status = "AWAITING_MOVE_DIR"
					g.curr_cmd = "MOVE"
				}

			case "TAKE":
				btn.OnClick = func() {
					g.status = "AWAITING_TAKE_ITEM"
					g.curr_cmd = "TAKE"
				}

			case "DROP":
				btn.OnClick = func() {
					g.status = "AWAITING_DROP_ITEM"
					g.curr_cmd = "DROP"
				}

			case "EQUIP":
				btn.OnClick = func() {
					g.status = "AWAITING_EQUIP_ITEM"
					g.curr_cmd = "EQUIP"
				}

			case "TALK":
				btn.OnClick = func() {
					g.status = "AWAITING_TALK_NPC"
					g.curr_cmd = "TALK"
				}

			case "ATTACK":
				btn.OnClick = func() {
					g.status = "AWAITING_ATTACK_NPC"
					g.curr_cmd = "ATTACK"
				}

			case "DEFEND":
				btn.OnClick = func() {
					g.curr_cmd = "DEFEND"
					SendMessage(conn, "DEFEND")
				}

			case "STATUS":
				btn.OnClick = func() {
					g.status = ""
					g.curr_cmd = "STATUS"
					SendMessage(conn, "STATUS")
				}

			case "QUEST":
				btn.OnClick = func() {
					g.status = "AWAITING_QUEST_NPC"
					g.curr_cmd = "QUEST"
				}

			case "QUESTS":
				btn.OnClick = func() {
					g.status = "AWAITING_QUEST_NPC"
					g.curr_cmd = "QUESTS"
					if chat != nil {
						chat.AddLine(g.name + " > QUESTS")
					}
					SendMessage(conn, "QUESTS")
				}

			case "WHO":
				btn.OnClick = func() {
					g.status = "WHO"
					g.curr_cmd = "WHO"
					SendMessage(conn, "WHO")
				}

			case "GROUP":
				btn.OnClick = func() {
					g.status = "GROUP"
					g.curr_cmd = "GROUP"
					if chat != nil {
						chat.AddLine(g.name + " > GROUP")
					}
					SendMessage(conn, "GROUP")
				}

			default:
				btn.OnClick = func() {
					g.status = "IDLE"
					g.curr_cmd = btn.Label
					if g.curr_cmd != "LOOK" {
						if chat != nil {
							chat.AddLine(g.name + " > " + btn.Label)
						}
					}
					SendMessage(conn, btn.Label)
				}
			}
		}
	}
}


func (g *Game) sendGroupButton(conn net.Conn, chat *components.ChatWidget) {
	if buttons := g.centerLeft.GetButtons(); buttons != nil {
		for _, button := range buttons {
			btn := button

			switch btn.Label {
			case "CREATE":
				btn.OnClick = func() {
					msg := "group create"
					SendMessage(conn, msg)
				}

			case "LEAVE":
				btn.OnClick = func() {
					msg := "group leave"
					SendMessage(conn, msg)
				}
			}
		}
	}
}
