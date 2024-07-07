package main

import (
	"fmt"
	"github.com/arifimran5/order-processing-demo/pkg/shared/rabbitmq"
	"time"
)

var (
	exchangeName = "ecommerce_exchange"
	queueName    = "email_queue"
	routingKey   = "email.#"
)

func main() {
	rq, err := rabbitmq.New("amqp://admin:admin@localhost:5672")
	if err != nil {
		fmt.Println("failed to connect to rabbitmq", err)
		panic(err)
	}
	defer rq.Close()

	_, err = rq.DeclareQueue(rabbitmq.QueueConfig{
		Name:       queueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	})
	if err != nil {
		fmt.Println("failed to declare queue", err)
		panic(err)
	}

	if err = rq.BindQueue(queueName, routingKey, exchangeName, nil); err != nil {
		fmt.Println("failed to bind queue", err)
	}

	forever := make(chan struct{})

	// process email messages

	messages, err := rq.Consume(queueName, "email", true, false, false, false, nil)
	if err != nil {
		fmt.Println("failed to consume", err)
	}

	go func() {
		for message := range messages {
			fmt.Printf("Received a message: %s \n", string(message.Body))
			fmt.Println("sending email")
			time.Sleep(1 * time.Second)
			fmt.Println("sent email 📧")

			fmt.Println()
			fmt.Println()
		}
	}()

	fmt.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}
