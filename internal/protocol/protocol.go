package protocol

import "fmt"

const (
	CmdConnect   = "CONNECT"
	CmdLook      = "LOOK"
	CmdMove      = "MOVE"
	CmdQuit      = "QUIT"
	CmdChat      = "CHAT"
	CmdWho       = "WHO"
	CmdGroup     = "GROUP"
	CmdTake      = "TAKE"
	CmdDrop      = "DROP"
	CmdInventory = "INVENTORY"
	CmdTalk      = "TALK"
	CmdAttack    = "ATTACK"
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
	MsgBye       = "OK bye\n"
)

const (
	CodeNameInUse          = 201
	CodeNoExit             = 301
	CodeUnknownCommand	   = 400
	CodeNotInGroup         = 401
	CodeAlreadyInGroup     = 402
	CodeItemNotFound       = 404
	CodeItemNotInInventory = 404
	CodeNPCNotFound        = 404
	CodeNPCNotHostile      = 405
	CodeNoQuestAvailable   = 406
	CodeAlreadyAuthenticated = 407
	CodeNotAuthenticated 	= 408
	CodeConnectionFailed   = 900
	CodeSendFailed         = 901
)

const (
	MsgNameInUse          = "NAME_IN_USE"
	MsgNoExit             = "NO_EXIT"
	MsgUnknownCommand	  = "UNKNOWN_COMMAND"
	MsgNotInGroup         = "NOT_IN_GROUP"
	MsgAlreadyInGroup     = "ALREADY_IN_GROUP"
	MsgItemNotFound       = "ITEM_NOT_FOUND"
	MsgItemNotInInventory = "ITEM_NOT_IN_INVENTORY"
	MsgNPCNotFound        = "NPC_NOT_FOUND"
	MsgNPCNotHostile      = "NPC_NOT_HOSTILE"
	MsgNoQuestAvailable   = "NO_QUEST_AVAILABLE"
	MsgAlreadyAuthenticated = "ALREADY_AUTHENTICATED"
	MsgNotAuthenticated		= "NOT_AUTHENTICATED"
	MsgConnectionFailed   = "CONNECTION_FAILED"
	MsgSendFailed         = "SEND_FAILED"
)

func FormatCmd(commandName, arguments string) string {
	if arguments == "" {
		return fmt.Sprintf("%s\n", commandName)
	}
	return fmt.Sprintf("%s %s\n", commandName, arguments)
}

func FormatOK(responseData string) string {
	if responseData == "" {
		return "OK\n"
	}
	return fmt.Sprintf("OK %s\n", responseData)
}

func FormatErr(errorCode int, errorMessage string) string {
	return fmt.Sprintf("ERR %03d %s\n", errorCode, errorMessage)
}

func FormatEvt(eventType, eventData string) string {
	return fmt.Sprintf("EVT %s %s\n", eventType, eventData)
}

type State int
const (
	// est un identificateur pré-déclaré qui agit comme un compteur automatique
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


