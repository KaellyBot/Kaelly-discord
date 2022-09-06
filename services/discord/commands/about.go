package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-discord/models"
	i18n "github.com/kaysoro/discordgo-i18n"
	"github.com/rs/zerolog/log"
)

const (
	CommandNameAbout = "about"
)

func About(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{getAboutEmbed(i.Locale)},
		},
	})
	if err != nil {
		log.Error().Err(err).Msgf("Cannot handle about reponse")
	}
}

func getAboutEmbed(locale discordgo.Locale) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       i18n.Get(locale, "about.title", i18n.Vars{"name": models.Name, "version": models.Version}),
		Description: i18n.Get(locale, "about.desc", i18n.Vars{"game": models.Game}),
		Color:       models.Color,
		Image:       &discordgo.MessageEmbedImage{URL: models.AvatarImage},
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: models.AvatarIcon},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   i18n.Get(locale, "about.invite.title"),
				Value:  i18n.Get(locale, "about.invite.desc", i18n.Vars{"invite": models.Invite}),
				Inline: false,
			},
			{
				Name:   i18n.Get(locale, "about.support.title"),
				Value:  i18n.Get(locale, "about.support.desc", i18n.Vars{"discord": models.Discord}),
				Inline: false,
			},
			{
				Name:   i18n.Get(locale, "about.twitter.title"),
				Value:  i18n.Get(locale, "about.twitter.desc", i18n.Vars{"twitter": models.Twitter}),
				Inline: false,
			},
			{
				Name:   i18n.Get(locale, "about.opensource.title"),
				Value:  i18n.Get(locale, "about.opensource.desc", i18n.Vars{"github": models.Github}),
				Inline: false,
			},
			{
				Name:   i18n.Get(locale, "about.free.title"),
				Value:  i18n.Get(locale, "about.free.desc", i18n.Vars{"paypal": models.Paypal}),
				Inline: false,
			},
			{
				Name:   i18n.Get(locale, "about.graphist.title"),
				Value:  i18n.Get(locale, "about.graphist.desc", i18n.Vars{"graphist": models.Elycann}),
				Inline: false,
			},
			{
				Name:   i18n.Get(locale, "about.donors.title"),
				Value:  i18n.Get(locale, "about.donors.desc", i18n.Vars{"donors": models.Donors}),
				Inline: false,
			},
		},
	}
}
