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
	handlers map[string]map[*func(*IRCMessage)]struct{}
	socket   *util.Socket
}

func NewIRCListener() *IRCListener {
	listener := &IRCListener{
		handlers: map[string]map[*func(*IRCMessage)]struct{}{
			"CLEARCHAT": {},
			"CLEARMSG": {},
			"GLOBALUSERSTATE": {},
			"USERSTATE": {},
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

	pingHandler := listener.handlePing
	listener.handlers["PING"][&pingHandler] = struct{}{}

	return listener
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
func (irc *IRCListener) AddEventListener(event string, handler func(*IRCMessage)) (func(), error) {
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
	// [@tags(comma-seperated) ] :sender COMMAND #channel :data
	ircMessage := parseIRCMessage(mStr)

	irc.dispatchHandlers(ircMessage)
}

func (irc *IRCListener) handlePing(message *IRCMessage) {
	irc.socket.SendMessage(fmt.Appendf([]byte{}, "PONG :%s", message.Data))
}

func (irc *IRCListener) dispatchHandlers(message *IRCMessage) {
	handlers := irc.handlers[message.Command]
	if handlers == nil { 
		log.Printf("[irc.dispatchHandlers]: Not handling %s", message.Command)
		return 
	}

	for fn := range handlers {
		(*fn)(message)
	}
}

type IRCMessage struct{
	Tags map[string]string
	Sender string
	Command string
	Channel string
	Data string
}

func parseIRCMessage(message string) *IRCMessage {
	tags, nextIndex := parseTags(message, 0)
	sender, nextIndex := parseSender(message, nextIndex)
	command, nextIndex := parseCommand(message, nextIndex)
	channel, nextIndex := parseChannel(message, nextIndex)
	data, _ := parseData(message, nextIndex)

	return &IRCMessage{
		Tags: tags,
		Sender: sender,
		Command: command,
		Channel: channel,
		Data: data,
	}
}

func parseTags(message string, start int) (map[string]string, int) {
	tags := map[string]string{}
	if start >= len(message) || message[start] != '@' {
		return tags, start
	}

	nextSpaceIndex := getNextSpace(message, start)

	for kvstr := range strings.SplitSeq(message[start+1:nextSpaceIndex], ";") {
		kv := strings.Split(kvstr, "=")
		key := kv[0]
		var val string
		if len(kv) > 1 {
			val = kv[1]
		}

		tags[key] = val
	}

	return tags, nextSpaceIndex+1
}

func parseSender(message string, start int) (string, int) {
	if start >= len(message) || message[start] != ':' {
		return "", start
	}

	nextSpaceIndex := getNextSpace(message, start)
	return message[start:nextSpaceIndex], nextSpaceIndex+1
}

func parseCommand(message string, start int) (string, int) {
	if start >= len(message) {
		return "", start
	}

	nextSpaceIndex := getNextSpace(message, start)
	return message[start:nextSpaceIndex], nextSpaceIndex+1
}

func parseChannel(message string, start int) (string, int) {
	if start >= len(message) || message[start] != '#' {
		return "", start
	}

	nextSpaceIndex := getNextSpace(message, start)

	return message[start+1:nextSpaceIndex], nextSpaceIndex+1
}

func parseData(message string, start int) (string, int) {
	if start >= len(message) || message[start] != ':' {
		return "", start
	}

	return message[start+1:], len(message)
}

func getNextSpace(message string, start int) int {
	nextSpaceIndex := strings.Index(message[start:], " ")
	if nextSpaceIndex == -1 {
		return len(message)
	}

	return start + nextSpaceIndex
}
