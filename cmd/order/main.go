package main

import (
	"fmt"
	"github.com/arifimran5/order-processing-demo/pkg/shared/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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

	// Before creating order - let's check the inventory
	productIDs := []string{"Nike Air Jordan", "Fidget Spinner XYZ", "Hand Gloves"}
	quantities := []int32{5, 3, 120}

	res, err := CheckInventory(productIDs, quantities)
	if err != nil {
		log.Fatalf("failed to check inventory %v", err)
	}
	for _, availability := range res.Availabilities {
		log.Printf("Product %s: Available: %v, Quantity: %d", availability.ProductId, availability.IsAvailable, availability.AvailableQuantity)
		if !availability.IsAvailable {
			log.Printf("Product %s: Unavailable", availability.ProductId)
			return
		}
	}

	fmt.Println("publishing orders....")
	for i := 0; i < len(productIDs); i++ {
		err = rq.Publish(exchangeName, "payment.init.order", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(fmt.Sprintf("orderid: %v", productIDs[i])),
		})
		if err != nil {
			fmt.Println("failed to publish order", err)
		}
	}
}
