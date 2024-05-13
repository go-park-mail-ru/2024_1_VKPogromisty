package post

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"socio/domain"
	"socio/errors"
	postspb "socio/internal/grpc/post/proto"
	"socio/pkg/utils"

	"socio/usecase/posts"

	"github.com/google/uuid"
)

const (
	staticFilePath = "."
)

type PostManager struct {
	postspb.UnimplementedPostServer

	PostsService *posts.Service
}

func NewPostManager(postsStorage posts.PostsStorage, attachmentStorage posts.AttachmentStorage) *PostManager {
	return &PostManager{
		PostsService: posts.NewPostsService(postsStorage, attachmentStorage),
	}
}

func (p *PostManager) GetPostByID(ctx context.Context, in *postspb.GetPostByIDRequest) (res *postspb.GetPostByIDResponse, err error) {
	postID := in.GetPostId()

	post, err := p.PostsService.GetPostByID(ctx, uint(postID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetPostByIDResponse{
		Post: postspb.ToPostResponse(post),
	}

	return
}

func (p *PostManager) GetUserPosts(ctx context.Context, in *postspb.GetUserPostsRequest) (res *postspb.GetUserPostsResponse, err error) {
	userID := in.GetUserId()
	lastPostID := in.GetLastPostId()
	postsAmount := in.GetPostsAmount()

	posts, err := p.PostsService.GetUserPosts(ctx, uint(userID), uint(lastPostID), uint(postsAmount))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetUserPostsResponse{
		Posts: postspb.ToPostsResponse(posts),
	}

	return
}

func (p *PostManager) GetUserFriendsPosts(ctx context.Context, in *postspb.GetUserFriendsPostsRequest) (res *postspb.GetUserFriendsPostsResponse, err error) {
	userID := in.GetUserId()
	lastPostID := in.GetLastPostId()
	postsAmount := in.GetPostsAmount()

	posts, err := p.PostsService.GetUserFriendsPosts(ctx, uint(userID), uint(lastPostID), uint(postsAmount))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetUserFriendsPostsResponse{
		Posts: postspb.ToPostsResponse(posts),
	}

	return
}

func (p *PostManager) CreatePost(ctx context.Context, in *postspb.CreatePostRequest) (res *postspb.CreatePostResponse, err error) {
	authorID := in.GetAuthorId()
	content := in.GetContent()
	attachments := in.GetAttachments()

	post, err := p.PostsService.CreatePost(ctx, posts.PostInput{
		AuthorID:    uint(authorID),
		Content:     content,
		Attachments: attachments,
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.CreatePostResponse{
		Post: postspb.ToPostResponse(post),
	}

	return
}

func (p *PostManager) UpdatePost(ctx context.Context, in *postspb.UpdatePostRequest) (res *postspb.UpdatePostResponse, err error) {
	postID := in.GetPostId()
	content := in.GetContent()
	userID := in.GetUserId()

	post, err := p.PostsService.UpdatePost(ctx, uint(userID), posts.PostUpdateInput{
		PostID:  uint(postID),
		Content: content,
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.UpdatePostResponse{
		Post: postspb.ToPostResponse(post),
	}

	return
}

func (p *PostManager) DeletePost(ctx context.Context, in *postspb.DeletePostRequest) (res *postspb.DeletePostResponse, err error) {
	postID := in.GetPostId()
	userID := in.GetUserId()

	err = p.PostsService.DeletePost(ctx, uint(userID), uint(postID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.DeletePostResponse{}

	return
}

func (p *PostManager) GetLikedPosts(ctx context.Context, in *postspb.GetLikedPostsRequest) (res *postspb.GetLikedPostsResponse, err error) {
	userID := in.GetUserId()
	lastLikeID := in.GetLastLikeId()
	postsAmount := in.GetPostsAmount()

	posts, err := p.PostsService.GetLikedPosts(ctx, uint(userID), uint(lastLikeID), uint(postsAmount))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetLikedPostsResponse{
		LikedPosts: postspb.ToLikedPosts(posts),
	}

	return
}

func (p *PostManager) LikePost(ctx context.Context, in *postspb.LikePostRequest) (res *postspb.LikePostResponse, err error) {
	postID := in.GetPostId()
	userID := in.GetUserId()

	like, err := p.PostsService.LikePost(ctx, &domain.PostLike{
		PostID: uint(postID),
		UserID: uint(userID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.LikePostResponse{
		Like: postspb.ToPostLikeResponse(like),
	}

	return
}

func (p *PostManager) UnlikePost(ctx context.Context, in *postspb.UnlikePostRequest) (res *postspb.UnlikePostResponse, err error) {
	postID := in.GetPostId()
	userID := in.GetUserId()

	err = p.PostsService.UnlikePost(ctx, &domain.PostLike{
		PostID: uint(postID),
		UserID: uint(userID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.UnlikePostResponse{}

	return
}

func (p *PostManager) Upload(stream postspb.Post_UploadServer) (err error) {
	file, err := os.Create(filepath.Join(staticFilePath, uuid.NewString()))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	fileName := ""
	contentType := ""

	var fileSize uint64
	fileSize = 0
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
			customErr := errors.NewCustomError(err)
			err = customErr.GRPCStatus().Err()
		}
	}()
	for {
		req, err := stream.Recv()
		if fileName == "" {
			fileName = req.GetFileName()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			customErr := errors.NewCustomError(err)
			err = customErr.GRPCStatus().Err()
			return err
		}
		chunk := req.GetChunk()
		fileSize += uint64(len(chunk))
		if _, err = file.Write(chunk); err != nil {
			customErr := errors.NewCustomError(err)
			err = customErr.GRPCStatus().Err()
			return err
		}
		contentType = req.GetContentType()
	}

	p.PostsService.UploadAttachment(fileName, file.Name(), contentType)

	if err = os.Remove(file.Name()); err != nil {
		return
	}

	return stream.SendAndClose(&postspb.UploadResponse{
		FileName: fileName,
		Size:     fileSize,
	})
}

func (p *PostManager) CreateGroupPost(ctx context.Context, in *postspb.CreateGroupPostRequest) (res *postspb.CreateGroupPostResponse, err error) {
	postID := in.GetPostId()
	groupID := in.GetGroupId()

	_, err = p.PostsService.CreateGroupPost(ctx, &domain.GroupPost{
		PostID:  uint(postID),
		GroupID: uint(groupID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.CreateGroupPostResponse{}

	return
}

func (p *PostManager) GetPostsOfGroup(ctx context.Context, in *postspb.GetPostsOfGroupRequest) (res *postspb.GetPostsOfGroupResponse, err error) {
	groupID := in.GetGroupId()
	lastPostID := in.GetLastPostId()
	postsAmount := in.GetPostsAmount()

	posts, err := p.PostsService.GetPostsOfGroup(ctx, uint(groupID), uint(lastPostID), uint(postsAmount))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetPostsOfGroupResponse{
		Posts: postspb.ToPostsResponse(posts),
	}

	return
}

func (p *PostManager) GetGroupPostsBySubscriptionIDs(ctx context.Context, in *postspb.GetGroupPostsBySubscriptionIDsRequest) (res *postspb.GetGroupPostsBySubscriptionIDsResponse, err error) {
	subIDs := in.GetSubscriptionIds()
	lastPostID := in.GetLastPostId()
	postsAmount := in.GetPostsAmount()

	subIDsUint := utils.Uint64ToUintSlice(subIDs)

	posts, err := p.PostsService.GetGroupPostsBySubscriptionIDs(ctx, subIDsUint, uint(lastPostID), uint(postsAmount))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetGroupPostsBySubscriptionIDsResponse{
		Posts: postspb.ToPostsResponse(posts),
	}

	return
}

func (p *PostManager) GetPostsByGroupSubIDsAndUserSubIDs(ctx context.Context, in *postspb.GetPostsByGroupSubIDsAndUserSubIDsRequest) (res *postspb.GetPostsByGroupSubIDsAndUserSubIDsResponse, err error) {
	groupSubIDs := in.GetGroupSubscriptionIds()
	userSubIDs := in.GetUserSubscriptionIds()
	lastPostID := in.GetLastPostId()
	postsAmount := in.GetPostsAmount()

	groupSubIDsUint := utils.Uint64ToUintSlice(groupSubIDs)
	userSubIDsUint := utils.Uint64ToUintSlice(userSubIDs)

	posts, err := p.PostsService.GetPostsByGroupSubIDsAndUserSubIDs(ctx, groupSubIDsUint, userSubIDsUint, uint(lastPostID), uint(postsAmount))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetPostsByGroupSubIDsAndUserSubIDsResponse{
		Posts: postspb.ToPostsResponse(posts),
	}

	return
}

func (p *PostManager) GetNewPosts(ctx context.Context, in *postspb.GetNewPostsRequest) (res *postspb.GetNewPostsResponse, err error) {
	lastPostID := in.GetLastPostId()
	postsAmount := in.GetPostsAmount()

	posts, err := p.PostsService.GetNewPosts(ctx, uint(lastPostID), uint(postsAmount))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetNewPostsResponse{
		Posts: postspb.ToPostsResponse(posts),
	}

	return
}

func (p *PostManager) GetCommentsByPostID(ctx context.Context, in *postspb.GetCommentsByPostIDRequest) (res *postspb.GetCommentsByPostIDResponse, err error) {
	postID := in.GetPostId()

	comments, err := p.PostsService.GetCommentsByPostID(ctx, uint(postID))
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.GetCommentsByPostIDResponse{
		Comments: postspb.ToCommentsResponse(comments),
	}

	return
}

func (p *PostManager) CreateComment(ctx context.Context, in *postspb.CreateCommentRequest) (res *postspb.CreateCommentResponse, err error) {
	postID := in.GetPostId()
	authorID := in.GetAuthorId()
	content := in.GetContent()

	comment, err := p.PostsService.CreateComment(ctx, &domain.Comment{
		PostID:   uint(postID),
		AuthorID: uint(authorID),
		Content:  content,
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.CreateCommentResponse{
		Comment: postspb.ToCommentResponse(comment),
	}

	return
}

func (p *PostManager) UpdateComment(ctx context.Context, in *postspb.UpdateCommentRequest) (res *postspb.UpdateCommentResponse, err error) {
	commentID := in.GetCommentId()
	content := in.GetContent()
	userID := in.GetUserId()

	comment, err := p.PostsService.UpdateComment(ctx, &domain.Comment{
		ID:       uint(commentID),
		AuthorID: uint(userID),
		Content:  content,
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.UpdateCommentResponse{
		Comment: postspb.ToCommentResponse(comment),
	}

	return
}

func (p *PostManager) DeleteComment(ctx context.Context, in *postspb.DeleteCommentRequest) (res *postspb.DeleteCommentResponse, err error) {
	commentID := in.GetCommentId()
	userID := in.GetUserId()

	err = p.PostsService.DeleteComment(ctx, &domain.Comment{
		ID:       uint(commentID),
		AuthorID: uint(userID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.DeleteCommentResponse{}

	return
}

func (p *PostManager) LikeComment(ctx context.Context, in *postspb.LikeCommentRequest) (res *postspb.LikeCommentResponse, err error) {
	commentID := in.GetCommentId()
	userID := in.GetUserId()

	like, err := p.PostsService.LikeComment(ctx, &domain.CommentLike{
		CommentID: uint(commentID),
		UserID:    uint(userID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.LikeCommentResponse{
		Like: postspb.ToCommentLikeResponse(like),
	}

	return
}

func (p *PostManager) UnlikeComment(ctx context.Context, in *postspb.UnlikeCommentRequest) (res *postspb.UnlikeCommentResponse, err error) {
	commentID := in.GetCommentId()
	userID := in.GetUserId()

	err = p.PostsService.UnlikeComment(ctx, &domain.CommentLike{
		CommentID: uint(commentID),
		UserID:    uint(userID),
	})
	if err != nil {
		customErr := errors.NewCustomError(err)
		err = customErr.GRPCStatus().Err()
		return
	}

	res = &postspb.UnlikeCommentResponse{}

	return
}
