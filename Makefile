@PHONY: build run-order run-payment run-inventory run-email

build:
	go build -o bin/order ./cmd/order
	go build -o bin/payment ./cmd/payment
	go build -o bin/inventory ./cmd/inventory
	go build -o bin/email ./cmd/email

run-order:
	go run ./cmd/order

run-payment:
	go run ./cmd/payment

run-inventory:
	go run ./cmd/inventory

run-email:
	go run ./cmd/email

grpc-inventory:
	protoc --go_out=cmd/inventory --go-grpc_out=cmd/inventory cmd/inventory/inventory.proto