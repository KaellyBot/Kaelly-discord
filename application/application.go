package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-discord/commands"
	"github.com/kaellybot/kaelly-discord/commands/about"
	"github.com/kaellybot/kaelly-discord/commands/align"
	"github.com/kaellybot/kaelly-discord/commands/almanax"
	"github.com/kaellybot/kaelly-discord/commands/config"
	"github.com/kaellybot/kaelly-discord/commands/help"
	"github.com/kaellybot/kaelly-discord/commands/item"
	"github.com/kaellybot/kaelly-discord/commands/job"
	"github.com/kaellybot/kaelly-discord/commands/pos"
	"github.com/kaellybot/kaelly-discord/commands/set"
	"github.com/kaellybot/kaelly-discord/models/constants"
	"github.com/kaellybot/kaelly-discord/repositories/areas"
	characRepo "github.com/kaellybot/kaelly-discord/repositories/characteristics"
	"github.com/kaellybot/kaelly-discord/repositories/cities"
	"github.com/kaellybot/kaelly-discord/repositories/dimensions"
	emojiRepo "github.com/kaellybot/kaelly-discord/repositories/emojis"
	feedRepo "github.com/kaellybot/kaelly-discord/repositories/feeds"
	guildRepo "github.com/kaellybot/kaelly-discord/repositories/guilds"
	"github.com/kaellybot/kaelly-discord/repositories/jobs"
	"github.com/kaellybot/kaelly-discord/repositories/orders"
	serverRepo "github.com/kaellybot/kaelly-discord/repositories/servers"
	streamerRepo "github.com/kaellybot/kaelly-discord/repositories/streamers"
	"github.com/kaellybot/kaelly-discord/repositories/subareas"
	"github.com/kaellybot/kaelly-discord/repositories/transports"
	videastRepo "github.com/kaellybot/kaelly-discord/repositories/videasts"
	"github.com/kaellybot/kaelly-discord/services/books"
	"github.com/kaellybot/kaelly-discord/services/characteristics"
	"github.com/kaellybot/kaelly-discord/services/discord"
	"github.com/kaellybot/kaelly-discord/services/emojis"
	"github.com/kaellybot/kaelly-discord/services/feeds"
	"github.com/kaellybot/kaelly-discord/services/guilds"
	"github.com/kaellybot/kaelly-discord/services/portals"
	"github.com/kaellybot/kaelly-discord/services/servers"
	"github.com/kaellybot/kaelly-discord/services/streamers"
	"github.com/kaellybot/kaelly-discord/services/videasts"
	"github.com/kaellybot/kaelly-discord/utils/databases"
	"github.com/kaellybot/kaelly-discord/utils/requests"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	db, err := databases.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("DB instantiation failed, shutting down.")
	}

	broker, err := amqp.New(constants.GetRabbitMQClientID(), viper.GetString(constants.RabbitMQAddress), getBindings())
	if err != nil {
		log.Fatal().Err(err).Msgf("Broker instantiation failed, shutting down.")
	}

	dimensionRepo := dimensions.New(db)
	areaRepo := areas.New(db)
	subAreaRepo := subareas.New(db)
	transportTypeRepo := transports.New(db)
	portalService, err := portals.New(dimensionRepo, areaRepo, subAreaRepo, transportTypeRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Dimension Service instantiation failed, shutting down.")
	}

	serverRepo := serverRepo.New(db)
	serverService, err := servers.New(serverRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Server Service instantiation failed, shutting down.")
	}

	jobRepo := jobs.New(db)
	cityRepo := cities.New(db)
	orderRepo := orders.New(db)
	bookService, err := books.New(jobRepo, cityRepo, orderRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Book Service instantiation failed, shutting down.")
	}

	feedRepo := feedRepo.New(db)
	feedService, err := feeds.New(feedRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Feed Service instantiation failed, shutting down.")
	}

	videastRepo := videastRepo.New(db)
	videastService, err := videasts.New(videastRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Vidast Service instantiation failed, shutting down.")
	}

	streamerRepo := streamerRepo.New(db)
	streamerService, err := streamers.New(streamerRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Streamer Service instantiation failed, shutting down.")
	}

	characRepo := characRepo.New(db)
	characService, err := characteristics.New(characRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Characteristic Service instantiation failed, shutting down.")
	}

	emojiRepo := emojiRepo.New(db)
	emojiService, err := emojis.New(emojiRepo)
	if err != nil {
		log.Fatal().Err(err).Msgf("Emoji Service instantiation failed, shutting down.")
	}

	guildRepo := guildRepo.New(db)
	guildService := guilds.New(guildRepo)
	requestsManager := requests.New(broker)

	commands := make([]commands.DiscordCommand, 0)
	commands = append(commands,
		about.New(),
		align.New(bookService, guildService, serverService, requestsManager),
		almanax.New(emojiService, requestsManager),
		config.New(guildService, feedService, serverService, streamerService, videastService, requestsManager),
		help.New(&commands),
		item.New(characService, emojiService, requestsManager),
		job.New(bookService, guildService, serverService, requestsManager),
		pos.New(guildService, portalService, serverService, requestsManager),
		set.New(characService, emojiService, requestsManager),
	)

	discordService, err := discord.New(
		viper.GetString(constants.Token),
		viper.GetInt(constants.ShardID),
		viper.GetInt(constants.ShardCount), commands)
	if err != nil {
		return nil, err
	}

	return &Impl{
		db:             db,
		broker:         broker,
		requestManager: requestsManager,
		bookService:    bookService,
		guildService:   guildService,
		portalService:  portalService,
		serverService:  serverService,
		discordService: discordService,
	}, nil
}

func (app *Impl) Run() error {
	err := app.requestManager.Listen()
	if err != nil {
		return err
	}

	err = app.discordService.Listen()
	if err != nil {
		return err
	}

	return nil
}

func (app *Impl) Shutdown() {
	app.db.Shutdown()
	app.broker.Shutdown()
	app.discordService.Shutdown()
	log.Info().Msgf("Application is no longer running")
}

func getBindings() []amqp.Binding {
	return []amqp.Binding{
		requests.GetBinding(),
	}
}
