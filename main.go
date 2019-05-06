package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type listStrValue struct {
	strings *[]string
}

func (v listStrValue) String() string {
	return strings.Join(*v.strings, ",")
}
func (v listStrValue) Set(value string) error {
	*v.strings = strings.Split(value, ",")
	return nil
}

var (
	token     string
	channelID string
	question  string
	choices   []string
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&channelID, "c", "", "Channel ID")
	flag.StringVar(&question, "q", "", "Question")
	flag.Var(listStrValue{&choices}, "cs", "Choices")
	flag.Parse()
}

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

func main() {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("creating Discord session failed: ", err)
		return
	}

	err = sendPoll(s, channelID, question, choices)
	if err != nil {
		log.Fatalln("sending poll failed: ", err)
		return
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
