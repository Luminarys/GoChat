package gcModules

import (
	"github.com/Luminarys/gochat"
	"strings"
)

//Allows for Ignores/Unignores of the specified Nick. Good for dealing with bot feedback loops and spammers. Requires operator status
type IgnoreMod struct {
}

func (m *IgnoreMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	parts := strings.Split(msg.Text, " ")
	if isOp, regd := c.Ops[msg.Nick]; isOp && regd {
		return (parts[0] == ".ignore" || parts[0] == ".unignore") && len(parts) > 1
	}
	return false
}

func (m *IgnoreMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	parts := strings.Split(msg.Text, " ")

	if isOp, regd := c.Ops[parts[1]]; isOp && regd {
		return "Sorry, but I can't ignore operators"
	}
	if parts[0] == ".ignore" {
		c.IgnoreNick(parts[1])
		return "Ignored " + parts[1]
	} else if parts[0] == ".unignore" {
		c.UnignoreNick(parts[1])
		return "Unignored " + parts[1]
	}
	return ""
}