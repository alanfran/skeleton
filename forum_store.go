package main

import (
	"fmt"
	"time"

	pg "gopkg.in/pg.v4"
)

// ForumBoard describes a board.
type ForumBoard struct {
	ID          int
	Name        string
	Description string

	// Last Post info
	// user
	// text excerpt
	// time
}

// ForumThread describes a thread.
type ForumThread struct {
	ID      int
	BoardID int
	Topic   string

	// Last Post Info
	// user
	// text excerpt
	// time
}

// ForumPost stores a post.
type ForumPost struct {
	ID       int
	ThreadID int
	UserID   int
	Body     string
	Created  time.Time
}

// ForumStore contains a reference to the database and provides methods that
// interact with it.
type ForumStore struct {
	db *pg.DB
}

// NewForumStore ensures the required tables exist and returns an initialized ForumStore.
func NewForumStore(db *pg.DB) *ForumStore {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS forum_boards (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT )`)
	if err != nil {
		fmt.Println(err)
		panic("Error creating forum_boards table.")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS forum_threads (
    id       SERIAL PRIMARY KEY,
    board_id INT REFERENCES forum_boards (id),
    topic    TEXT NOT NULL )`)
	if err != nil {
		panic("Error creating forum_threads table.")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS forum_posts (
    id        SERIAL PRIMARY KEY,
    thread_id INT REFERENCES forum_threads (id),
    user_id   INT NOT NULL,
    body      TEXT NOT NULL,
    created   TIMESTAMP NOT NULL )`)
	if err != nil {
		panic("Error creating forum_posts table.")
	}
	return &ForumStore{db}
}

// CreateBoard inserts a board into the database or returns an error.
func (s ForumStore) CreateBoard(b ForumBoard) (ForumBoard, error) {
	err := s.db.Create(&b)
	return b, err
}

// CreateThread inserts a thread into the database or returns an error.
func (s ForumStore) CreateThread(t ForumThread) (ForumThread, error) {
	err := s.db.Create(&t)
	return t, err
}

// CreatePost inserts a post into the database or returns an error.
func (s ForumStore) CreatePost(p ForumPost) (ForumPost, error) {
	p.Created = time.Now()
	err := s.db.Create(&p)
	return p, err
}

// GetBoards returns either a slice of all the boards, or an error.
func (s ForumStore) GetBoards() (boards []ForumBoard, err error) {
	err = s.db.Model(&boards).Select()
	return boards, err
}

// GetThreads returns a slice containing all of the threads in a board, or an error.
func (s ForumStore) GetThreads(bid int) (threads []ForumThread, err error) {
	err = s.db.Model(&threads).Select()
	return threads, err
}

// GetThreadRange returns a range of threads from a board, or an error.
func (s ForumStore) GetThreadRange(bid, from, len int) (threads []ForumThread, err error) {
	err = s.db.Model(&threads).Where("board_id = ? AND id <= ?", bid, from).Order("id DESC").Limit(len).Select()
	return threads, err
}

// GetPosts returns a slice containing all of the posts in a thread, or an error.
func (s ForumStore) GetPosts(tid int) (posts []ForumPost, err error) {
	err = s.db.Model(&posts).Where("thread_id = ?", tid).Select()
	return posts, err
}

// GetPostRange returns a slice containing a range of posts from a thread, or an error.
func (s ForumStore) GetPostRange(tid, from, len int) (posts []ForumPost, err error) {
	err = s.db.Model(&posts).Where("thread_id = ? AND id >= ?", tid, from).Order("id ASC").Limit(len).Select()
	return posts, err
}

// PutBoard updates a board's information.
func (s ForumStore) PutBoard(b ForumBoard) (ForumBoard, error) {
	err := s.db.Update(&b)
	return b, err
}

// PutThread updates a thread's information.
func (s ForumStore) PutThread(t ForumThread) (ForumThread, error) {
	err := s.db.Update(&t)
	return t, err
}

// PutPost updates a post.
func (s ForumStore) PutPost(p ForumPost) (ForumPost, error) {
	err := s.db.Update(&p)
	return p, err
}

// DelBoard deletes a board.
func (s ForumStore) DelBoard(bid int) error {
	err := s.db.Delete(&ForumBoard{ID: bid})
	return err
}

// DelThread deletes a thread.
func (s ForumStore) DelThread(tid int) error {
	err := s.db.Delete(&ForumThread{ID: tid})
	return err
}

// DelPost deletes a post.
func (s ForumStore) DelPost(pid int) error {
	err := s.db.Delete(&ForumPost{ID: pid})
	return err
}
