package store

import (
	"errors"
	"sync"

	"github.com/Chetas1/grpc-blog-service/proto"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

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

func NewBlogStore() BlogStore {
	return &blogStore{
		posts: make(map[string]*proto.Post),
	}
}

func (s *blogStore) Create(post *proto.Post) {
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

func (s *blogStore) Update(id string, title, content, author string, tags []string) (*proto.Post, error) {
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

	response := make([]*proto.Post, 0)
	for _, post := range s.posts {
		response = append(response, post)
	}
	return response, nil
}
