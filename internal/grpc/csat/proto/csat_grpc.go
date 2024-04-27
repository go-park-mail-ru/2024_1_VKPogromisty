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
