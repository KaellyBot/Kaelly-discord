package set

import (
	"github.com/kaellybot/kaelly-discord/commands"
	"github.com/kaellybot/kaelly-discord/services/characteristics"
	"github.com/kaellybot/kaelly-discord/utils/requests"
)

const (
	setRequestRoutingKey = "requests.encyclopedias"

	defaultIconSize = "128"
)

type Command struct {
	commands.AbstractCommand
	characteristicService characteristics.Service
	requestManager        requests.RequestManager
	handlers              commands.DiscordHandlers
}