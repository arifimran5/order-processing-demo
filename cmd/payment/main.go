package main

import (
	"fmt"
	"github.com/arifimran5/order-processing-demo/pkg/shared/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var (
	exchangeName = "ecommerce_exchange"
	queueName    = "payment_queue"
	routingKey   = "payment.init.#"
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

	messages, err := rq.Consume(queueName, "payment", true, false, false, false, nil)
	if err != nil {
		fmt.Println("failed to register consumer", err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range messages {
			fmt.Printf("received message: %v \n", string(d.Body))
			fmt.Println("processing payment...")
			time.Sleep(time.Second * 2)
			fmt.Println("payment processed")

			fmt.Println("initiating inventory management and email service ...")

			// send message to inventory queue
			err = rq.Publish(exchangeName, "payment.done.inventory", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf(string(d.Body))),
			})
			if err != nil {
				fmt.Println("Failed to publish inventory update", err)
			}

			//send message to email queue
			//err = rq.Publish(exchangeName, "payment.done.email", false, false, amqp.Publishing{
			//	ContentType: "text/plain",
			//	Body:        []byte(fmt.Sprintf("send email: %v", string(d.Body))),
			//})
			//if err != nil {
			//	fmt.Println("Failed to publish email", err)
			//}
			//fmt.Println("INITIATED INVENTORY & EMAIL 💥")
			//fmt.Println()
			//fmt.Println()
		}
	}()
	fmt.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}
