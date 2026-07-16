package irc

import (
	"chatter-wails/internal/util"
	"context"
	"fmt"
	"log"
	"strings"
)

const IRC_ENDPOINT = "wss://irc-ws.chat.twitch.tv:443"

const COMMANDS_CAPABILITY = "twitch.tv/commands"
const MEMBERSHIP_CAPABILITY = "twitch.tv/membership"
const TAGS_CAPABILITY = "twitch.tv/tags"

// CLEARCHAT -> /clear /ban or /timeout
// CLEARMSG -> message delete


type IRCListener struct {
	// event -> funcs
	handlers map[string]map[*func()]struct{}
	socket   *util.Socket
}

func NewIRCListener() *IRCListener {
	return &IRCListener{
		handlers: map[string]map[*func()]struct{}{
			"CLEARCHAT": {},
			"CLEARMSG": {},
			"GLOBALUSERSTATE": {},
			"NOTICE": {},
			"PART": {},
			"PING": {},
			"PRIVMSG": {},
			"RECONNEC": {},
			"ROOMSTAT": {},
			"USERNOTICE": {},
			"USERSTAT": {},
		},
	}
}

func (irc *IRCListener) Connect(accessToken, userLogin string, commands, membership, tags bool) error {
	ctx := context.Background()
	socket, err := util.NewSocket(ctx, IRC_ENDPOINT, irc.handleSocketMessage)
	if err != nil {
		return err
	}

	irc.socket = socket

	irc.SendMessage(fmt.Sprintf("PASS oauth:%s", accessToken))
	irc.SendMessage((fmt.Sprintf("NICK %s", userLogin)))

	capabilities := []string{}
	if commands {
		capabilities = append(capabilities, COMMANDS_CAPABILITY)
	}
	if membership {
		capabilities = append(capabilities, MEMBERSHIP_CAPABILITY)
	}
	if tags {
		capabilities = append(capabilities, TAGS_CAPABILITY)
	}

	irc.RequestCapabilities(capabilities)

	return nil
}

func (irc *IRCListener) SendMessage(message string) {
	irc.socket.SendMessage([]byte(message))
}

func (irc *IRCListener) RequestCapabilities(capabilities []string) {
	if len(capabilities) == 0 { return }
	joined := strings.Join(capabilities, " ")
	irc.SendMessage(fmt.Sprintf("CAP REQ :%s", joined))
}

func (irc *IRCListener) JoinChannel(channelName string) {
	irc.SendMessage(fmt.Sprintf("JOIN #%s", channelName))
}

func (irc *IRCListener) PartChannel(channelName string) {
	irc.SendMessage(fmt.Sprintf("PART #%s", channelName))
}

// returns func to remove listener
func (irc *IRCListener) AddEventListener(event string, handler func()) (func(), error) {
	if _, exists := irc.handlers[event]; !exists {
		return nil, fmt.Errorf("Event %s doesn't exist in IRC", event)
	}

	irc.handlers[event][&handler] = struct{}{}

	var removeListener func()
	removeListener = func() {
		delete(irc.handlers[event], &handler)
	}

	return removeListener, nil
}

func (irc *IRCListener) handleSocketMessage(message []byte) {
	mStr := string(message)
	
	words := strings.Split(mStr, " ")
	if words[0] == "PING" {
		pong := "PONG "
		if len(words) > 1 {
			pong += words[1]
		}
		irc.socket.SendMessage([]byte(pong))
	}

	log.Printf("[IRC]: %s", mStr)
}
