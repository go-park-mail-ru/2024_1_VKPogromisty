package posts

import (
	"context"
	"socio/domain"
	"socio/errors"
)

func (s *Service) GetCommentsByPostID(ctx context.Context, postID uint) (comments []*domain.Comment, err error) {
	_, err = s.PostsStorage.GetPostByID(ctx, postID)
	if err != nil {
		return
	}

	comments, err = s.PostsStorage.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return
	}

	for _, comment := range comments {
		s.Sanitizer.SanitizeComment(comment)
	}

	return
}

func (s *Service) CreateComment(ctx context.Context, comment *domain.Comment) (newComment *domain.Comment, err error) {
	s.Sanitizer.SanitizeComment(comment)

	if len(comment.Content) == 0 {
		err = errors.ErrInvalidData
		return
	}

	_, err = s.PostsStorage.GetPostByID(ctx, comment.PostID)
	if err != nil {
		return
	}

	newComment, err = s.PostsStorage.StoreComment(ctx, comment)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizeComment(newComment)

	return
}

func (s *Service) UpdateComment(ctx context.Context, comment *domain.Comment) (updatedComment *domain.Comment, err error) {
	oldComment, err := s.PostsStorage.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return
	}

	if oldComment.AuthorID != comment.AuthorID {
		err = errors.ErrForbidden
		return
	}

	if len(comment.Content) == 0 {
		err = errors.ErrInvalidData
		return
	}

	updatedComment, err = s.PostsStorage.UpdateComment(ctx, comment)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizeComment(updatedComment)

	return
}

func (s *Service) DeleteComment(ctx context.Context, comment *domain.Comment) (err error) {
	oldComment, err := s.PostsStorage.GetCommentByID(ctx, comment.ID)
	if err != nil {
		return
	}

	if oldComment.AuthorID != comment.AuthorID {
		err = errors.ErrForbidden
		return
	}

	err = s.PostsStorage.DeleteComment(ctx, comment.ID)
	if err != nil {
		return
	}

	return
}

func (s *Service) LikeComment(ctx context.Context, commentLike *domain.CommentLike) (newLike *domain.CommentLike, err error) {
	_, err = s.PostsStorage.GetCommentByID(ctx, commentLike.CommentID)
	if err != nil {
		return
	}

	newLike, err = s.PostsStorage.StoreCommentLike(ctx, commentLike)
	if err != nil {
		return
	}

	return
}

func (s *Service) UnlikeComment(ctx context.Context, commentLike *domain.CommentLike) (err error) {
	err = s.PostsStorage.DeleteCommentLike(ctx, commentLike)
	if err != nil {
		return
	}

	return
}
