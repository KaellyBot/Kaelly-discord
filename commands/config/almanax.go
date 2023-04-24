package config

import (
	"context"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-discord/models/mappers"
)

func (command *ConfigCommand) almanaxRequest(ctx context.Context, s *discordgo.Session,
	i *discordgo.InteractionCreate, lg discordgo.Locale) {

	channelId, enabled, locale, err := command.getWebhookAlmanaxOptions(ctx)
	if err != nil {
		panic(err)
	}

	webhook, err := command.createWebhook(s, channelId)
	if err != nil {
		panic(err)
		// TODO
	}

	msg := mappers.MapConfigurationWebhookAlmanaxRequest(webhook, enabled, locale, lg)
	err = command.requestManager.Request(s, i, configurationRequestRoutingKey, msg, command.almanaxRespond)
	if err != nil {
		panic(err)
	}
}

func (command *ConfigCommand) almanaxRespond(ctx context.Context, s *discordgo.Session,
	i *discordgo.InteractionCreate, message *amqp.RabbitMQMessage, properties map[string]any) {

	// TODO respond
}
