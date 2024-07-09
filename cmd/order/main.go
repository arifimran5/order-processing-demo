package main

import (
	"fmt"
	"github.com/arifimran5/order-processing-demo/pkg/shared/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"math/rand"
)

var (
	exchangeName = "ecommerce_exchange"
)

func main() {
	rq, err := rabbitmq.New("amqp://admin:admin@localhost:5672")
	if err != nil {
		fmt.Println("failed to connect to rabbitmq", err)
		panic(err)
	}
	defer rq.Close()

	err = rq.DeclareExchange(rabbitmq.ExchangeConfig{
		Name:       exchangeName,
		Durable:    true,
		Type:       "topic",
		AutoDelete: false,
		NoWait:     false,
	})

	if err != nil {
		fmt.Println("failed to declare exchange", err)
		panic(err)
	}

	fmt.Println("publishing orders....")
	for i := 0; i < 5; i++ {
		randomId := rand.Intn(1000) + 1
		err = rq.Publish(exchangeName, "payment.init.order", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("orderid: %d", randomId)),
		})
		if err != nil {
			fmt.Println("failed to publish order", err)
		}
	}
}
