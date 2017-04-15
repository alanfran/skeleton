package blog

import "time"

// Post stores a blog post.
type Post struct {
	TableName struct{} `sql:"blog_posts"` // ORMs should look in `blog_posts` instead of `posts`

	ID         int
	Author     int
	AuthorName string `sql:"-"`
	Title      string
	Body       string
	Date       time.Time
	DateString string `sql:"-"`
}

// Storer defines the behavior of a blog store.
type Storer interface {
	CreatePost(b Post) (Post, error)
	GetPost(id int) (Post, error)
	GetRecentPosts(limit int) ([]Post, error)
	GetPostRange(begin, len int) ([]Post, error)
	PutPost(updated Post) error
	DelPost(id int) error
}
