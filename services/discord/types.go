package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-discord/commands"
	"github.com/kaellybot/kaelly-discord/services/guilds"
)

const (
	welcomeMessagePermissions = discordgo.PermissionSendMessages | discordgo.PermissionEmbedLinks
)

var (
	ErrInvalidInteractionType = errors.New("interaction type is not handled")
)

type Service interface {
	Listen() error
	IsConnected() bool
	Shutdown()
}

type Impl struct {
	session      *discordgo.Session
	commands     []commands.DiscordCommand
	guildService guilds.Service
	broker       amqp.MessageBroker
}
