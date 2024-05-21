package sanitizer

import (
	"socio/domain"

	"github.com/microcosm-cc/bluemonday"
)

type Sanitizer struct {
	sanitizer *bluemonday.Policy
}

func NewSanitizer(sanitizer *bluemonday.Policy) *Sanitizer {
	return &Sanitizer{
		sanitizer: sanitizer,
	}
}

func (s *Sanitizer) Sanitize(input string) string {
	return s.sanitizer.Sanitize(input)
}

func (s *Sanitizer) SanitizeUser(user *domain.User) {
	if user == nil {
		return
	}

	user.FirstName = s.Sanitize(user.FirstName)
	user.LastName = s.Sanitize(user.LastName)
	user.Email = s.Sanitize(user.Email)
	user.Avatar = s.Sanitize(user.Avatar)
}

func (s *Sanitizer) SanitizePost(post *domain.Post) {
	if post == nil {
		return
	}

	post.Content = s.Sanitize(post.Content)

	for i := range post.Attachments {
		post.Attachments[i] = s.Sanitize(post.Attachments[i])
	}
}

func (s *Sanitizer) SanitizePostWithAuthor(post *domain.PostWithAuthor) {
	if post == nil {
		return
	}

	s.SanitizePost(post.Post)
	s.SanitizeUser(post.Author)
}

func (s *Sanitizer) SanitizePersonalMessage(message *domain.PersonalMessage) {
	if message == nil {
		return
	}

	message.Content = s.Sanitize(message.Content)
}

func (s *Sanitizer) SanitizeDialog(dialog *domain.Dialog) {
	if dialog == nil {
		return
	}

	s.SanitizeUser(dialog.User1)
	s.SanitizeUser(dialog.User2)
	s.SanitizePersonalMessage(dialog.LastMessage)
}

func (s *Sanitizer) SanitizePublicGroup(publicGroup *domain.PublicGroup) {
	if publicGroup == nil {
		return
	}

	publicGroup.Name = s.Sanitize(publicGroup.Name)
	publicGroup.Description = s.Sanitize(publicGroup.Description)
	publicGroup.Avatar = s.Sanitize(publicGroup.Avatar)
}

func (s *Sanitizer) SanitizeComment(comment *domain.Comment) {
	if comment == nil {
		return
	}

	comment.Content = s.Sanitize(comment.Content)
}

func (s *Sanitizer) SanitizeComments(comments []*domain.Comment) {
	for _, comment := range comments {
		s.SanitizeComment(comment)
	}
}

func (s *Sanitizer) SanitizeSticker(sticker *domain.Sticker) {
	if sticker == nil {
		return
	}

	sticker.Name = s.Sanitize(sticker.Name)
}

func (s *Sanitizer) SanitizeStickers(stickers []*domain.Sticker) {
	for _, sticker := range stickers {
		s.SanitizeSticker(sticker)
	}
}
