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
		Content:     post.Content,
		Attachments: post.Attachments,
		LikedByIds:  post.LikedByIDs,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: post.CreatedAt.Unix(),
			Nanos:   int32(post.CreatedAt.Nanosecond()),
		},
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
		Id:     uint64(like.ID),
		PostId: uint64(like.PostID),
		UserId: uint64(like.UserID),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: like.CreatedAt.Unix(),
			Nanos:   int32(like.CreatedAt.Nanosecond()),
		},
	}
}

func ToPost(res *PostResponse) *domain.Post {
	return &domain.Post{
		ID:          uint(res.Id),
		AuthorID:    uint(res.AuthorId),
		Content:     res.Content,
		Attachments: res.Attachments,
		LikedByIDs:  res.LikedByIds,
		CreatedAt: customtime.CustomTime{
			Time: res.CreatedAt.AsTime(),
		},
	}
}

func ToPosts(res *GetUserPostsResponse) (posts []*domain.Post) {
	posts = make([]*domain.Post, 0)

	for _, post := range res.Posts {
		posts = append(posts, &domain.Post{
			ID:          uint(post.Id),
			AuthorID:    uint(post.AuthorId),
			Content:     post.Content,
			Attachments: post.Attachments,
			LikedByIDs:  post.LikedByIds,
			CreatedAt: customtime.CustomTime{
				Time: post.CreatedAt.AsTime(),
			},
		})
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
