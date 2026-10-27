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
						g.curr_cmd = "MOVE North"
						sendMessage(conn, g.curr_cmd)
						g.direction = 0
						g.status = ""
					}
				}
			case "East":
				btn.OnClick = func() {
					if g.status == "AWAITING_MOVE_DIR" {
						g.curr_cmd = "MOVE East"
						sendMessage(conn, g.curr_cmd)
						g.direction = 1
						g.status = ""
					}
				}
			case "South":
				btn.OnClick = func() {
					if g.status == "AWAITING_MOVE_DIR" {
						g.curr_cmd = "MOVE South"
						sendMessage(conn, g.curr_cmd)
						g.direction = 2
						g.status = ""
					}
				}
			case "West":
				btn.OnClick = func() {
					if g.status == "AWAITING_MOVE_DIR" {
						g.curr_cmd = "MOVE West"
						sendMessage(conn, g.curr_cmd)
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
					sendMessage(conn, "DEFEND")
				}

			case "STATUS":
				btn.OnClick = func() {
					g.status = ""
					g.curr_cmd = "STATUS"
					sendMessage(conn, "STATUS")
				}

			case "QUEST":
				btn.OnClick = func() {
					g.status = "AWAITING_QUEST_NPC"
					g.curr_cmd = "QUEST"
				}

			case "QUESTS":
				btn.OnClick = func() {
					g.status = ""
					g.curr_cmd = "QUESTS"
					if chat != nil {
						chat.AddLine(g.name + " > QUESTS")
					}
					sendMessage(conn, "QUESTS")
				}

			case "WHO":
				btn.OnClick = func() {
					g.status = ""
					g.curr_cmd = "WHO"
					sendMessage(conn, "WHO")
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
					sendMessage(conn, btn.Label)
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
					g.curr_cmd = "CREATE GROUP"
					sendMessage(conn, "group create")
				}

			case "LEAVE":
				btn.OnClick = func() {
					g.curr_cmd = "LEAVE GROUP"
					sendMessage(conn, "group leave")
				}
			}
		}
	}
}
