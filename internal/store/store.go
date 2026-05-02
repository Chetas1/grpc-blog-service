// Package store provides the BlogStore interface and an in-memory
// implementation suitable for development, testing, and reference.
//
// Production deployments should provide a Postgres / Redis / DynamoDB
// implementation behind the same interface; the gRPC handlers in
// `server/` depend only on BlogStore.
package store

import (
	"errors"
	"sync"

	"github.com/Chetas-Patil/grpc-blog-service/proto"
)

// ErrPostNotFound is returned by Get, Update, and Delete when the requested
// post ID does not exist in the store.
var ErrPostNotFound = errors.New("post not found")

// BlogStore is the persistence interface for blog posts.
//
// All implementations must be safe for concurrent use by multiple goroutines.
type BlogStore interface {
	Create(post *proto.Post)
	Get(id string) (*proto.Post, error)
	Update(id string, title, content, author string, tags []string) (*proto.Post, error)
	Delete(id string) error
	ReadAll() ([]*proto.Post, error)
}

type blogStore struct {
	mu    sync.RWMutex
	posts map[string]*proto.Post
}

// NewBlogStore returns an in-memory BlogStore. Data is lost when the process
// exits; use a persistent backend (Postgres, Redis, etc.) implementing the
// BlogStore interface for production traffic.
func NewBlogStore() BlogStore {
	return &blogStore{
		posts: make(map[string]*proto.Post),
	}
}

func (s *blogStore) Create(post *proto.Post) {
	if post == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.posts[post.PostId] = post
}

func (s *blogStore) Get(id string) (*proto.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (s *blogStore) Update(id, title, content, author string, tags []string) (*proto.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, ErrPostNotFound
	}

	post.Title = title
	post.Content = content
	post.Author = author
	post.Tags = tags

	return post, nil
}

func (s *blogStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return ErrPostNotFound
	}
	delete(s.posts, id)
	return nil
}

func (s *blogStore) ReadAll() ([]*proto.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	response := make([]*proto.Post, 0, len(s.posts))
	for _, post := range s.posts {
		response = append(response, post)
	}
	return response, nil
}
