package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Chetas1/grpc-blog-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func runTest(ctx context.Context, client proto.BlogServiceClient) {
	fmt.Println("--- Starting gRPC CRUD Test ---")
	// Create Post
	createReq := &proto.CreatePostRequest{
		Title:           "Cloud Computing",
		Content:         "Learn fundamentals of cloud computing",
		Author:          "Chetas",
		PublicationDate: timestamppb.Now(),
		Tags:            []string{"cloud", "azure"},
	}

	createdPost, err := client.CreatePost(ctx, createReq)
	if err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	fmt.Printf("[CREATE] Created Post ID: %s\n", createdPost.PostId)

	// Read Post
	readPost, err := client.ReadPost(ctx, &proto.ReadPostRequest{PostId: createdPost.PostId})
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	fmt.Printf("[READ] Retrieved Post: %s by %s\n", readPost.Title, readPost.Author)

	// Update Post
	updateReq := &proto.UpdatePostRequest{
		PostId:  createdPost.PostId,
		Title:   "Cloud Computing(updated)",
		Content: "Updated content for the blog post.",
		Author:  "Chetas",
		Tags:    []string{"az", "csp", "updated"},
	}
	updatedPost, err := client.UpdatePost(ctx, updateReq)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	fmt.Printf("[UPDATE] Updated Title: %s\n", updatedPost.Title)

	// Delete Post
	deleteRes, err := client.DeletePost(ctx, &proto.DeletePostRequest{PostId: createdPost.PostId})
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	fmt.Printf("[DELETE] Result: %s\n", deleteRes.Message)
	fmt.Println("--- Test Completed Successfully ---")
}
