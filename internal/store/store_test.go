package store

import (
	"sync"
	"testing"

	"github.com/Chetas1/grpc-blog-service/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewBlogStore(t *testing.T) {
	store := NewBlogStore()
	assert.NotNil(t, store)
}

func TestCreate(t *testing.T) {
	store := NewBlogStore()

	post := &proto.Post{
		PostId:  "post-1",
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  "John Doe",
		Tags:    []string{"test", "golang"},
	}

	store.Create(post)

	retrievedPost, err := store.Get("post-1")
	require.NoError(t, err)
	assert.Equal(t, post.PostId, retrievedPost.PostId)
	assert.Equal(t, post.Title, retrievedPost.Title)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, post.Author, retrievedPost.Author)
	assert.Equal(t, post.Tags, retrievedPost.Tags)
}

func TestCreateMultiplePosts(t *testing.T) {
	store := NewBlogStore()

	posts := []*proto.Post{
		{
			PostId:  "post-1",
			Title:   "First Post",
			Content: "Content 1",
			Author:  "Author 1",
			Tags:    []string{"tag1"},
		},
		{
			PostId:  "post-2",
			Title:   "Second Post",
			Content: "Content 2",
			Author:  "Author 2",
			Tags:    []string{"tag2"},
		},
		{
			PostId:  "post-3",
			Title:   "Third Post",
			Content: "Content 3",
			Author:  "Author 3",
			Tags:    []string{"tag3"},
		},
	}

	for _, post := range posts {
		store.Create(post)
	}

	for _, post := range posts {
		retrievedPost, err := store.Get(post.PostId)
		require.NoError(t, err)
		assert.Equal(t, post, retrievedPost)
	}
}

func TestGet(t *testing.T) {
	store := NewBlogStore()

	post := &proto.Post{
		PostId:  "test-post",
		Title:   "Get Test",
		Content: "Testing get method",
		Author:  "Tester",
		Tags:    []string{"test"},
	}

	store.Create(post)

	retrievedPost, err := store.Get("test-post")
	require.NoError(t, err)
	assert.Equal(t, post, retrievedPost)
}

func TestGetNotFound(t *testing.T) {
	store := NewBlogStore()

	_, err := store.Get("non-existent")
	assert.Error(t, err)
	assert.Equal(t, ErrPostNotFound, err)
}

func TestUpdate(t *testing.T) {
	store := NewBlogStore()

	post := &proto.Post{
		PostId:  "update-test",
		Title:   "Original Title",
		Content: "Original Content",
		Author:  "Original Author",
		Tags:    []string{"original"},
	}

	store.Create(post)

	newTitle := "Updated Title"
	newContent := "Updated Content"
	newAuthor := "Updated Author"
	newTags := []string{"updated", "modified"}

	updatedPost, err := store.Update("update-test", newTitle, newContent, newAuthor, newTags)
	require.NoError(t, err)

	assert.Equal(t, "update-test", updatedPost.PostId)
	assert.Equal(t, newTitle, updatedPost.Title)
	assert.Equal(t, newContent, updatedPost.Content)
	assert.Equal(t, newAuthor, updatedPost.Author)
	assert.Equal(t, newTags, updatedPost.Tags)

	// Verify the update persisted
	retrievedPost, err := store.Get("update-test")
	require.NoError(t, err)
	assert.Equal(t, newTitle, retrievedPost.Title)
	assert.Equal(t, newContent, retrievedPost.Content)
	assert.Equal(t, newAuthor, retrievedPost.Author)
	assert.Equal(t, newTags, retrievedPost.Tags)
}

func TestUpdateNotFound(t *testing.T) {
	store := NewBlogStore()

	_, err := store.Update("non-existent", "title", "content", "author", []string{})
	assert.Error(t, err)
	assert.Equal(t, ErrPostNotFound, err)
}

