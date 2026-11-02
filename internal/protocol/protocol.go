package protocol

import "fmt"

const (
	CmdConnect   = "CONNECT"
	CmdLook      = "LOOK"
	CmdMove      = "MOVE"
	CmdQuit      = "QUIT"
	CmdChat      = "CHAT"
	CmdWho       = "WHO"
	CmdWhoAll    = "WHO_ALL"
	CmdGroup     = "GROUP"
	CmdTake      = "TAKE"
	CmdDrop      = "DROP"
	CmdInventory = "INVENTORY"
	CmdEquip     = "EQUIP"
	CmdTalk      = "TALK"
	CmdAttack    = "ATTACK"
	CmdDefend    = "DEFEND"
	CmdFlee      = "FLEE"
	CmdStatus    = "STATUS"
	CmdQuest     = "QUEST"
	CmdQuests    = "QUESTS"
)

const (
	GroupCreate = "CREATE"
	GroupInvite = "INVITE"
	GroupJoin   = "JOIN"
	GroupLeave  = "LEAVE"
)

const (
	ChatGlobal = "GLOBAL"
	ChatRoom   = "ROOM"
	ChatGroup  = "GROUP"
)

const (
	PrefixOK  = "OK"
	PrefixERR = "ERR"
	PrefixEVT = "EVT"

	MsgHello     = "OK hello proto=1\n"
	MsgConnected = "OK connected\n"
	MsgBye       = "bye"
)

const (
	EvtRoomEnter   = "ROOM PRESENCE ENTER"
	EvtRoomLeave   = "ROOM PRESENCE LEAVE"
	EvtGroupInvite = "GROUP INVITE"
	EvtGroupJoin   = "GROUP JOIN"
	EvtGroupLeave  = "GROUP LEAVE"
	EvtPlayers     = "STATS"
	EvtChatGlobal  = "GLOBAL CHAT"
	EvtChatRoom    = "ROOM CHAT"
	EvtChatGroup   = "GROUP CHAT"
)

const (
	StatusCombat          = "combat"
	StatusHealthy         = "healthy"
	StatusPerfectDefense  = "perfect_defense"
	StatusWeakenedDefense = "weakened_defense"
	StatusBrokenDefense   = "broken_defense"
	StatusFlee            = "flee"
	StatusDefeated        = "defeated"
	StatusVictory         = "victory"
)

const (
	CodeNameInUse          = 201
	CodeNoExit             = 301
	CodeBadRequest         = 400
	CodeNotInGroup         = 401
	CodeAlreadyInGroup     = 402
	CodeUnknownPlayer      = 403
	CodeItemNotFound       = 404
	CodeItemNotInInventory = 404
	CodeNPCNotFound        = 404
	CodeNPCNotHostile      = 405
	CodeNoQuestAvailable   = 406
	CodeUnknownGroup       = 407
	CodeConnectionFailed   = 900
	CodeSendFailed         = 901
)

const (
	MsgNameInUse          = "NAME_IN_USE"
	MsgNoExit             = "NO_EXIT"
	MsgBadRequest         = "BAD_REQUEST"
	MsgNotInGroup         = "NOT_IN_GROUP"
	MsgAlreadyInGroup     = "ALREADY_IN_GROUP"
	MsgUnknownPlayer      = "UNKNOWN_PLAYER"
	MsgItemNotFound       = "ITEM_NOT_FOUND"
	MsgItemNotInInventory = "ITEM_NOT_IN_INVENTORY"
	MsgUnknownGroup       = "UNKNOWN_GROUP"
	MsgNPCNotFound        = "NPC_NOT_FOUND"
	MsgNPCNotHostile      = "NPC_NOT_HOSTILE"
	MsgNoQuestAvailable   = "NO_QUEST_AVAILABLE"
	MsgConnectionFailed   = "CONNECTION_FAILED"
	MsgSendFailed         = "SEND_FAILED"

	MsgArgsError            = "ARGS_ERROR"
	MsgUnknownCommand       = "UNKNOWN_COMMAND"
	MsgNotAuthenticated     = "NOT_AUTHENTICATED"
	MsgAlreadyAuthenticated = "ALREADY_AUTHENTICATED"
	MsgNoInvitation         = "NO_INVITATION"
	MsgGroupNotFound        = "GROUP_NOT_FOUND"
	MsgMoveError            = "NO_MOVE_DURING_COMBAT"
)

// FormatCmd formats a client command with optional arguments into a newline-terminated string.
func FormatCmd(commandName, arguments string) string {
	if arguments == "" {
		return fmt.Sprintf("%s\n", commandName)
	}
	return fmt.Sprintf("%s %s\n", commandName, arguments)
}

// FormatOK formats a success response, appending data if provided.
func FormatOK(responseData string) string {
	if responseData == "" {
		return "OK\n"
	}
	return fmt.Sprintf("OK %s\n", responseData)
}

// FormatErr formats an error response with a 3-digit error code and message.
func FormatErr(errorCode int, errorMessage string) string {
	return fmt.Sprintf("ERR %03d %s\n", errorCode, errorMessage)
}

// FormatEvt formats a server push event.
func FormatEvt(eventType, eventData string) string {
	return fmt.Sprintf("EVT %s %s\n", eventType, eventData)
}

type State int

const (
	StateDisconnected State = iota
	StateConnected
	StateAuthenticated
	StateTerminated
)

func (s State) String() string {
	switch s {
	case StateDisconnected:
		return "DISCONNECTED"
	case StateConnected:
		return "CONNECTED"
	case StateAuthenticated:
		return "AUTHENTICATED"
	case StateTerminated:
		return "TERMINATED"
	default:
		return "UNKNOWN"
	}
}
