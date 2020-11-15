package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	fmt.Println("+ Testing Config...")
	check := assert.New(t)

	c := &Config{}
	err := c.Load("../../test/config.yaml")
	check.Nil(err)

	// general
	fmt.Println("Checking general")
	check.Equal(3, c.General.LogLevel)

	// irc
	check.True(c.IrcConfigured)
	check.Equal("irc.server.net:6697", c.IRC.Server)
	check.Equal("kkeeyy", c.IRC.Key)
	check.True(c.IRC.UseSSL)
	check.Equal("something", c.IRC.NickServPassword)
	check.Equal("user", c.IRC.Username)
	check.Equal("mybot", c.IRC.BotName)
	check.Equal("Bee", c.IRC.GateKeeper)
	check.Equal(defaultChannel, c.IRC.Channel)
	check.Equal(nodeRole, c.IRC.Role)
	check.Equal(centralBot, c.IRC.CentralBot)

	// tracker
	check.True(c.TrackerConfigured)
	check.Equal("blue", c.Tracker.Site)
	check.Equal("mytoken", c.Tracker.Token)
	check.Equal("https://blue.it", c.Tracker.URL)
	check.Equal("apkikey", c.Tracker.APIKey)
	check.Equal(8080, c.Tracker.Port)
	check.Equal([]string{"thisguy", "thisotherguy"}, c.Tracker.BlacklistedUploaders)
	check.Equal([]string{"tag1", "tag2"}, c.Tracker.ExcludedTags)
	check.Equal([]string{"niceguy:2", "veryniceperson:-1"}, c.Tracker.Users)
	check.Equal(2, len(c.Tracker.WhitelistedUsers))
	check.Equal(2, c.Tracker.WhitelistedUsers["niceguy"])
	check.Equal(-1, c.Tracker.WhitelistedUsers["veryniceperson"])

	fmt.Println(c.String())
}
