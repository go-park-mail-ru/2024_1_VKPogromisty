package posts

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/sanitizer"

	"github.com/microcosm-cc/bluemonday"
)

const (
	DefaultPostsAmount      = uint(20)
	DefaultLikedPostsAmount = uint(20)
)

//easyjson:json
type PostInput struct {
	Content     string   `json:"content"`
	AuthorID    uint     `json:"authorId"`
	Attachments []string `json:"attachments"`
}

//easyjson:json
type ListUserPostsInput struct {
	UserID      uint `json:"userId"`
	LastPostID  uint `json:"lastPostId"`
	PostsAmount uint `json:"postsAmount"`
}

//easyjson:json
type ListUserFriendsPostsInput struct {
	LastPostID  uint `json:"lastPostId"`
	PostsAmount uint `json:"postsAmount"`
}

//easyjson:json
type PostUpdateInput struct {
	PostID              uint     `json:"postId"`
	Content             string   `json:"content"`
	AttachmentsToAdd    []string `json:"attachmentsToAdd"`
	AttachmentsToDelete []string `json:"attachmentsToDelete"`
}

//easyjson:json
type DeletePostInput struct {
	PostID uint `json:"postId"`
}

//easyjson:json
type LikeWithPost struct {
	Like *domain.PostLike `json:"like"`
	Post *domain.Post     `json:"post"`
}

//easyjson:json
type LikeWithPostAndUser struct {
	LikeWithPost
	User *domain.User `json:"likedBy"`
}

type PostsStorage interface {
	GetPostByID(ctx context.Context, postID uint) (post *domain.Post, err error)
	GetUserPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error)
	GetUserFriendsPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error)
	StorePost(ctx context.Context, post *domain.Post) (newPost *domain.Post, err error)
	UpdatePost(ctx context.Context, post *domain.Post, attachmentsToDelete []string) (updatedPost *domain.Post, err error)
	DeletePost(ctx context.Context, postID uint) (err error)
	GetLikedPosts(ctx context.Context, userID uint, lastLikeID uint, limit uint) (likedPosts []LikeWithPost, err error)
	StorePostLike(ctx context.Context, likeData *domain.PostLike) (like *domain.PostLike, err error)
	DeletePostLike(ctx context.Context, likeData *domain.PostLike) (err error)
	StoreGroupPost(ctx context.Context, groupPost *domain.GroupPost) (newGroupPost *domain.GroupPost, err error)
	GetGroupPostByPostID(ctx context.Context, postID uint) (groupPost *domain.GroupPost, err error)
	DeleteGroupPost(ctx context.Context, postID uint) (err error)
	GetPostsOfGroup(ctx context.Context, groupID, lastPostID, postsAmount uint) (posts []*domain.Post, err error)
	GetGroupPostsBySubscriptionIDs(ctx context.Context, subIDs []uint, lastPostID, postsAmount uint) (posts []*domain.Post, err error)
	GetPostsByGroupSubIDsAndUserSubIDs(ctx context.Context, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) (posts []*domain.Post, err error)
	GetNewPosts(ctx context.Context, lastPostID, postsAmount uint) (posts []*domain.Post, err error)
	GetCommentsByPostID(ctx context.Context, postID uint) (comments []*domain.Comment, err error)
	GetCommentByID(ctx context.Context, id uint) (comment *domain.Comment, err error)
	StoreComment(ctx context.Context, comment *domain.Comment) (newComment *domain.Comment, err error)
	UpdateComment(ctx context.Context, comment *domain.Comment) (updatedComment *domain.Comment, err error)
	DeleteComment(ctx context.Context, id uint) (err error)
	GetCommentLikeByCommentIDAndUserID(ctx context.Context, data *domain.CommentLike) (commentLike *domain.CommentLike, err error)
	StoreCommentLike(ctx context.Context, commentLike *domain.CommentLike) (newCommentLike *domain.CommentLike, err error)
	DeleteCommentLike(ctx context.Context, commentLike *domain.CommentLike) (err error)
}

type AttachmentStorage interface {
	Store(fileName string, filePath string, contentType string) (err error)
	Delete(fileName string) (err error)
}

type Service struct {
	PostsStorage      PostsStorage
	AttachmentStorage AttachmentStorage
	Sanitizer         *sanitizer.Sanitizer
}

//easyjson:json
type ListPostsResponse struct {
	Posts []*domain.Post `json:"posts"`
}

func NewPostsService(postsStorage PostsStorage, attachmentStorage AttachmentStorage) (postsService *Service) {
	postsService = &Service{
		PostsStorage:      postsStorage,
		AttachmentStorage: attachmentStorage,
		Sanitizer:         sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
	}

	return
}

func (s *Service) GetPostByID(ctx context.Context, postID uint) (post *domain.Post, err error) {
	post, err = s.PostsStorage.GetPostByID(ctx, postID)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizePost(post)

	return
}

func (s *Service) GetUserPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error) {
	if postsAmount == 0 {
		postsAmount = DefaultPostsAmount
	}

	posts, err = s.PostsStorage.GetUserPosts(ctx, userID, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePost(post)
	}

	return
}

func (s *Service) GetUserFriendsPosts(ctx context.Context, userID uint, lastPostID uint, postsAmount uint) (posts []*domain.Post, err error) {
	if postsAmount == 0 {
		postsAmount = DefaultPostsAmount
	}

	posts, err = s.PostsStorage.GetUserFriendsPosts(ctx, userID, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePost(post)
	}

	return
}

