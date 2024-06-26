package rabbit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jis4nx/go-ecom/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func ConnectRabbitMQ(user, passw, host, vhost string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", user, passw, host, vhost))
	return conn, err
}

func NewRabbitClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}
	return RabbitClient{conn: conn, ch: ch}, nil
}

func NewConsumer() (RabbitClient, error) {
	config := config.GetConfig()
	conn, err := ConnectRabbitMQ(config.RQ.USER, config.RQ.PASSWORD, config.RQ.HOST, config.RQ.VHOST)
	if err != nil {
		return RabbitClient{}, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}
	return RabbitClient{conn: conn, ch: ch}, nil
}

// Close the channel
func (rc *RabbitClient) Close() {
	rc.ch.Close()
}

func (rc *RabbitClient) NewQueue(name string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(name, durable, autoDelete, false, false, nil)
	return err
}

func (rc *RabbitClient) CreateBinding(name, routingKey, exchangeKey string) error {
	return rc.ch.QueueBind(name, routingKey, exchangeKey, false, nil)
}

func (rc *RabbitClient) Consume(queueName, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queueName, consumer, autoAck, false, false, false, nil)
}

func (rc *RabbitClient) CreateExchange(name, exchangeType string, durable, autoDelete bool) error {
	return rc.ch.ExchangeDeclare(name, exchangeType, durable, autoDelete, false, false, nil)
}

func (rc *RabbitClient) PublishMsgWithContext(ctx context.Context, exchangeKey, routingKey string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return rc.ch.PublishWithContext(
		ctx,
		exchangeKey,
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Transient,
			Body:         data,
		},
	)
}
