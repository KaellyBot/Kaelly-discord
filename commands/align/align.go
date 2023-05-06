package align

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	contract "github.com/kaellybot/kaelly-commands"
	"github.com/kaellybot/kaelly-discord/commands"
	"github.com/kaellybot/kaelly-discord/models/constants"
	"github.com/kaellybot/kaelly-discord/models/entities"
	"github.com/kaellybot/kaelly-discord/services/books"
	"github.com/kaellybot/kaelly-discord/services/guilds"
	"github.com/kaellybot/kaelly-discord/services/servers"
	"github.com/kaellybot/kaelly-discord/utils/middlewares"
	"github.com/kaellybot/kaelly-discord/utils/requests"
)

//nolint:exhaustive // only useful handlers must be implemented, it will panic also
func New(bookService books.Service, guildService guilds.Service,
	serverService servers.Service, requestManager requests.RequestManager) *Command {
	cmd := Command{
		bookService:    bookService,
		guildService:   guildService,
		serverService:  serverService,
		requestManager: requestManager,
	}

	subCommandHandlers := cmd.HandleSubCommand(commands.SubCommandHandlers{
		contract.AlignGetSubCommandName: middlewares.
			Use(cmd.checkOptionalCity, cmd.checkOptionalOrder, cmd.checkServer, cmd.getRequest),
		contract.AlignSetSubCommandName: middlewares.
			Use(cmd.checkMandatoryCity, cmd.checkMandatoryOrder, cmd.checkLevel, cmd.checkServer, cmd.setRequest),
	})

	cmd.slashHandlers = commands.DiscordHandlers{
		discordgo.InteractionApplicationCommand:             subCommandHandlers,
		discordgo.InteractionApplicationCommandAutocomplete: cmd.autocomplete,
	}
	cmd.userHandlers = commands.DiscordHandlers{
		discordgo.InteractionApplicationCommand: middlewares.Use(cmd.checkServer, cmd.userRequest),
	}
	return &cmd
}

func (command *Command) Matches(i *discordgo.InteractionCreate) bool {
	return len(i.ApplicationCommandData().TargetID) == 0 &&
		contract.AlignSlashCommandName == i.ApplicationCommandData().Name ||
		contract.AlignUserCommandName == i.ApplicationCommandData().Name
}

func (command *Command) Handle(s *discordgo.Session, i *discordgo.InteractionCreate, lg discordgo.Locale) {
	if len(i.ApplicationCommandData().TargetID) == 0 {
		command.CallHandler(s, i, lg, command.slashHandlers)
	} else {
		command.CallHandler(s, i, lg, command.userHandlers)
	}
}

func (command *Command) getGetOptions(ctx context.Context) (entities.City, entities.Order, entities.Server, error) {
	city, ok := ctx.Value(constants.ContextKeyCity).(entities.City)
	if !ok {
		city = entities.City{}
	}

	order, ok := ctx.Value(constants.ContextKeyOrder).(entities.Order)
	if !ok {
		order = entities.Order{}
	}

	server, ok := ctx.Value(constants.ContextKeyServer).(entities.Server)
	if !ok {
		return entities.City{}, entities.Order{}, entities.Server{},
			fmt.Errorf("cannot cast %v as entities.Server", ctx.Value(constants.ContextKeyServer))
	}

	return city, order, server, nil
}

func (command *Command) getSetOptions(ctx context.Context) (entities.City, entities.Order,
	int64, entities.Server, error) {
	city, ok := ctx.Value(constants.ContextKeyCity).(entities.City)
	if !ok {
		return entities.City{}, entities.Order{}, 0, entities.Server{},
			fmt.Errorf("cannot cast %v as entities.City", ctx.Value(constants.ContextKeyCity))
	}

	order, ok := ctx.Value(constants.ContextKeyOrder).(entities.Order)
	if !ok {
		return entities.City{}, entities.Order{}, 0, entities.Server{},
			fmt.Errorf("cannot cast %v as entities.Order", ctx.Value(constants.ContextKeyOrder))
	}

	server, ok := ctx.Value(constants.ContextKeyServer).(entities.Server)
	if !ok {
		return entities.City{}, entities.Order{}, 0, entities.Server{},
			fmt.Errorf("cannot cast %v as entities.Server", ctx.Value(constants.ContextKeyServer))
	}

	level, ok := ctx.Value(constants.ContextKeyLevel).(int64)
	if !ok {
		return entities.City{}, entities.Order{}, 0, entities.Server{},
			fmt.Errorf("cannot cast %v as uint", ctx.Value(constants.ContextKeyLevel))
	}

	return city, order, level, server, nil
}

func (command *Command) getUserOptions(ctx context.Context) (entities.Server, error) {
	server, ok := ctx.Value(constants.ContextKeyServer).(entities.Server)
	if !ok {
		return entities.Server{}, fmt.Errorf("cannot cast %v as entities.Server", ctx.Value(constants.ContextKeyServer))
	}

	return server, nil
}