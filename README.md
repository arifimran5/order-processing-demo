## Ecommerce Order Processing Demo with Rabbitmq and Go

> [!NOTE]  
> Services are for demo purpose, it's not feature complete. Goal is to show how to use rabbitmq with go and create a microservices architecture.

This repository contains 4 services that process an ecommerce order and a shared pkg for rabbitmq:
1. Order
2. Payment
3. Inventory
4. Email

### A shared package for rabbitmq utility functions
- Creating connection
- Declaring exchanges
- Declaring queues
- Binding queues to exchanges
- Publishing messages to exchanges
- Consuming messages from queues
- Setting up dead-letter queues


### Prerequisites
1. Rabbitmq server
2. Go
3. [rabbitmq/amqp go client](https://github.com/rabbitmq/amqp091-go)

It's a cli demo. All the required commands to run the services are provided in the Makefile.
Or run services individually by calling `go run cmd/order/main.go` etc.


### Place Order by running - Order service
```bash
make run-order
```

### To consume messages run following services
```bash
make run-payment
make run-inventory
make run-email
```