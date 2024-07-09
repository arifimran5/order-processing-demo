package main

import (
	"context"
	pb "github.com/arifimran5/order-processing-demo/cmd/inventory/inventory"
	"google.golang.org/grpc"
	"log"
	"time"
)

func CheckInventory(productIds []string, quantities []int32) (*pb.InventoryResponse, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewInventoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	request := &pb.InventoryRequest{}
	for i, productId := range productIds {
		item := &pb.InventoryItem{
			ProductId: productId,
			Quantity:  quantities[i],
		}
		request.Items = append(request.Items, item)
	}
	return c.CheckInventory(ctx, request)
}
