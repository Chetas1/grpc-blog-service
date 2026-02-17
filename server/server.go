package main

import (
	"context"
	"log"

	"github.com/Chetas1/grpc-blog-service/internal/store"
	"github.com/Chetas1/grpc-blog-service/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedBlogServiceServer
	store store.BlogStore
}

func (s *server) CreatePost(ctx context.Context, req *proto.CreatePostRequest) (*proto.Post, error) {

	if req.Title == "" || req.Author == "" {
		return nil, status.Error(codes.InvalidArgument, "title and author are required")
	}

	post := &proto.Post{
		PostId:          uuid.New().String(),
		Title:           req.Title,
		Content:         req.Content,
		Author:          req.Author,
		PublicationDate: req.PublicationDate,
		Tags:            req.Tags,
	}
	s.store.Create(post)
	log.Printf("[CREATE] Successfully created post ID: %s", post.PostId)
	return post, nil
}

func (s *server) ReadPost(ctx context.Context, req *proto.ReadPostRequest) (*proto.Post, error) {
	post, err := s.store.Get(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found: %s", req.PostId)
	}
	log.Printf("[READ] Successfully retrieved post ID: %s", req.PostId)
	return post, nil
}

func (s *server) UpdatePost(ctx context.Context, req *proto.UpdatePostRequest) (*proto.Post, error) {
	updatedPost, err := s.store.Update(req.PostId, req.Title, req.Content, req.Author, req.Tags)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to update, post not found: %s", req.PostId)
	}

	log.Printf("[UPDATE] Successfully updated post ID: %s", req.PostId)
	return updatedPost, nil
}

func (s *server) DeletePost(ctx context.Context, req *proto.DeletePostRequest) (*proto.DeletePostResponse, error) {
	err := s.store.Delete(req.PostId)
	if err != nil {
		return &proto.DeletePostResponse{Success: false, Message: err.Error()}, status.Errorf(codes.NotFound, "post not found: %s", req.PostId)
	}

	log.Printf("[DELETE] Successfully deleted post ID: %s", req.PostId)
	return &proto.DeletePostResponse{Success: true, Message: "Post deleted successfully"}, nil
}

func (s *server) ReadAll(ctx context.Context, req *proto.ReadAllRequest) (*proto.ReadAllResponse, error) {

	posts, err := s.store.ReadAll()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "post not found")
	}
	log.Printf("[READALL] Successfully retrieved post ")
	return &proto.ReadAllResponse{
		Posts: posts,
	}, err
}
