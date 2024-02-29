package services

import "socio/utils"

type Post struct {
	ID           uint             `json:"postId"`
	AuthorID     uint             `json:"authorId"`
	Text         string           `json:"text"`
	Attachments  []string         `json:"attachments"`
	CreationDate utils.CustomTime `json:"creationDate,omitempty"`
}
