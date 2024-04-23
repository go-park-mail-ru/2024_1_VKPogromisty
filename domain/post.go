package domain

import (
	customtime "socio/pkg/time"
)

type Post struct {
	ID          uint                  `json:"postId"`
	AuthorID    uint                  `json:"authorId"`
	Content     string                `json:"content"`
	Attachments []string              `json:"attachments"`
	LikedByIDs  []uint64              `json:"likedBy"`
	CreatedAt   customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt   customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type PostLike struct {
	ID        uint                  `json:"likeId"`
	PostID    uint                  `json:"postId"`
	UserID    uint                  `json:"userId"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type PostWithAuthor struct {
	Post   *Post `json:"post"`
	Author *User `json:"author"`
}
