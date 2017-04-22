package blog

import (
	"errors"
	"time"
)

// MockStore stores Posts in memory and provides CRUD methods.
type MockStore struct {
	mem    map[int]Post
	serial int
}

// NewMockStore returns an initialized MockStore.
func NewMockStore() *MockStore {
	return &MockStore{
		mem: map[int]Post{},
	}
}

// GetPost returns the blog post with the provided ID, or an error.
func (s *MockStore) GetPost(id int) (p Post, err error) {
	p, ok := s.mem[id]
	if !ok {
		err = errors.New("Could not find a post with that id.")
	}
	return p, err
}

// GetRecentPosts returns the last n posts, or an error.
func (s *MockStore) GetRecentPosts(limit int) (ps []Post, err error) {
	for i := s.serial - 1; i > 0; i-- {
		if len(ps) == limit {
			break
		}
		p, ok := s.mem[i]
		if ok {
			ps = append(ps, p)
		}
	}
	return ps, err
}

// GetPostRange returns a slice of `len` posts starting from the ID `begin`, or an error.
func (s *MockStore) GetPostRange(begin, length int) (ps []Post, err error) {
	for i := begin; i > 0; i-- {
		if len(ps) == length {
			break
		}
		p, ok := s.mem[i]
		if ok {
			ps = append(ps, p)
		}
	}
	return ps, err
}

// CreatePost inserts the provided post into the memory store.
func (s *MockStore) CreatePost(b Post) (p Post, err error) {
	p.Date = time.Now()
	p.ID = s.serial
	s.serial++

	s.mem[p.ID] = p
	return p, err
}

// PutPost updates the provided post in the memory store.
func (s *MockStore) PutPost(updated Post) error {
	p, err := s.GetPost(updated.ID)
	if err != nil {
		return err
	}
	p.Title = updated.Title
	p.Body = updated.Body

	s.mem[p.ID] = p

	return err
}

// DelPost deletes the post with the provided ID.
func (s *MockStore) DelPost(id int) error {
	_, err := s.GetPost(id)
	delete(s.mem, id)

	return err
}
