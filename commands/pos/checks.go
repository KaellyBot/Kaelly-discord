package pos

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-discord/models"
	"github.com/kaellybot/kaelly-discord/utils/middlewares"
	"github.com/kaellybot/kaelly-discord/utils/validators"
	i18n "github.com/kaysoro/discordgo-i18n"
	"github.com/rs/zerolog/log"
)

func (command *PosCommand) checkDimension(s *discordgo.Session, i *discordgo.InteractionCreate, lg discordgo.Locale, next middlewares.NextFunc) {

	data := i.ApplicationCommandData()

	// Filled case, expecting [1, 1] dimension
	for _, option := range data.Options {
		if option.Name == dimensionOptionName {
			dimensions := command.dimensionService.FindDimensions(option.StringValue(), lg)
			response, checkSuccess := validators.ExpectOnlyOneElement("pos.dimension.check", option.StringValue(), dimensions, lg)
			if checkSuccess {
				next()
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &response,
				})
				if err != nil {
					log.Error().Err(err).Msg("Dimension check response ignored")
				}
			}

			return
		}
	}

	// Option not filled, ANY dimension is then retrieved
	next()
}

func (command *PosCommand) checkServer(s *discordgo.Session, i *discordgo.InteractionCreate, lg discordgo.Locale, next middlewares.NextFunc) {

	data := i.ApplicationCommandData()

	// Filled case, expecting [1, 1] server
	for _, option := range data.Options {
		if option.Name == serverOptionName {
			servers := command.serverService.FindServers(option.StringValue(), lg)
			response, checkSuccess := validators.ExpectOnlyOneElement("pos.server.check", option.StringValue(), servers, lg)
			if checkSuccess {
				next()
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &response,
				})
				if err != nil {
					log.Error().Err(err).Msg("Server check response ignored")
				}
			}

			return
		}
	}

	// Option not filled (refers to guild and/or channel)
	server := command.guildService.GetServer()
	if server == nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: i18n.Get(lg, "pos.server.check.required", i18n.Vars{"game": models.Game}),
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("Server check response ignored")
		}
	} else {
		next()
	}
}
