package csat

import (
	"socio/domain"
	customtime "socio/pkg/time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCSATQuestionResponse(question *domain.CSATQuestion) (res *QuestionResponse) {
	return &QuestionResponse{
		Id:        uint64(question.ID),
		PoolId:    uint64(question.PoolID),
		Question:  question.Question,
		WorstCase: question.WorstCase,
		BestCase:  question.BestCase,
		CreatedAt: timestamppb.New(question.CreatedAt.Time),
		UpdatedAt: timestamppb.New(question.UpdatedAt.Time),
	}
}

func ToCSATPoolResponse(pool *domain.CSATPool) (res *PoolResponse) {
	return &PoolResponse{
		Id:        uint64(pool.ID),
		Name:      pool.Name,
		IsActive:  pool.IsActive,
		CreatedAt: timestamppb.New(pool.CreatedAt.Time),
		UpdatedAt: timestamppb.New(pool.UpdatedAt.Time),
	}
}

func ToCSATQuestion(question *QuestionResponse) (res *domain.CSATQuestion) {
	return &domain.CSATQuestion{
		PoolID:    uint(question.PoolId),
		Question:  question.Question,
		BestCase:  question.BestCase,
		WorstCase: question.WorstCase,
		CreatedAt: customtime.CustomTime{Time: question.CreatedAt.AsTime()},
		UpdatedAt: customtime.CustomTime{Time: question.UpdatedAt.AsTime()},
	}
}

func ToCSATPool(pool *PoolResponse) (res *domain.CSATPool) {
	return &domain.CSATPool{
		ID:        uint(pool.Id),
		Name:      pool.Name,
		IsActive:  pool.IsActive,
		CreatedAt: customtime.CustomTime{Time: pool.CreatedAt.AsTime()},
		UpdatedAt: customtime.CustomTime{Time: pool.UpdatedAt.AsTime()},
	}
}

func ToCSATReplyResponse(reply *domain.CSATReply) (res *ReplyResponse) {
	return &ReplyResponse{
		Id:         uint64(reply.ID),
		QuestionId: uint64(reply.QuestionID),
		UserId:     uint64(reply.UserID),
		Score:      int64(reply.Score),
		CreatedAt:  timestamppb.New(reply.CreatedAt.Time),
		UpdatedAt:  timestamppb.New(reply.UpdatedAt.Time),
	}
}

func ToCSATStatsResponse(stats *domain.CSATStat) (res *StatsResponse) {
	return &StatsResponse{
		Question:     ToCSATQuestionResponse(stats.Question),
		TotalReplies: uint64(stats.TotalReplies),
		AvgScore:     float32(stats.AvgScore),
	}
}

func ToCSATPools(res []*PoolResponse) (pools []*domain.CSATPool) {
	for _, pool := range pools {
		res = append(res, ToCSATPoolResponse(pool))
	}
	return
}

func ToCSATQuestions(res []*QuestionResponse) (questions []*domain.CSATQuestion) {
	for _, question := range res {
		questions = append(questions, ToCSATQuestion(question))
	}
	return
}

func ToCSATReply(res *ReplyResponse) (reply *domain.CSATReply) {
	return &domain.CSATReply{
		QuestionID: uint(res.QuestionId),
		UserID:     uint(res.UserId),
		Score:      int(res.Score),
		CreatedAt:  customtime.CustomTime{Time: res.CreatedAt.AsTime()},
		UpdatedAt:  customtime.CustomTime{Time: res.UpdatedAt.AsTime()},
	}
}

func ToCSATStats(res []*StatsResponse) (stats []*domain.CSATStat) {
	for _, stat := range res {
		stats = append(stats, &domain.CSATStat{
			Question:     ToCSATQuestion(stat.Question),
			TotalReplies: int(stat.TotalReplies),
			AvgScore:     float64(stat.AvgScore),
		})
	}
	return
}
