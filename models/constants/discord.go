package constants

import (
	"github.com/bwmarrin/discordgo"
)

const (
	InvisibleCharacter     = "\u200b"
	MaxButtonPerActionRow  = 5
	MaxCharacterPerField   = 10
	MaxEquipmentPerField   = 8
	MaxIngredientsPerField = 8
)

func GetIntents() discordgo.Intent {
	return discordgo.IntentMessageContent |
		discordgo.IntentGuildMembers |
		discordgo.IntentGuilds |
		discordgo.IntentGuildMessages |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildWebhooks
}
