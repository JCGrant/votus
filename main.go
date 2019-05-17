package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

var reactions = []string{
	"🇦",
	"🇧",
	"🇨",
	"🇩",
	"🇪",
	"🇫",
	"🇬",
	"🇭",
	"🇮",
	"🇯",
	"🇰",
	"🇱",
	"🇲",
	"🇳",
	"🇴",
	"🇵",
	"🇶",
	"🇷",
	"🇸",
	"🇹",
}

// Question represents a single question for Votus to ask
type Question struct {
	Text    string
	Choices []string
}

// Config for Votus
type Config struct {
	Token     string
	Admins    []string
	Questions []Question
}

// Bot represents the Votus bot
type Bot struct {
	Config
	session *discordgo.Session
}

func main() {
	configFileName := "config.toml"
	conf, err := LoadConfig(configFileName)
	if err != nil {
		log.Fatalln("loading config failed: ", err)
	}
	b, err := New(conf)
	if err != nil {
		log.Fatalln("creating bot failed: ", err)
	}
	defer b.Stop()

	err = b.Run()
	if err != nil {
		log.Fatalln("running bot failed: ", err)
	}
}

// LoadConfig loads the config
func LoadConfig(filename string) (conf Config, err error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file failed: %s", err)
	}
	if _, err := toml.Decode(string(bs), &conf); err != nil {
		return Config{}, fmt.Errorf("decoding toml failed: %s", err)
	}
	return conf, nil
}

// New makes a new bot
func New(conf Config) (*Bot, error) {
	session, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		return nil, fmt.Errorf("creating Discord session failed: %s", err)
	}
	b := &Bot{conf, session}
	session.AddHandler(b.sendPollHandler)
	return b, nil
}

// Run runs the bot
func (b *Bot) Run() error {
	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("opening connection failed: %s", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Votus is now running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-c
	log.Println("Shutting down Votus...")
	return nil
}

// Stop stops the bot
func (b *Bot) Stop() {
	b.session.Close()
}

func (b *Bot) sendPollHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Ignore all messages not sent by an admin
	if !stringInSlice(m.Author.Username, b.Admins) {
		return
	}

	if m.Content == "votus" {
		err := sendPoll(s, m.ChannelID, b.Questions[0].Text, b.Questions[0].Choices)
		if err != nil {
			log.Fatalln("sending poll failed: ", err)
		}
	}
}

func sendPoll(s *discordgo.Session, channelID string, question string, choices []string) error {
	msgContent := "@here " + question + "\n\n"

	for i, choice := range choices {
		reaction := string(reactions[i])
		msgContent += fmt.Sprintf("%s: %s\n", reaction, choice)
	}

	msg, err := s.ChannelMessageSend(channelID, msgContent)
	if err != nil {
		return fmt.Errorf("sending message '%s' failed: %s", msgContent, err)
	}

	for i := range choices {
		reaction := string(reactions[i])
		err := s.MessageReactionAdd(channelID, msg.ID, reaction)
		if err != nil {
			return fmt.Errorf("adding reaction '%s' failed: %s", reaction, err)
		}
	}
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
