package channelserver

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Character struct {
	ID              uint32 `db:"id"`
	IsFemale        bool   `db:"is_female"`
	IsNewCharacter  bool   `db:"is_new_character"`
	SmallGRLevel    uint8  `db:"small_gr_level"`
	GROverrideMode  bool   `db:"gr_override_mode"`
	Name            string `db:"name"`
	UnkDescString   string `db:"unk_desc_string"`
	GROverrideLevel uint16 `db:"gr_override_level"`
	GROverrideUnk0  uint8  `db:"gr_override_unk0"`
	GROverrideUnk1  uint8  `db:"gr_override_unk1"`
	Exp             uint16 `db:"exp"`
	Weapon          uint16 `db:"weapon"`
	LastLogin       uint32 `db:"last_login"`
}

var weapons = []string{
	"<:gs:970861408227049477>",
	"<:hbg:970861408281563206>",
	"<:hm:970861408239628308>",
	"<:lc:970861408298352660>",
	"<:sns:970861408319315988>",
	"<:lbg:970861408327725137>",
	"<:ds:970861408277368883>",
	"<:ls:970861408319311872>",
	"<:hh:970861408222863400>",
	"<:gl:970861408327720980>",
	"<:bw:970861408294174780>",
	"<:tf:970861408424177744>",
	"<:sw:970861408340283472>",
	"<:ms:970861408411594842>",
}

func (s *Server) getCharacterForUser(uid int) (*Character, error) {
	character := Character{}
	err := s.db.Get(&character, "SELECT id, is_female, is_new_character, small_gr_level, gr_override_mode, name, unk_desc_string, gr_override_level, gr_override_unk0, gr_override_unk1, exp, weapon, last_login FROM characters WHERE id = $1", uid)
	if err != nil {
		return nil, err
	}

	return &character, nil
}

func CountChars(s *Server) string {
	count := 0
	for _, stage := range s.stages {
		count += len(stage.clients)
	}

	message := fmt.Sprintf("Server [%s]: %d players;", s.name, count)

	return message
}

func PlayerList(s *Server) string {
	list := ""
	count := 0
	for _, stage := range s.stages {
		for client := range stage.clients {
			char, err := s.getCharacterForUser(int(client.charID))

			status := ""
			if stage.isCharInQuestByID(client.charID) {
				status = "(in quest)"
			}

			if err == nil {
				list = list + weapons[char.Weapon] + " " + char.Name + " " + status + "\n"
				count += 1
			}
		}
	}

	message := fmt.Sprintf("<:5658sus:969620902385946777> Invasores in Server: [%s ] <:5658sus:969620902385946777> \n========== Total %d ==========\n", s.name, count)
	message += list

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

	if commandName == "!status" && s.enable {
		ds.ChannelMessageSend(m.ChannelID, CountChars(s))
		return
	}

	if commandName == "!players" && s.enable {
		ds.ChannelMessageSend(m.ChannelID, PlayerList(s))
		return
	}

	if m.ChannelID == s.erupeConfig.Discord.RealtimeChannelID {
		message := fmt.Sprintf("[DISCORD] %s: %s", m.Author.Username, s.discordBot.NormalizeDiscordMessage(m.Content))
		s.BroadcastChatMessage(message)
	}
}