func TestDelete(t *testing.T) {
	store := NewBlogStore()

	post := &proto.Post{
		PostId:  "delete-test",
		Title:   "To Be Deleted",
		Content: "This post will be deleted",
		Author:  "Deleter",
		Tags:    []string{"delete"},
	}

	store.Create(post)

	// Verify post exists
	_, err := store.Get("delete-test")
	require.NoError(t, err)

	// Delete the post
	err = store.Delete("delete-test")
	require.NoError(t, err)

	// Verify post is deleted
	_, err = store.Get("delete-test")
	assert.Error(t, err)
	assert.Equal(t, ErrPostNotFound, err)
}

func TestDeleteNotFound(t *testing.T) {
	store := NewBlogStore()

	err := store.Delete("non-existent")
	assert.Error(t, err)
	assert.Equal(t, ErrPostNotFound, err)
}

func TestConcurrentCreate(t *testing.T) {
	store := NewBlogStore()
	numGoroutines := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 1; i <= numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			postID := generatePostID(id)
			post := &proto.Post{
				PostId:  postID,
				Title:   "Concurrent Post",
				Content: "Concurrent content",
				Author:  "Concurrent Author",
				Tags:    []string{"concurrent"},
			}
			store.Create(post)
		}(i)
	}

	wg.Wait()

	// Verify all posts were created
	for i := 1; i <= numGoroutines; i++ {
		postID := generatePostID(i)
		_, err := store.Get(postID)
		require.NoError(t, err)
	}
}

func TestConcurrentGetAndCreate(t *testing.T) {
	store := NewBlogStore()

	post := &proto.Post{
		PostId:  "shared-post",
		Title:   "Shared Post",
		Content: "Shared content",
		Author:  "Shared Author",
		Tags:    []string{"shared"},
	}

	store.Create(post)

	numGoroutines := 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := store.Get("shared-post")
			require.NoError(t, err)
		}()
	}

	wg.Wait()
}

func TestConcurrentUpdateAndRead(t *testing.T) {
	store := NewBlogStore()

	post := &proto.Post{
		PostId:  "concurrent-update",
		Title:   "Original",
		Content: "Original content",
		Author:  "Original author",
		Tags:    []string{"original"},
	}

	store.Create(post)

	numGoroutines := 50
	var wg sync.WaitGroup

	// Half goroutines will read
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Get("concurrent-update")
			require.NoError(t, err)
		}()
	}

	// Half goroutines will update
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := store.Update("concurrent-update", "Updated", "New content", "New author", []string{"new"})
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Verify final state
	retrievedPost, err := store.Get("concurrent-update")
	require.NoError(t, err)
	assert.Equal(t, "Updated", retrievedPost.Title)
}

func TestConcurrentDelete(t *testing.T) {
	store := NewBlogStore()

	post := &proto.Post{
		PostId:  "delete-concurrent",
		Title:   "To delete",
		Content: "Content",
		Author:  "Author",
		Tags:    []string{"delete"},
	}

	store.Create(post)

	// Try to delete from multiple goroutines
	numGoroutines := 10
	var wg sync.WaitGroup
	var successCount int
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := store.Delete("delete-concurrent")
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				successCount++
			}
		}()
	}

	wg.Wait()

	// Only one delete should succeed
	assert.Equal(t, 1, successCount)

	// Verify post is deleted
	_, err := store.Get("delete-concurrent")
	assert.Equal(t, ErrPostNotFound, err)
}

func TestCreateAndUpdateWithTimestamp(t *testing.T) {
	store := NewBlogStore()

	timestamp := timestamppb.Now()
	post := &proto.Post{
		PostId:          "timestamp-test",
		Title:           "Timestamp Test",
		Content:         "Testing timestamp",
		Author:          "Author",
		Tags:            []string{"timestamp"},
		PublicationDate: timestamp,
	}

	store.Create(post)

	retrievedPost, err := store.Get("timestamp-test")
	require.NoError(t, err)
	assert.Equal(t, timestamp, retrievedPost.PublicationDate)
}

// Helper function to generate post IDs for concurrent tests
func generatePostID(id int) string {
	return "post-" + string(rune(id))
}
