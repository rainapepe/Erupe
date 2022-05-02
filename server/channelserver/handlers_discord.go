package channelserver

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func CountChars(s *Server) string {
	count := 0
	for _, stage := range s.stages {
		count += len(stage.clients)
	}

	message := fmt.Sprintf("Server [%s]: %d players;", s.name, count)

	return message
}

// onDiscordMessage handles receiving messages from discord and forwarding them in game.
func (s *Server) onDiscordMessage(ds *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore messages from our bot
	if m.Author.ID == ds.State.User.ID {
		return
	}

	// Split on comma.
	args := strings.Split(m.Content, " ")
	commandName := args[0]

	if commandName == "!status" {
		ds.ChannelMessageSend(m.ChannelID, CountChars(s))
	}

	if m.ChannelID == s.erupeConfig.Discord.RealtimeChannelID {
		message := fmt.Sprintf("[DISCORD] %s: %s", m.Author.Username, s.discordBot.NormalizeDiscordMessage(m.Content))
		s.BroadcastChatMessage(message)
	}
}
