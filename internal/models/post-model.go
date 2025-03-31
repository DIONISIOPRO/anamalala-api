package models

import (
	"time"
)

// PostType represents the type of post
type PostType string

const (
	PostTypeText  PostType = "text"
	PostTypeImage PostType = "image"
	PostTypeVideo PostType = "video"
	PostTypeFile  PostType = "file"
)

// Post represents a post in the chat room
type Post struct {
	ID      string   `bson:"_id,omitempty" json:"id,omitempty"`
	UserID  string   `bson:"user_id" json:"user_id"`
	Author  Author   `bson:"author" json:"author"`
	Content string   `bson:"content" json:"content" validate:"required"`
	Type    PostType `bson:"type" json:"type"`
	// Attachments []string  `bson:"attachments,omitempty" json:"attachments,omitempty"`
	Comments    []Comment `bson:"comments" json:"comments"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	Likes       int       `bson:"likes" json:"likes"`
	LikedUserId []string  `bson:"likeduserid" json:"likeduserid"`
}

type Author struct {
	Name string `bson:"name" json:"name"`
	ID   string `bson:"id" json:"id"`
}

// PostResponse represents the post data returned to clients
type PostResponse struct {
	ID          string            `json:"id"`
	User        UserResponse      `json:"user"`
	Content     string            `json:"content"`
	Type        PostType          `json:"type"`
	Attachments []string          `json:"attachments"`
	Comments    []CommentResponse `json:"comments"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	LikedUserId []User            `json:"likeduserid"`
	Likes       int               `json:"likes"`
}

// PostCreation represents data for creating a new post
type PostCreation struct {
	Content     string   `json:"content" validate:"required"`
	Type        PostType `json:"type" validate:"required"`
	Attachments []string `json:"attachments"`
}

// Posts represents a slice of Post
type Posts []Post
