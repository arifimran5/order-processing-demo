package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type ExchangeConfig struct {
	Name string
	// Type is the exchange type
	Type string
	// Durable indicates that the exchange will survive a server restart
	Durable bool
	// AutoDelete indicates that the exchange will be deleted once the last queue is unbound
	AutoDelete bool
	// Internal indicates that the exchange is only for use by the server
	Internal bool
	// NoWait indicates that the server should not respond to the method
	NoWait bool
	Args   amqp.Table
}

type QueueConfig struct {
	Name string
	// Durable indicates that the queue will survive a server restart
	Durable bool
	// AutoDelete indicates that the queue will be deleted once the last consumer unsubscribes
	AutoDelete bool
	// Exclusive indicates that only one consumer can access the queue
	Exclusive bool
	// NoWait indicates that the server should not respond to the method
	NoWait bool
	Args   amqp.Table
}

// New creates a new rabbitmq connection
func New(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close closes the channel and connection and returns the error if any
func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

// DeclareExchange declares an exchange
func (r *RabbitMQ) DeclareExchange(config ExchangeConfig) error {
	return r.channel.ExchangeDeclare(
		config.Name,
		config.Type,
		config.Durable,
		config.AutoDelete,
		config.Internal,
		config.NoWait,
		config.Args,
	)
}

// DeclareQueue declares a queue
func (r *RabbitMQ) DeclareQueue(config QueueConfig) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		config.Args,
	)
}

// BindQueue binds a Queue to an exchange
func (r *RabbitMQ) BindQueue(name, key, exchange string, args amqp.Table) error {
	return r.channel.QueueBind(name, key, exchange, false, args)
}

// Publish publishes the message to an exchange
func (r *RabbitMQ) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

// Consume starts consuming messages from a queue
func (r *RabbitMQ) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)
}

// SetupDeadLetterQueue sets up a dead-letter queue and binds it to an exchange
func (r *RabbitMQ) SetupDeadLetterQueue(queueName, exchangeName, routingKey string) error {
	err := r.DeclareExchange(ExchangeConfig{
		Name:    exchangeName,
		Type:    "direct",
		Durable: true,
	})
	if err != nil {
		return fmt.Errorf("failed to declare exchange for a dead letter queue: %w", err)
	}

	_, err = r.DeclareQueue(QueueConfig{
		Name:    queueName,
		Durable: true,
		Args:    nil,
	})
	if err != nil {
		return fmt.Errorf("failed to declare dead-letter queue: %w", err)
	}
	err = r.BindQueue(queueName, routingKey, exchangeName, nil)
	if err != nil {
		return fmt.Errorf("failed to bind dead-letter queue: %w", err)
	}

	return nil
}
