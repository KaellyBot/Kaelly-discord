package requests

import (
	"context"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	AnswersRoutingKey = "answers.*"
	AnswersQueueName  = "answers"
)

type RequestManager interface {
	Request(s *discordgo.Session, i *discordgo.InteractionCreate, routingKey string,
		message *amqp.RabbitMQMessage, callback RequestCallback, properties ...map[string]any) error
	Listen() error
}

type RequestManagerImpl struct {
	broker   amqp.MessageBroker
	requests map[string]discordRequest
}

type discordRequest struct {
	session     *discordgo.Session
	interaction *discordgo.InteractionCreate
	callback    RequestCallback
	properties  map[string]any
}

type RequestCallback func(ctx context.Context, s *discordgo.Session,
	i *discordgo.InteractionCreate, message *amqp.RabbitMQMessage, properties map[string]any)
