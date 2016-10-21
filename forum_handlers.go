package main

import "github.com/gin-gonic/gin"

// Forum Home: Board List
func forumHomeH(c *gin.Context) {
	c.String(200, "Ok")
}

// Board View: Thread List
func forumBoardH(c *gin.Context) {
	c.String(200, "Ok")
}

// Thread View: Post List
func forumThreadH(c *gin.Context) {
	c.String(200, "Ok")
}

// API

// Create Board
func postBoardH(c *gin.Context) {
	c.String(200, "Ok")
}

// Create Thread
func postThreadH(c *gin.Context) {
	c.String(200, "Ok")
}

// Create Post
func postPostH(c *gin.Context) {
	c.String(200, "Ok")
}

// Get Boards
func getBoardsH(c *gin.Context) {
	c.String(200, "Ok")
}

// Get Threads
func getThreadsH(c *gin.Context) {
	c.String(200, "Ok")
}

// Get Posts
func getPostsH(c *gin.Context) {
	c.String(200, "Ok")
}

// Get Post
func getPostH(c *gin.Context) {
	c.String(200, "Ok")
}

// Update Board
func putBoardH(c *gin.Context) {
	c.String(200, "Ok")
}

// Update Thread
func putThreadH(c *gin.Context) {
	c.String(200, "Ok")
}

// Update Post
func putPostH(c *gin.Context) {
	c.String(200, "Ok")
}

// Delete Board
func deleteBoardH(c *gin.Context) {
	c.String(200, "Ok")
}

// Delete Thread
func deleteThreadH(c *gin.Context) {
	c.String(200, "Ok")
}

// Delete Post
func deletePostH(c *gin.Context) {
	c.String(200, "Ok")
}