func (s *Service) CreatePost(ctx context.Context, input PostInput) (newPost *domain.Post, err error) {
	if len(input.Content) == 0 && len(input.Attachments) == 0 {
		err = errors.ErrInvalidBody
		return
	}

	newPost, err = s.PostsStorage.StorePost(ctx, &domain.Post{AuthorID: input.AuthorID, Content: input.Content, Attachments: input.Attachments})
	if err != nil {
		return
	}

	s.Sanitizer.SanitizePost(newPost)

	return
}

func (s *Service) UpdatePost(ctx context.Context, userID uint, input PostUpdateInput) (post *domain.Post, err error) {
	oldPost, err := s.PostsStorage.GetPostByID(ctx, input.PostID)
	if err != nil {
		return
	}

	if oldPost.AuthorID != userID {
		err = errors.ErrForbidden
		return
	}

	if len(input.Content) == 0 && (len(oldPost.Attachments)+len(input.AttachmentsToAdd)-len(input.AttachmentsToDelete)) == 0 {
		err = errors.ErrInvalidBody
		return
	}

	oldPost.Content = input.Content
	oldPost.Attachments = input.AttachmentsToAdd

	_, err = s.PostsStorage.UpdatePost(ctx, oldPost, input.AttachmentsToDelete)
	if err != nil {
		return
	}

	for _, attachment := range input.AttachmentsToDelete {
		err = s.AttachmentStorage.Delete(attachment)
		if err != nil {
			return
		}
	}

	post, err = s.PostsStorage.GetPostByID(ctx, input.PostID)
	if err != nil {
		return
	}

	return
}

func (s *Service) DeletePost(ctx context.Context, userID uint, postID uint) (err error) {
	post, err := s.PostsStorage.GetPostByID(ctx, postID)
	if err != nil {
		return
	}

	if post.AuthorID != userID {
		err = errors.ErrForbidden
		return
	}

	for _, attachment := range post.Attachments {
		err = s.AttachmentStorage.Delete(attachment)
		if err != nil {
			return
		}
	}

	err = s.PostsStorage.DeleteGroupPost(ctx, postID)
	if err != nil {
		return
	}

	err = s.PostsStorage.DeletePost(ctx, postID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetLikedPosts(ctx context.Context, userID uint, lastLikeID uint, limit uint) (likedPosts []LikeWithPost, err error) {
	if limit == 0 {
		limit = DefaultLikedPostsAmount
	}

	likedPosts, err = s.PostsStorage.GetLikedPosts(ctx, userID, lastLikeID, limit)
	if err != nil {
		return
	}

	for _, likedPost := range likedPosts {
		s.Sanitizer.SanitizePost(likedPost.Post)
	}

	return
}

func (s *Service) LikePost(ctx context.Context, likeData *domain.PostLike) (like *domain.PostLike, err error) {
	like, err = s.PostsStorage.StorePostLike(ctx, likeData)
	if err != nil {
		return
	}

	return
}

func (s *Service) UnlikePost(ctx context.Context, likeData *domain.PostLike) (err error) {
	err = s.PostsStorage.DeletePostLike(ctx, likeData)
	if err != nil {
		return
	}

	return
}

func (s *Service) UploadAttachment(fileName string, filePath string, contentType string) (err error) {
	err = s.AttachmentStorage.Store(fileName, filePath, contentType)
	if err != nil {
		return
	}

	return
}

func (s *Service) CreateGroupPost(ctx context.Context, groupPost *domain.GroupPost) (newGroupPost *domain.GroupPost, err error) {
	newGroupPost, err = s.PostsStorage.StoreGroupPost(ctx, groupPost)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetGroupPostByPostID(ctx context.Context, postID uint) (groupPost *domain.GroupPost, err error) {
	_, err = s.PostsStorage.GetPostByID(ctx, postID)
	if err != nil {
		return
	}

	groupPost, err = s.PostsStorage.GetGroupPostByPostID(ctx, postID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetPostsOfGroup(ctx context.Context, groupID, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if postsAmount == 0 {
		postsAmount = DefaultPostsAmount
	}

	posts, err = s.PostsStorage.GetPostsOfGroup(ctx, groupID, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePost(post)
	}

	return
}

func (s *Service) GetGroupPostsBySubscriptionIDs(ctx context.Context, subIDs []uint, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if postsAmount == 0 {
		postsAmount = DefaultPostsAmount
	}

	posts, err = s.PostsStorage.GetGroupPostsBySubscriptionIDs(ctx, subIDs, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePost(post)
	}

	return
}

func (s *Service) GetPostsByGroupSubIDsAndUserSubIDs(ctx context.Context, groupSubIDs, userSubIDs []uint, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if postsAmount == 0 {
		postsAmount = DefaultPostsAmount
	}

	posts, err = s.PostsStorage.GetPostsByGroupSubIDsAndUserSubIDs(ctx, groupSubIDs, userSubIDs, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePost(post)
	}

	return
}

func (s *Service) GetNewPosts(ctx context.Context, lastPostID, postsAmount uint) (posts []*domain.Post, err error) {
	if postsAmount == 0 {
		postsAmount = DefaultPostsAmount
	}

	posts, err = s.PostsStorage.GetNewPosts(ctx, lastPostID, postsAmount)
	if err != nil {
		return
	}

	for _, post := range posts {
		s.Sanitizer.SanitizePost(post)
	}

	return
}
