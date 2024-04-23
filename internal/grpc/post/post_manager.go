package post

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"socio/domain"
	postspb "socio/internal/grpc/post/proto"
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
		return
	}

	res = &postspb.UpdatePostResponse{
		Post: postspb.ToPostResponse(post),
	}

	return
}

func (p *PostManager) DeletePost(ctx context.Context, in *postspb.DeletePostRequest) (res *postspb.DeletePostResponse, err error) {
	postID := in.GetPostId()
	userId := in.GetUserId()

	err = p.PostsService.DeletePost(ctx, uint(userId), uint(postID))
	if err != nil {
		return
	}

	res = &postspb.DeletePostResponse{}

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
		return
	}

	res = &postspb.UnlikePostResponse{}

	return
}

func (p *PostManager) UploadAvatar(stream postspb.Post_UploadAttachmentServer) (err error) {
	file, err := os.Create(filepath.Join(staticFilePath, uuid.NewString()))
	if err != nil {
		return
	}

	fileName := ""

	var fileSize uint64
	fileSize = 0
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println(err)
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
			return err
		}
		chunk := req.GetData()
		fileSize += uint64(len(chunk))
		if _, err = file.Write(chunk); err != nil {
			return err
		}
	}

	p.PostsService.UploadAttachment(fileName, file.Name())

	if err = os.Remove(file.Name()); err != nil {
		return
	}

	return stream.SendAndClose(&postspb.UploadAttachmentResponse{
		FileName: fileName,
		Size:     fileSize,
	})
}
