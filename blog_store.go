package main

import (
	"time"

	"gopkg.in/pg.v4"
)

// BlogPost stores a blog post.
type BlogPost struct {
	ID         int
	Author     int
	AuthorName string `sql:"-"`
	Title      string
	Body       string
	Date       time.Time
	DateString string `sql:"-"`
}

// BlogStore contains a reference to the database and provides CRUD methods.
type BlogStore struct {
	db  *pg.DB
	mem map[int]BlogPost
}

// NewBlogStore ensures the `blog_posts` table exists in the database and returns an initialized BlogStore.
func NewBlogStore(db *pg.DB) *BlogStore {
	// Create Table in db
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS blog_posts (
		id SERIAL PRIMARY KEY,
		author INT NOT NULL,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		date TIMESTAMP NOT NULL)`)
	if err != nil {
		panic(err)
	}

	return &BlogStore{
		db:  db,
		mem: make(map[int]BlogPost),
	}
}

// GetPost returns the blog post with the provided ID, or an error.
func (s BlogStore) GetPost(id int) (BlogPost, error) {
	var p BlogPost
	p.ID = id
	p.Date = time.Now()

	err := s.db.Select(&p)

	return p, err
}

// GetRecentPosts returns the last n posts, or an error.
func (s BlogStore) GetRecentPosts(limit int) ([]BlogPost, error) {
	var posts []BlogPost
	err := s.db.Model(&posts).Order("id DESC").Limit(limit).Select()

	for k, p := range posts {
		u, _ := users.Get(p.Author)
		posts[k].AuthorName = u.Name
		posts[k].DateString = p.Date.Format("Jan 2, 2006 at 15:04")
	}

	return posts, err
}

// GetPostRange returns a slice of `len` posts starting from the ID `begin`, or an error.
func (s BlogStore) GetPostRange(begin, len int) ([]BlogPost, error) {
	var posts []BlogPost
	err := s.db.Model(&posts).Where("id <= ?", begin).Order("id DESC").Limit(len).Select()

	return posts, err
}

// CreatePost inserts the provided post into the database.
func (s BlogStore) CreatePost(b BlogPost) (BlogPost, error) {

	b.Date = time.Now()
	err := s.db.Create(&b)

	return b, err
}

// PutPost updates the provided post in the database.
func (s BlogStore) PutPost(updated BlogPost) error {
	p, err := s.GetPost(updated.ID)
	if err != nil {
		return err
	}
	p.Title = updated.Title
	p.Body = updated.Body

	err = s.db.Update(&p)

	return err
}

// DelPost deletes the post with the provided ID.
func (s BlogStore) DelPost(id int) error {
	err := s.db.Delete(&BlogPost{ID: id})

	return err
}
