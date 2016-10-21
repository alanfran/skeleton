package main

import (
	"testing"

	pg "gopkg.in/pg.v4"
)

func init() {
	if db == nil {
		db = pg.Connect(&pg.Options{
			User:     dbUser,
			Password: dbPassword,
			Database: dbTestDatabase,
		})
		// verify connection
		_, err := db.Exec(`SELECT 1`)
		if err != nil {
			panic("Error connecting to the database.")
		}
	}

	forum = NewForumStore(db)
}

func TestForumCRUD(t *testing.T) {
	// create a board
	testBoard := ForumBoard{
		Name:        "Test Board",
		Description: "This board will be used in unit tests.",
	}
	board, err := forum.CreateBoard(testBoard)
	if err != nil {
		t.Error("Error creating board.")
		t.Error(err)
	}

	// Create a thread
	testThread := ForumThread{
		BoardID: board.ID,
		Topic:   "This thread is for testing purposes.",
	}
	thread, err := forum.CreateThread(testThread)
	if err != nil {
		t.Error("Error creating thread.")
		t.Error(err)
	}

	// Create a post
	testPost := ForumPost{
		ThreadID: thread.ID,
		Body:     "Post content goes here",
	}
	post, err := forum.CreatePost(testPost)
	if err != nil {
		t.Error("Error creating post.")
		t.Error(err)
	}

	// GET

	boards, err := forum.GetBoards()
	if err != nil {
		t.Error("Error getting boards.")
		t.Error(err)
	}
	if len(boards) < 1 {
		t.Error("Not enough boards returned.")
	}

	threads, err := forum.GetThreads(board.ID)
	if err != nil {
		t.Error("Error getting threads.")
		t.Error(err)
	}
	if len(threads) < 1 {
		t.Error("Not enough threads returned.")
	}

	threadRange, err := forum.GetThreadRange(board.ID, 99999, 1)
	if err != nil {
		t.Error("Error getting thread range.")
		t.Error(err)
	}
	if len(threadRange) < 1 {
		t.Error("Not enough threads returned in range.")
	}

	posts, err := forum.GetPosts(thread.ID)
	if err != nil {
		t.Error("Error getting posts.")
		t.Error(err)
	}
	if len(posts) < 1 {
		t.Error("Not enough posts returned.")
	}

	postRange, err := forum.GetPostRange(thread.ID, 1, 1)
	if err != nil {
		t.Error("Error getting post range.")
		t.Error(err)
	}
	if len(postRange) < 1 {
		t.Error("Not enough posts returned.")
	}

	// PUT

	b2, err := forum.PutBoard(ForumBoard{
		ID:          board.ID,
		Name:        "Modified board.",
		Description: "This board has been modified.",
	})
	if err != nil {
		t.Error("Error updating board.")
		t.Error(err)
	}
	if b2.Name == board.Name || b2.Description == board.Description {
		t.Error("Board did not properly update.")
	}

	t2, err := forum.PutThread(ForumThread{
		ID:      thread.ID,
		BoardID: thread.BoardID,
		Topic:   "Modified thread.",
	})
	if err != nil {
		t.Error("Error updating thread.")
		t.Error(err)
	}
	if t2.Topic == thread.Topic {
		t.Error("Thread did not properly update.")
	}

	p2, err := forum.PutPost(ForumPost{
		ID:       post.ID,
		ThreadID: post.ThreadID,
		Body:     "Modified post.",
		Created:  post.Created,
	})
	if err != nil {
		t.Error("Error updating post.")
		t.Error(err)
	}
	if p2.Body == post.Body {
		t.Error("Post did not properly update.")
	}

	// Delete test post
	err = forum.DelPost(post.ID)
	if err != nil {
		t.Error("Error deleting post.")
		t.Error(err)
	}

	// Delete test thread
	err = forum.DelThread(thread.ID)
	if err != nil {
		t.Error("Error deleting thread.")
		t.Error(err)
	}

	// Delete test board
	err = forum.DelBoard(board.ID)
	if err != nil {
		t.Error("Error deleting board.")
		t.Error(err)
	}
}
