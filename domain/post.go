package domain

import (
	customtime "socio/pkg/time"
)

//easyjson:json
type Post struct {
	ID          uint                  `json:"postId"`
	AuthorID    uint                  `json:"authorId"`
	GroupID     uint                  `json:"groupId,omitempty"`
	Content     string                `json:"content"`
	Attachments []string              `json:"attachments"`
	LikedByIDs  []uint64              `json:"likedBy"`
	CreatedAt   customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt   customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

//easyjson:json
type PostLike struct {
	ID        uint                  `json:"likeId"`
	PostID    uint                  `json:"postId"`
	UserID    uint                  `json:"userId"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

//easyjson:json
type PostWithAuthor struct {
	Post   *Post `json:"post"`
	Author *User `json:"author"`
}

//easyjson:json
type PostWithAuthorAndGroup struct {
	Post   *Post        `json:"post"`
	Author *User        `json:"author"`
	Group  *PublicGroup `json:"group"`
}

//easyjson:json
type GroupPost struct {
	ID        uint                  `json:"id"`
	PostID    uint                  `json:"postId"`
	GroupID   uint                  `json:"groupId"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
