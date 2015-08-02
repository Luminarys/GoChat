package gochat

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}
	bot.Quit()
}

func TestChanJoin(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	c.Part()
	bot.Quit()
}

func TestBroadcast(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}
	c := bot.JoinChan("#go-bot-test")
	bot.Broadcast("broadcast test")
	c.Part()
	bot.Quit()
}

func TestMessage(t *testing.T) {
	bot, err := NewBot("irc.rizon.net:6666", "go-bot")
	if err != nil {
		t.Error(err.Error())
	}

	c := bot.JoinChan("#go-bot-test")
	c.Say("Message test")
	time.Sleep(7 * time.Second)
	c.Part()
	bot.Quit()
}