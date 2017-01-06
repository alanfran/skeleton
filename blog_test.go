package main

import (
	"strconv"
	"testing"

	"gopkg.in/pg.v4"
)

func init() {
	if db == nil {
		db = pg.Connect(&pg.Options{
			Addr:     dbAddr,
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

	blog = NewBlogStore(db)
}

func TestBlogCRUD(t *testing.T) {
	oTitle := "Test Title"
	oBody := "Test Body. Lorem ipsum etc."

	// create
	p, err := blog.CreatePost(BlogPost{
		Author: 123,
		Title:  oTitle,
		Body:   oBody})
	if err != nil {
		t.Error("Error creating blog post.")
		t.Error(err)
	}

	// update
	p.Title = "Something different."
	p.Body = "This should be different"
	err = blog.PutPost(p)
	if err != nil {
		t.Error("Error updating blog post.")
		t.Error(err)
	}

	// get
	p2, err := blog.GetPost(p.ID)
	if err != nil {
		t.Error("Error retrieving updated blog post.")
		t.Error(err)
	}
	if !(p2.Title != oTitle && p2.Title == p.Title && p2.Body != oBody && p2.Body == p.Body) {
		t.Error("Retrieved post does not match updated data.")
	}

	// delete
	err = blog.DelPost(p.ID)
	if err != nil {
		t.Error("Error deleting blog post.")
		t.Error(err)
	}

}

func TestGetPostRange(t *testing.T) {
	posts := []BlogPost{
		BlogPost{
			Author: 1,
			Title:  "Title 1",
			Body:   "Body 1",
		},
		BlogPost{
			Author: 2,
			Title:  "Title 2",
			Body:   "Body 2",
		},
		BlogPost{
			Author: 3,
			Title:  "Title 3",
			Body:   "Body 3",
		},
		BlogPost{
			Author: 4,
			Title:  "Title 4",
			Body:   "Body 4",
		},
		BlogPost{
			Author: 5,
			Title:  "Title 5",
			Body:   "Body 5",
		},
	}

	// create posts
	for _, p := range posts {
		_, err := blog.CreatePost(p)
		if err != nil {
			t.Error("Error creating blog post in range.")
			t.Error(err)
		}
	}

	// get range
	postRange, err := blog.GetPostRange(999999, 5)
	if err != nil {
		t.Error("Error getting post range.")
		t.Error(err)
	}
	if len(postRange) < 5 {
		t.Error("GetPostRange returned too few posts. " + strconv.Itoa(len(postRange)))
	}

	// get recent
	latest, err := blog.GetRecentPosts(3)
	if err != nil {
		t.Error("Error retrieving recent posts.")
		t.Error(err)
	}
	if len(latest) < 3 {
		t.Error("GetRecentPosts returned too few posts. " + strconv.Itoa(len(postRange)))
	}

	// delete
	for _, p := range postRange {
		err = blog.DelPost(p.ID)
		if err != nil {
			t.Error("Error deleting post.")
			t.Error(err)
		}
	}

}
