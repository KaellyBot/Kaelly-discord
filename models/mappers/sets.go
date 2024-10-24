package mappers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
	contract "github.com/kaellybot/kaelly-commands"
	"github.com/kaellybot/kaelly-discord/models/constants"
	"github.com/kaellybot/kaelly-discord/services/characteristics"
	"github.com/kaellybot/kaelly-discord/services/emojis"
	"github.com/kaellybot/kaelly-discord/utils/discord"
	i18n "github.com/kaysoro/discordgo-i18n"
	"github.com/rs/zerolog/log"
)

func MapSetListRequest(query string, lg discordgo.Locale) *amqp.RabbitMQMessage {
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_ENCYCLOPEDIA_LIST_REQUEST,
		Language: constants.MapDiscordLocale(lg),
		Game:     constants.GetGame().AMQPGame,
		EncyclopediaListRequest: &amqp.EncyclopediaListRequest{
			Query: query,
			Type:  amqp.EncyclopediaListRequest_SET,
		},
	}
}

func MapSetToDefaultWebhookEdit(answer *amqp.EncyclopediaItemAnswer,
	characService characteristics.Service, emojiService emojis.Service,
	locale amqp.Language) *discordgo.WebhookEdit {
	return MapSetToWebhookEdit(answer, len(answer.GetSet().GetEquipments()),
		characService, emojiService, locale)
}

func MapSetToWebhookEdit(answer *amqp.EncyclopediaItemAnswer, itemNumber int,
	characService characteristics.Service, emojiService emojis.Service,
	locale amqp.Language) *discordgo.WebhookEdit {
	set := answer.GetSet()
	lg := constants.MapAMQPLocale(locale)
	bonus := &amqp.EncyclopediaItemAnswer_Set_Bonus{ItemNumber: 0}
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
		Embeds:     mapSetToEmbeds(answer, bonus, characService, lg),
		Components: mapSetToComponents(answer, bonus, emojiService, lg),
	}
}

func mapSetToEmbeds(answer *amqp.EncyclopediaItemAnswer,
	bonus *amqp.EncyclopediaItemAnswer_Set_Bonus, service characteristics.Service,
	lg discordgo.Locale) *[]*discordgo.MessageEmbed {
	set := answer.GetSet()
	fields := discord.SliceFields(set.GetEquipments(), constants.MaxEquipmentPerField,
		func(i int, items []*amqp.EncyclopediaItemAnswer_Set_Equipment) *discordgo.MessageEmbedField {
			name := constants.InvisibleCharacter
			if i == 0 {
				name = i18n.Get(lg, "set.items.title")
			}

			return &discordgo.MessageEmbedField{
				Name: name,
				Value: i18n.Get(lg, "set.items.description", i18n.Vars{
					"items": mapSetItems(items, lg),
				}),
				Inline: true,
			}
		})

	if bonus != nil {
		i18nEffects := mapEffects(bonus.GetEffects(), service)
		bonusFields := discord.SliceFields(i18nEffects, constants.MaxCharacterPerField,
			func(i int, items []i18nCharacteristic) *discordgo.MessageEmbedField {
				name := constants.InvisibleCharacter
				if i == 0 {
					name = i18n.Get(lg, "set.effects.title", i18n.Vars{
						"itemNumber": bonus.GetItemNumber(),
					})
				}

				return &discordgo.MessageEmbedField{
					Name: name,
					Value: i18n.Get(lg, "set.effects.description", i18n.Vars{
						"effects": items,
					}),
					Inline: true,
				}
			})

		if len(bonusFields) > 0 {
			fields = append(fields, discord.GhostInlineField())
			fields = append(fields, bonusFields...)
		}
	}

	return &[]*discordgo.MessageEmbed{
		{
			Title:       set.Name,
			Description: i18n.Get(lg, "set.description", i18n.Vars{"level": set.Level}),
			Color:       constants.Color,
			URL:         i18n.Get(lg, "set.url", i18n.Vars{"id": set.Id}),
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: set.GetIcon()},
			Fields:      fields,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    answer.Source.Name,
				URL:     answer.Source.Url,
				IconURL: answer.Source.Icon,
			},
			Footer: discord.BuildDefaultFooter(lg),
		},
	}
}

func mapSetToComponents(answer *amqp.EncyclopediaItemAnswer,
	bonus *amqp.EncyclopediaItemAnswer_Set_Bonus, service emojis.Service,
	lg discordgo.Locale) *[]discordgo.MessageComponent {
	set := answer.GetSet()
	components := make([]discordgo.MessageComponent, 0)

	var maxItemNumber int
	for _, currentBonus := range set.Bonuses {
		itemNumber := int(currentBonus.ItemNumber)
		if itemNumber > maxItemNumber {
			maxItemNumber = itemNumber
		}
	}
	bonuses := make([]discordgo.SelectMenuOption, 0)
	for _, currentBonus := range set.Bonuses {
		emoji := service.GetSetBonusEmoji(int(currentBonus.ItemNumber), maxItemNumber)
		bonuses = append(bonuses, discordgo.SelectMenuOption{
			Label: i18n.Get(lg, "set.effects.option", i18n.Vars{
				"itemNumber": currentBonus.ItemNumber,
				"itemCount":  maxItemNumber,
			}),
			Value:   fmt.Sprintf("%v", currentBonus.ItemNumber),
			Default: currentBonus.ItemNumber == bonus.ItemNumber,
			Emoji:   emoji,
		})
	}

	if len(bonuses) > 0 {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    contract.CraftSetBonusCustomID(set.Id),
					MenuType:    discordgo.StringSelectMenu,
					Placeholder: i18n.Get(lg, "set.effects.placeholder"),
					Options:     bonuses,
				},
			},
		})
	}

	items := make([]discordgo.SelectMenuOption, 0)
	for _, item := range set.Equipments {
		items = append(items, discordgo.SelectMenuOption{
			Label: item.Name,
			Value: item.Id,
			Emoji: service.GetEquipmentEmoji(item.Type),
		})
	}

	var itemType amqp.ItemType
	if set.GetIsCosmetic() {
		itemType = amqp.ItemType_COSMETIC_TYPE
	} else {
		itemType = amqp.ItemType_EQUIPMENT_TYPE
	}

	components = append(components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    contract.CraftItemCustomID(itemType.String()),
				MenuType:    discordgo.StringSelectMenu,
				Placeholder: i18n.Get(lg, "set.items.placeholder"),
				Options:     items,
			},
		},
	})

	return &components
}

type i18nSetItem struct {
	Name  string
	URL   string
	Level int64
}

func mapSetItems(items []*amqp.EncyclopediaItemAnswer_Set_Equipment,
	lg discordgo.Locale) []i18nSetItem {
	result := make([]i18nSetItem, 0)
	for _, item := range items {
		result = append(result, i18nSetItem{
			Name:  item.GetName(),
			URL:   i18n.Get(lg, "item.url", i18n.Vars{"id": item.GetId()}),
			Level: item.GetLevel(),
		})
	}

	return result
}
