package main

import (
	"context"
	"fmt"
	pb "github.com/arifimran5/order-processing-demo/cmd/inventory/inventory"
	"github.com/arifimran5/order-processing-demo/pkg/shared/rabbitmq"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	exchangeName = "ecommerce_exchange"
	queueName    = "inventory_queue"
	routingKey   = "payment.done.#"
)

func main() {
	go runInventoryService()

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

	// process inventory messages

	messages, err := rq.Consume(queueName, "inventory", true, false, false, false, nil)
	if err != nil {
		fmt.Println("failed to consume", err)
	}

	go func() {
		for message := range messages {
			fmt.Printf("Received a message: %s \n", string(message.Body))
			fmt.Println("updating inventory")
			time.Sleep(1 * time.Second)
			fmt.Println("updated inventory 🥬")
			fmt.Println()
			fmt.Println()
		}
	}()

	fmt.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}

type server struct {
	pb.UnimplementedInventoryServiceServer
}

func (s *server) CheckInventory(ctx context.Context, req *pb.InventoryRequest) (*pb.InventoryResponse, error) {
	response := &pb.InventoryResponse{}

	for _, item := range req.Items {
		// check you database

		isAvailable := false
		var availableQuantity int32 = item.Quantity + 100

		if item.Quantity <= availableQuantity {
			isAvailable = true
		}

		availability := &pb.ItemAvailability{
			ProductId:         item.ProductId,
			IsAvailable:       isAvailable,
			AvailableQuantity: availableQuantity,
		}
		response.Availabilities = append(response.Availabilities, availability)
	}
	return response, nil
}

func runInventoryService() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, &server{})
	log.Printf("grpc server listening on %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
