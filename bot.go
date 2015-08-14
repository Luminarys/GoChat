package gochat

import (
	"errors"
	"os"
)

type Bot struct {
	Server   string
	Nick     string
	Channels map[string]*Channel
	Modules  []Module
	Conn     *connection
}

//Creates a new bot for a server, and returns it once it is ready
func NewBot(server, nick string, hijack bool) (*Bot, error) {
	logInit(os.Stdout, os.Stderr, os.Stderr)
	var conn *connection
	var err error

	bot := &Bot{
		Server:   server,
		Nick:     nick,
		Channels: make(map[string]*Channel),
		Modules:  make([]Module, 0),
	}

	conn, err = makeConn(server, nick, false, hijack) //Create new irc connection
	if err != nil {
		return nil, errors.New("Error! Could not connect")
	}
	bot.Conn = conn
	/*ready once we get the welcome message(001)
	conn.AddCallback("001", func(e *irc.Event) {
		ready <- true
	})
	close(ready)*/

	readyChan := make(chan bool)
	go bot.handleMessages(readyChan)

	bot.Conn.user(nick)
	bot.Conn.nick(nick)

	<-readyChan
	close(readyChan)

	LTrace.Println("Succesfully intialized a bot!")
	return bot, nil
}

func (bot *Bot) handleMessages(ready chan bool) {
	r := false
	LTrace.Println("Starting message handling loop")
	for msg := range bot.Conn.ReadChan {
		LTrace.Println("Receieved message: ", msg)
		if msg.Cmd == "PRIVMSG" {
			//Check that channel is valid
			if _, ok := bot.Channels[msg.Arguments[0]]; ok {
				bot.Channels[msg.Arguments[0]].HandleMessage(msg)
			} else {
				//Otherwise it's probably a user, so instantiate a
				//temporary user channel to handle things
				c := bot.NewChannel(msg.Arguments[0])
				c.HandleMessage(msg)
			}
		} else if msg.Cmd == "MODE" {
			if _, ok := bot.Channels[msg.Arguments[0]]; ok {
				bot.Channels[msg.Arguments[0]].ModeChange(msg)
			} else if !r {
				LTrace.Println("Got mode message")
				ready <- true
				r = true
			}
		} else if msg.Cmd == "353" {
			bot.Channels[msg.Arguments[2]].SetOps(msg.Text)
		}
	}
}

//Loads a module into the bot
func (bot *Bot) AddModule(mod Module) {
	bot.Modules = append(bot.Modules, mod)
}

//Joins a channel
func (bot *Bot) JoinChan(chanName string) *Channel {
	LTrace.Println("Creating and registering a channel interally")
	c := bot.NewChannel(chanName)
	bot.Channels[chanName] = c

	LTrace.Println("Sending the JOIN message to the server")
	bot.Conn.send("JOIN " + chanName)

	for !c.Ready {
	}
	LTrace.Println("Channel is ready!")

	return c
}

func (bot *Bot) Register(pass string, email string) {
	bot.Conn.register(pass, email)
}

//Disconnects and destroys the bot
func (bot *Bot) Quit() {
	bot.Conn.quit()
}

//Broadcasts a message to all chans
func (bot *Bot) Broadcast(message string) {
	for _, c := range bot.Channels {
		bot.Conn.privmsg(c.Name, message)
	}
}

//Leaves a channel and destroys it
func (bot *Bot) Part(channel string) {
	bot.Conn.send("PART " + channel)
	bot.Channels[channel] = nil
}
