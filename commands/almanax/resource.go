package almanax

import (
	"context"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-discord/commands"
	"github.com/kaellybot/kaelly-discord/models/constants"
	"github.com/kaellybot/kaelly-discord/models/mappers"
	"github.com/kaellybot/kaelly-discord/utils/middlewares"
	"github.com/rs/zerolog/log"
)

func (command *Command) getResources(ctx context.Context, s *discordgo.Session,
	i *discordgo.InteractionCreate, _ middlewares.NextFunc) {
	duration, err := getDurationOption(ctx)
	if err != nil {
		panic(err)
	}

	msg := mappers.MapAlmanaxResourceRequest(duration, i.Locale)
	err = command.requestManager.Request(s, i, almanaxRequestRoutingKey, msg, command.resourceRespond)
	if err != nil {
		panic(err)
	}
}

func (command *Command) resourceRespond(_ context.Context, s *discordgo.Session,
	i *discordgo.InteractionCreate, message *amqp.RabbitMQMessage, _ map[string]any) {
	if !isAlmanaxResourceAnswerValid(message) {
		panic(commands.ErrInvalidAnswerMessage)
	}

	// TODO 1 PERSONNAGE ATM
	webhookEdit := mappers.MapAlmanaxResourceToWebhook(message.GetEncyclopediaAlmanaxResourceAnswer(),
		1, constants.MapAMQPLocale(message.Language), command.emojiService)

	_, err := s.InteractionResponseEdit(i.Interaction, webhookEdit)
	if err != nil {
		log.Warn().Err(err).
			Msgf("Cannot respond to interaction after receiving internal answer, ignoring request")
	}
}

func isAlmanaxResourceAnswerValid(message *amqp.RabbitMQMessage) bool {
	return message.Status == amqp.RabbitMQMessage_SUCCESS &&
		message.Type == amqp.RabbitMQMessage_ENCYCLOPEDIA_ALMANAX_RESOURCE_ANSWER &&
		message.GetEncyclopediaAlmanaxResourceAnswer() != nil
}
