package post

import (
	"socio/domain"
	customtime "socio/pkg/time"
	"socio/usecase/posts"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToPostResponse(post *domain.Post) *PostResponse {
	return &PostResponse{
		Id:          uint64(post.ID),
		AuthorId:    uint64(post.AuthorID),
		GroupId:     uint64(post.GroupID),
		Content:     post.Content,
		Attachments: post.Attachments,
		LikedByIds:  post.LikedByIDs,
		CreatedAt:   timestamppb.New(post.CreatedAt.Time),
		UpdatedAt:   timestamppb.New(post.UpdatedAt.Time),
	}
}

func ToPostsResponse(posts []*domain.Post) (res []*PostResponse) {
	res = make([]*PostResponse, 0)

	for _, post := range posts {
		res = append(res, ToPostResponse(post))
	}

	return
}

func ToPostLikeResponse(like *domain.PostLike) *PostLikeResponse {
	return &PostLikeResponse{
		Id:        uint64(like.ID),
		PostId:    uint64(like.PostID),
		UserId:    uint64(like.UserID),
		CreatedAt: timestamppb.New(like.CreatedAt.Time),
	}
}

func ToPost(res *PostResponse) *domain.Post {
	return &domain.Post{
		ID:          uint(res.Id),
		AuthorID:    uint(res.AuthorId),
		GroupID:     uint(res.GroupId),
		Content:     res.Content,
		Attachments: res.Attachments,
		LikedByIDs:  res.LikedByIds,
		CreatedAt: customtime.CustomTime{
			Time: res.CreatedAt.AsTime(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: res.UpdatedAt.AsTime(),
		},
	}
}

func ToPosts(res []*PostResponse) (posts []*domain.Post) {
	posts = make([]*domain.Post, 0)

	for _, post := range res {
		posts = append(posts, ToPost(post))
	}

	return
}

func ToPostLike(res *PostLikeResponse) *domain.PostLike {
	return &domain.PostLike{
		ID:     uint(res.Id),
		PostID: uint(res.PostId),
		UserID: uint(res.UserId),
		CreatedAt: customtime.CustomTime{
			Time: res.CreatedAt.AsTime(),
		},
	}
}

func ToLikeWithPost(res *LikedPostResponse) *posts.LikeWithPost {
	return &posts.LikeWithPost{
		Post: ToPost(res.Post),
		Like: ToPostLike(res.Like),
	}
}

func ToLikesWithPosts(res *GetLikedPostsResponse) (likesWithPosts []*posts.LikeWithPost) {
	likesWithPosts = make([]*posts.LikeWithPost, 0)

	for _, likedPost := range res.LikedPosts {
		likesWithPosts = append(likesWithPosts, ToLikeWithPost(likedPost))
	}

	return
}

func ToLikedPosts(likesWithPosts []posts.LikeWithPost) (res []*LikedPostResponse) {
	res = make([]*LikedPostResponse, 0)

	for _, likeWithPost := range likesWithPosts {
		res = append(res, &LikedPostResponse{
			Post: ToPostResponse(likeWithPost.Post),
			Like: ToPostLikeResponse(likeWithPost.Like),
		})
	}

	return
}

func ToCommentResponse(comment *domain.Comment) (res *CommentResponse) {
	if comment == nil {
		return nil
	}

	return &CommentResponse{
		Id:         uint64(comment.ID),
		AuthorId:   uint64(comment.AuthorID),
		PostId:     uint64(comment.PostID),
		Content:    comment.Content,
		CreatedAt:  timestamppb.New(comment.CreatedAt.Time),
		UpdatedAt:  timestamppb.New(comment.UpdatedAt.Time),
		LikedByIds: comment.LikedByIDs,
	}
}

func ToCommentsResponse(comments []*domain.Comment) (res []*CommentResponse) {
	res = make([]*CommentResponse, 0, len(comments))

	for _, comment := range comments {
		res = append(res, ToCommentResponse(comment))
	}

	return
}

func ToCommentLikeResponse(like *domain.CommentLike) (res *CommentLikeResponse) {
	if like == nil {
		return nil
	}

	return &CommentLikeResponse{
		Id:        uint64(like.ID),
		CommentId: uint64(like.CommentID),
		UserId:    uint64(like.UserID),
		CreatedAt: timestamppb.New(like.CreatedAt.Time),
	}
}

func ToComment(res *CommentResponse) *domain.Comment {
	if res == nil {
		return nil
	}

	return &domain.Comment{
		ID:         uint(res.Id),
		AuthorID:   uint(res.AuthorId),
		PostID:     uint(res.PostId),
		Content:    res.Content,
		CreatedAt:  customtime.CustomTime{Time: res.CreatedAt.AsTime()},
		UpdatedAt:  customtime.CustomTime{Time: res.UpdatedAt.AsTime()},
		LikedByIDs: res.LikedByIds,
	}
}

func ToComments(res []*CommentResponse) (comments []*domain.Comment) {
	comments = make([]*domain.Comment, 0, len(res))

	for _, comment := range res {
		comments = append(comments, ToComment(comment))
	}

	return
}

func ToCommentLike(res *CommentLikeResponse) *domain.CommentLike {
	if res == nil {
		return nil
	}

	return &domain.CommentLike{
		ID:        uint(res.Id),
		CommentID: uint(res.CommentId),
		UserID:    uint(res.UserId),
		CreatedAt: customtime.CustomTime{Time: res.CreatedAt.AsTime()},
	}
}

func ToGroupPostResponse(groupPost *domain.GroupPost) (res *GroupPostResponse) {
	if groupPost == nil {
		return nil
	}

	return &GroupPostResponse{
		Id:        uint64(groupPost.ID),
		PostId:    uint64(groupPost.PostID),
		GroupId:   uint64(groupPost.GroupID),
		CreatedAt: timestamppb.New(groupPost.CreatedAt.Time),
		UpdatedAt: timestamppb.New(groupPost.UpdatedAt.Time),
	}
}

func ToGroupPost(res *GroupPostResponse) (groupPost *domain.GroupPost) {
	if res == nil {
		return nil
	}

	return &domain.GroupPost{
		ID:      uint(res.GetId()),
		PostID:  uint(res.GetPostId()),
		GroupID: uint(res.GetGroupId()),
		CreatedAt: customtime.CustomTime{
			Time: res.CreatedAt.AsTime(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: res.UpdatedAt.AsTime(),
		},
	}
}
