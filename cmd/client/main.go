package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	newsv1 "github.com/sanjevscet/news-grpc/api/news/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client: %v\n", err)
	}

	newClient := newsv1.NewNewsServiceClient(conn)

	ctx := context.Background()

	res, err := newClient.Create(ctx, &newsv1.CreateRequest{
		Id:      uuid.NewString(),
		Author:  "John Doe",
		Title:   "New Article",
		Summary: "This is a summary of the new article.",
		Content: "This is the content of the new article.",
		Tags:    []string{"tag1", "tag2"},
	})

	if err != nil {
		log.Fatalf("failed to create article: %v\n", err)
	}

	log.Printf("Article created: %+v\n", res)

	getResponse, err := newClient.Get(ctx, &newsv1.GetRequest{
		Id: res.Id,
	})

	if err != nil {
		log.Fatalf("failed to get article: %v\n", err)
	}

	log.Printf("Article fetched: %+v\n", getResponse)

	defer conn.Close()
}
