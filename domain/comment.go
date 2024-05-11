package domain

import customtime "socio/pkg/time"

type Comment struct {
	ID        uint                  `json:"id"`
	Content   string                `json:"content"`
	PostID    uint                  `json:"postId"`
	AuthorID  uint                  `json:"authorId"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type CommentLike struct {
	ID        uint                  `json:"id"`
	CommentID uint                  `json:"commentId"`
	UserID    uint                  `json:"userId"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
