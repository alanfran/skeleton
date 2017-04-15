package blog

import (
	"time"

	"gopkg.in/pg.v4"
)

// PgStore contains a reference to the database and provides CRUD methods.
type PgStore struct {
	db *pg.DB
}

// NewPgStore ensures the `blog_posts` table exists in the database and returns an initialized PgStore.
func NewPgStore(db *pg.DB) *PgStore {
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

	return &PgStore{
		db: db,
	}
}

// GetPost returns the blog post with the provided ID, or an error.
func (s PgStore) GetPost(id int) (Post, error) {
	var p Post
	p.ID = id
	p.Date = time.Now()

	err := s.db.Select(&p)

	return p, err
}

// GetRecentPosts returns the last n posts, or an error.
func (s PgStore) GetRecentPosts(limit int) ([]Post, error) {
	var posts []Post
	err := s.db.Model(&posts).Order("id DESC").Limit(limit).Select()

	return posts, err
}

// GetPostRange returns a slice of `len` posts starting from the ID `begin`, or an error.
func (s PgStore) GetPostRange(begin, len int) ([]Post, error) {
	var posts []Post
	err := s.db.Model(&posts).Where("id <= ?", begin).Order("id DESC").Limit(len).Select()

	return posts, err
}

// CreatePost inserts the provided post into the database.
func (s PgStore) CreatePost(b Post) (Post, error) {

	b.Date = time.Now()
	err := s.db.Create(&b)

	return b, err
}

// PutPost updates the provided post in the database.
func (s PgStore) PutPost(updated Post) error {
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
func (s PgStore) DelPost(id int) error {
	err := s.db.Delete(&Post{ID: id})

	return err
}
