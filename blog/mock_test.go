package blog

import (
	"fmt"
	"strconv"
	"testing"
)

var (
	mock Storer
)

func init() {
	mock = NewMockStore()
}

func TestMockCRUD(t *testing.T) {
	oTitle := "Test Title"
	oBody := "Test Body. Lorem ipsum etc."

	// create
	p, err := mock.CreatePost(Post{
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
	err = mock.PutPost(p)
	if err != nil {
		t.Error("Error updating blog post.")
		t.Error(err)
	}

	// get
	p2, err := mock.GetPost(p.ID)
	if err != nil {
		t.Error("Error retrieving updated blog post.")
		t.Error(err)
	}
	if !(p2.Title != oTitle && p2.Title == p.Title && p2.Body != oBody && p2.Body == p.Body) {
		t.Error("Retrieved post does not match updated data.")
	}

	// delete
	err = mock.DelPost(p.ID)
	if err != nil {
		t.Error("Error deleting blog post.")
		t.Error(err)
	}

}

func TestMockGetRanges(t *testing.T) {
	posts := []Post{
		Post{
			Author: 1,
			Title:  "Title 1",
			Body:   "Body 1",
		},
		Post{
			Author: 2,
			Title:  "Title 2",
			Body:   "Body 2",
		},
		Post{
			Author: 3,
			Title:  "Title 3",
			Body:   "Body 3",
		},
		Post{
			Author: 4,
			Title:  "Title 4",
			Body:   "Body 4",
		},
		Post{
			Author: 5,
			Title:  "Title 5",
			Body:   "Body 5",
		},
	}

	// create posts
	for k := range posts {
		_, err := mock.CreatePost(posts[k])
		if err != nil {
			t.Error("Error creating blog post in range.")
			t.Error(err)
		}
	}

	// get range
	postRange, err := mock.GetPostRange(10, 5)
	if err != nil {
		t.Error("Error getting post range.")
		t.Error(err)
	}
	if len(postRange) < 5 {
		t.Error("GetPostRange returned too few posts. " + strconv.Itoa(len(postRange)))
	}

	// get recent
	latest, err := mock.GetRecentPosts(3)
	if err != nil {
		t.Error("Error retrieving recent posts.")
		t.Error(err)
	}
	if len(latest) < 3 {
		t.Error("GetRecentPosts returned too few posts. " + strconv.Itoa(len(postRange)))
	}

	// delete
	for k := range postRange {
		err = mock.DelPost(postRange[k].ID)
		if err != nil {
			t.Error("Error deleting post " + fmt.Sprintf("%v", postRange[k].ID))
			t.Error(err)
		}
	}

}
