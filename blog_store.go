package main

import (
  "time"
  "gopkg.in/pg.v4"
)
 type BlogPost struct{
  ID      int
  Author  int
  Title   string
  Body    string
  Date    time.Time
}

// type BlogStore
type BlogStore struct {
  db          *pg.DB
}

func NewBlogStore(db *pg.DB) *BlogStore{
  // Create Table in db
  _, err := db.Exec(`CREATE TABLE IF NOT EXISTS blog_posts (id SERIAL, author TEXT, title TEXT, body TEXT, date TIMESTAMP)`)
  if err != nil {
    panic(err)
  }

  return &BlogStore {db: db}
}

// Get Post (id)
func (s BlogStore) GetPost(id int) (BlogPost, error) {
  var p BlogPost
  p.ID = id

  err := s.db.Select(&p)

  return p, err
}

func (s BlogStore) GetRecentPosts(limit int) ([]BlogPost, error) {
  var posts []BlogPost
  err := s.db.Model(&posts).Order("id DESC").Limit(limit).Select()

  return posts, err
}

func (s BlogStore) GetPostRange(begin, len int) ([]BlogPost, error) {
  var posts []BlogPost
  err := s.db.Model(&posts).Where("id >= ?", begin).Order("id DESC").Limit(len).Select()

  return posts, err
}

// Create Post
func (s BlogStore) CreatePost(b BlogPost) (BlogPost, error) {
  b.Date = time.Now()

  err := s.db.Create(&b)

  return b, err
}

// Put Post
func (s BlogStore) PutPost(b BlogPost) error {
  err := s.db.Update(&b)

  return err
}

// Del Post
func (s BlogStore) DelPost(id int) error {
  err := s.db.Delete(&BlogPost{ID: id})

  return err
}
