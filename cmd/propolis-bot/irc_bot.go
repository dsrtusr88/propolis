package main

// forked from github.com/go-chat-bot/bot/irc

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/go-chat-bot/bot"
	ircevent "github.com/thoj/go-ircevent"
)

// BotConfig must contain the necessary data to connect to an IRC server.
type BotConfig struct {
	Server        string   // IRC server:port. Ex: ircevent.freenode.org:7000
	Channels      []string // Channels to connect. Ex: []string{"#go-bot", "#channel mypassword"}
	User          string   // The IRC username the bot will use
	Nick          string   // The nick the bot will use
	RealName      string   // The real name (longer string) the bot will use
	Password      string   // nickserv password
	Key           string   // IRC Key
	GateKeeper    string   // user access bot
	TLSServerName string   // Must supply if UseTLS is true
	UseTLS        bool     // Should connect using TLS?
	Debug         bool     // This will log all IRC communication to standad output
}

var (
	ircConn     *ircevent.Connection
	b           *bot.Bot
	nickStartRE *regexp.Regexp
)

const protocol = "irc"

func responseHandler(target string, message string, sender *bot.User) {
	channel := target
	if ircConn.GetNick() == target {
		channel = sender.Nick
	}

	if message != "" {
		for _, line := range strings.Split(message, "\n") {
			ircConn.Privmsg(channel, line)
		}
	}
}

func getServerName(server string) string {
	separatorIndex := strings.LastIndex(server, ":")
	if separatorIndex != -1 {
		return server[:separatorIndex]
	}
	return server
}

// SetUp returns a bot for irc according to the BotConfig, but does not run it.
// When you are ready to run the bot, call Run(nil).
// This is useful if you need a pointer to the bot, otherwise you can simply call Run().
func SetUp(c *BotConfig) *bot.Bot {
	ircConn = ircevent.IRC(c.User, c.Nick)
	ircConn.RealName = c.RealName
	ircConn.Password = c.Password
	ircConn.UseTLS = c.UseTLS
	ircConn.TLSConfig = &tls.Config{
		ServerName: getServerName(c.Server),
	}
	ircConn.VerboseCallbackHandler = c.Debug

	b = bot.New(&bot.Handlers{
		Response: responseHandler,
	},
		&bot.Config{
			Protocol: protocol,
			Server:   c.Server,
		},
	)
	// prepare regex to strip from messages - nick followed by colon/comma and spaces
	nickStartRE = regexp.MustCompile(fmt.Sprintf("%s[,:] *", c.Nick))

	ircConn.AddCallback("001", func(_ *ircevent.Event) {
		ircConn.Privmsg("NickServ", "IDENTIFY "+c.Password)
		time.Sleep(1000 * time.Millisecond)
		ircConn.Privmsg(c.GateKeeper, fmt.Sprintf("enter %s %s %s", c.Channels[0], c.Nick, c.Key))
		time.Sleep(1000 * time.Millisecond)
		for _, channel := range c.Channels {
			ircConn.Join(channel)
		}
	})
	ircConn.AddCallback("PRIVMSG", func(e *ircevent.Event) {
		if e.Nick != c.Nick {
			fmt.Println("Ignoring command from " + e.Nick)
			return
		}
		b.MessageReceived(
			&bot.ChannelData{
				Protocol:  protocol,
				Server:    ircConn.Server,
				Channel:   e.Arguments[0],
				IsPrivate: e.Arguments[0] == ircConn.GetNick()},
			&bot.Message{
				Text: nickStartRE.ReplaceAllString(e.Message(), ""),
			},
			&bot.User{
				ID:       e.Host,
				Nick:     e.Nick,
				RealName: e.User})
	})
	return b
}

// SetUpConn wraps SetUp and returns ircConn in addition to bot.
func SetUpConn(c *BotConfig) (*bot.Bot, *ircevent.Connection) {
	return SetUp(c), ircConn
}

// Run reads the BotConfig, connect to the specified IRC server and starts the bot.
// The bot will automatically join all the channels specified in the configuration.
func Run(c *BotConfig) {
	if c != nil {
		SetUp(c)
	}

	err := ircConn.Connect(c.Server)
	if err != nil {
		log.Fatal(err)
	}
	ircConn.Loop()
}
