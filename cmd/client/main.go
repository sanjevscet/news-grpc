package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/google/uuid"
	newsv1 "github.com/sanjevscet/news-grpc/api/news/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create client: %v\n", err)
	}

	newClient := newsv1.NewNewsServiceClient(conn)

	ctx := context.Background()

	println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	println("Creating a new article...")
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
	println("Article created successfully!")
	log.Printf("Article created: %+v\n", res)

	println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	println("Getting the created article...")
	getResponse, err := newClient.Get(ctx, &newsv1.GetRequest{
		Id: res.Id,
	})

	if err != nil {
		log.Fatalf("failed to get article: %v\n", err)
	}

	log.Printf("Article fetched: %+v\n", getResponse)
	println("Article fetched successfully!")

	println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	println("Creating Multiple articles...")
	for i := range 5 {
		res, err := newClient.Create(ctx, &newsv1.CreateRequest{
			Id:      uuid.NewString(),
			Author:  "John Doe " + strconv.Itoa(i),
			Title:   "New Article " + strconv.Itoa(i),
			Summary: fmt.Sprintf("This is a summary of the new article %d", i),
			Content: "This is the content of the new article.",
			Tags:    []string{"tag1", "tag2"},
		})

		if err != nil {
			log.Fatalf("failed to create article: %v\n", err)
		}
		log.Printf("Article created with ID: %s\n", res.Id)
	}
	println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	println("Getting all articles...")

	stream, err := newClient.GetAll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("failed to get articles: %v\n", err)
	}

	for {
		article, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive article: %v\n", err)
		}
		log.Printf("Article fetched: %s, summary: %s\n", article.Id, article.Summary)
	}
	println("All articles fetched successfully!")
	println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	println("Client finished successfully!")

	defer conn.Close()
}
