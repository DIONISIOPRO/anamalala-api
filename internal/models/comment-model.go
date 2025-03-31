package models

import (
	"time"
)

// Comment represents a comment on a post
type Comment struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	ReferenceID   string    `bson:"reference_id" json:"reference_id"`
	UserID      string    `bson:"user_id" json:"user_id"`
	Author      Author    `bson:"author" json:"author"`
	Content     string    `bson:"content" json:"content" validate:"required"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	DeletedAt   time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	Comments    Comments  `bson:"comments,omitempty" json:"comments,omitempty"`
	Likes       int       `bson:"likes" json:"likes"`
	Reference   string    `bson:"reference" json:"reference"`
	LikedUserId []string  `bson:"likeduserid,omitempty" json:"likeduserid,omitempty"`
}

// CommentResponse represents the comment data returned to clients
type CommentResponse struct {
	ID        string       `json:"id"`
	User      UserResponse `json:"user"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
	Comments  []Comment    `json:"comments,omitempty"`
}

// CommentCreation represents data for creating a new comment
type CommentCreation struct {
	PostID  string `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// Comments represents a slice of Comment
type Comments []Comment
