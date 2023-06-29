package mappers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-discord/models/constants"
	"github.com/kaellybot/kaelly-discord/services/characteristics"
	"github.com/kaellybot/kaelly-discord/utils/discord"
	i18n "github.com/kaysoro/discordgo-i18n"
	"github.com/rs/zerolog/log"
)

func MapSetListRequest(query string, lg discordgo.Locale) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_LIST_REQUEST,
		Language: constants.MapDiscordLocale(lg),
		EncyclopediaSetListRequest: &amqp.EncyclopediaSetListRequest{
			Query: query,
		},
	}
}

func MapSetRequest(query string, lg discordgo.Locale) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_SET_REQUEST,
		Language: constants.MapDiscordLocale(lg),
		EncyclopediaSetRequest: &amqp.EncyclopediaSetRequest{
			Query: query,
		},
	}
}

func MapSetToDefaultWebhookEdit(set *amqp.EncyclopediaSetAnswer, service characteristics.Service,
	locale amqp.Language) *discordgo.WebhookEdit {
	return MapSetToWebhookEdit(set, len(set.Equipments), service, locale)
}

func MapSetToWebhookEdit(set *amqp.EncyclopediaSetAnswer, itemNumber int,
	service characteristics.Service, locale amqp.Language) *discordgo.WebhookEdit {
	lg := constants.MapAMQPLocale(locale)
	bonus := &amqp.EncyclopediaSetAnswer_Bonus{ItemNumber: 0}
	for _, currentBonus := range set.Bonuses {
		if currentBonus.ItemNumber == int64(itemNumber) {
			bonus = currentBonus
			break
		} else if bonus.ItemNumber < currentBonus.ItemNumber {
			bonus = currentBonus
		}
	}

	if bonus != nil && bonus.ItemNumber != int64(itemNumber) {
		log.Warn().
			Str(constants.LogAnkamaID, set.Id).
			Int(constants.LogItemNumber, itemNumber).
			Msgf("Set bonus with specific item numbers was not found, returning the highest one...")
	}

	return &discordgo.WebhookEdit{
		Embeds:     mapSetToEmbeds(set, bonus, service, lg),
		Components: mapSetToComponents(set, bonus, lg),
	}
}

func mapSetToEmbeds(set *amqp.EncyclopediaSetAnswer, bonus *amqp.EncyclopediaSetAnswer_Bonus,
	service characteristics.Service, lg discordgo.Locale) *[]*discordgo.MessageEmbed {
	fields := discord.SliceFields(set.GetEquipments(), constants.MaxEquipmentPerField,
		func(i int, items []*amqp.EncyclopediaSetAnswer_Equipment) *discordgo.MessageEmbedField {
			name := constants.InvisibleCharacter
			if i == 0 {
				name = i18n.Get(lg, "set.items.title")
			}

			return &discordgo.MessageEmbedField{
				Name: name,
				Value: i18n.Get(lg, "set.items.description", i18n.Vars{
					"items": items,
				}),
				Inline: true,
			}
		})

	if bonus != nil {
		bonusFields := discord.SliceFields(bonus.GetEffects(), constants.MaxCharacterPerField,
			func(i int, items []*amqp.EncyclopediaSetAnswer_Effect) *discordgo.MessageEmbedField {
				name := constants.InvisibleCharacter
				if i == 0 {
					name = i18n.Get(lg, "set.effects.title", i18n.Vars{
						"itemNumber": bonus.GetItemNumber(),
					})
				}

				return &discordgo.MessageEmbedField{
					Name: name,
					Value: i18n.Get(lg, "set.effects.description", i18n.Vars{
						"effects": mapSetEffects(items, service),
					}),
					Inline: true,
				}
			})

		fields = append(fields, bonusFields...)
	}

	return &[]*discordgo.MessageEmbed{
		{
			Title:       set.Name,
			Description: i18n.Get(lg, "set.description", i18n.Vars{"level": set.Level, "items": set.Equipments}),
			Color:       constants.Color,
			URL:         i18n.Get(lg, "set.url", i18n.Vars{"id": set.Id}),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://static.ankama.com/dofus/renderer/look/7b317c38302c323132342c323531327c313d31363736353536342c323d31363335353838332c333d31363737373138352c343d323931303036342c353d31343536313739397c3134307d/full/1/200_200-10.png",
			}, // TODO URL
			Fields: fields,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    set.Source.Name,
				URL:     set.Source.Url,
				IconURL: set.Source.Icon,
			},
		},
	}
}

func mapSetToComponents(set *amqp.EncyclopediaSetAnswer, bonus *amqp.EncyclopediaSetAnswer_Bonus,
	lg discordgo.Locale) *[]discordgo.MessageComponent {
	components := make([]discordgo.MessageComponent, 0)

	bonuses := make([]discordgo.SelectMenuOption, 0)
	for _, currentBonus := range set.Bonuses {
		bonuses = append(bonuses, discordgo.SelectMenuOption{
			Label: i18n.Get(lg, "set.effects.option", i18n.Vars{
				"itemNumber": currentBonus.ItemNumber,
				"itemCount":  len(set.Equipments),
			}),
			Value:   fmt.Sprintf("%v", currentBonus.ItemNumber),
			Default: currentBonus.ItemNumber == bonus.ItemNumber,
			Emoji: discordgo.ComponentEmoji{
				Name: constants.EmojiEffectID,
			},
		})
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    fmt.Sprintf("/sets/%v/effects", set.Id),
				MenuType:    discordgo.StringSelectMenu,
				Placeholder: i18n.Get(lg, "set.effects.placeholder"),
				Options:     bonuses,
			},
		},
	})

	items := make([]discordgo.SelectMenuOption, 0)
	for _, item := range set.Equipments {
		items = append(items, discordgo.SelectMenuOption{
			Label: item.Name,
			Value: item.Id,
			Emoji: discordgo.ComponentEmoji{
				ID: constants.EmojiHatID,
			},
		})
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    fmt.Sprintf("/sets/%v/items", set.Id),
				MenuType:    discordgo.StringSelectMenu,
				Placeholder: i18n.Get(lg, "set.items.placeholder"),
				Options:     items,
			},
		},
	})

	return &components
}
