package main

import (
	"fmt"
	"io/ioutil"
	"log"

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
	ChannelID string `toml:"channel_id"`
	Questions []Question
}

func main() {
	configFileName := "config.toml"
	var conf Config
	bs, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatalln("reading config file failed: ", err)
	}
	if _, err := toml.Decode(string(bs), &conf); err != nil {
		log.Fatalln("loading config failed: ", err)
	}

	session, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Fatalln("creating Discord session failed: ", err)
	}

	err = sendPoll(session, conf.ChannelID, conf.Questions[0].Text, conf.Questions[0].Choices)
	if err != nil {
		log.Fatalln("sending poll failed: ", err)
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
