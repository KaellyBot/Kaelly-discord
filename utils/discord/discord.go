package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kaellybot/kaelly-discord/models/constants"
	"github.com/kaellybot/kaelly-discord/services/emojis"
	"github.com/kaellybot/kaelly-discord/utils/slicers"
	i18n "github.com/kaysoro/discordgo-i18n"
)

func GetPaginationButtons(page, pages int, crafter CraftPageCustomID,
	lg discordgo.Locale, emojiService emojis.Service) *[]discordgo.MessageComponent {
	previousPage := page - 1
	if previousPage < constants.DefaultPage {
		previousPage = constants.DefaultPage
	}

	lastPage := pages - 1
	nextPage := page + 1
	if nextPage > lastPage {
		nextPage = lastPage
	}

	return &[]discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: crafter(constants.DefaultPage),
					Label:    i18n.Get(lg, "default.page.first"),
					Style:    discordgo.PrimaryButton,
					Disabled: page <= constants.DefaultPage,
					Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDFirst),
				},
				discordgo.Button{
					CustomID: crafter(previousPage),
					Label:    i18n.Get(lg, "default.page.previous"),
					Style:    discordgo.PrimaryButton,
					Disabled: page <= constants.DefaultPage,
					Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDPrevious),
				},
				discordgo.Button{
					CustomID: crafter(nextPage),
					Label:    i18n.Get(lg, "default.page.next"),
					Style:    discordgo.PrimaryButton,
					Disabled: page >= lastPage,
					Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDNext),
				},
				discordgo.Button{
					CustomID: crafter(lastPage),
					Label:    i18n.Get(lg, "default.page.last"),
					Style:    discordgo.PrimaryButton,
					Disabled: page >= lastPage,
					Emoji:    emojiService.GetMiscEmoji(constants.EmojiIDLast),
				},
			},
		},
	}
}

func SliceFields[T any](items []T, limit int, toField ItemsToField[T]) []*discordgo.MessageEmbedField {
	fields := make([]*discordgo.MessageEmbedField, 0)
	slicedItems := slicers.Slice(items, limit)
	for i, items := range slicedItems {
		fields = append(fields, toField(i, items))
	}

	return fields
}

func SliceButtons[T any](items []T, toButton ItemToButton[T]) []discordgo.ActionsRow {
	actionsRow := make([]discordgo.ActionsRow, 0)
	slicedItems := slicers.Slice(items, constants.MaxButtonPerActionRow)
	for _, items := range slicedItems {
		buttons := make([]discordgo.MessageComponent, 0)
		for _, subItem := range items {
			buttons = append(buttons, toButton(subItem))
		}
		actionsRow = append(actionsRow, discordgo.ActionsRow{
			Components: buttons,
		})
	}

	return actionsRow
}

func GhostInlineField() *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   constants.InvisibleCharacter,
		Inline: false,
	}
}

func BuildDefaultFooter(lg discordgo.Locale) *discordgo.MessageEmbedFooter {
	return &discordgo.MessageEmbedFooter{
		Text: i18n.Get(lg, "default.footer", i18n.Vars{
			"name":      constants.Name,
			"version":   constants.Version,
			"changelog": i18n.Get(lg, "default.changelog"),
		}),
		IconURL: constants.AvatarIcon,
	}
}
