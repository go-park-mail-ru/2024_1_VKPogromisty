package csat

import (
	"context"
	"socio/domain"
	csatspb "socio/internal/grpc/csat/proto"
	"socio/usecase/csat"
)

type CSATManager struct {
	csatspb.UnimplementedCSATServer

	CSATService *csat.Service
}

func NewCSATManager(csatStorage csat.CSATStorage) *CSATManager {
	return &CSATManager{
		CSATService: csat.NewService(csatStorage),
	}
}

func (c *CSATManager) CreateQuestion(ctx context.Context, in *csatspb.CreateQuestionRequest) (res *csatspb.CreateQuestionResponse, err error) {
	poolID := in.GetPoolId()
	question := in.GetQuestion()
	worstCase := in.GetWorstCase()
	bestCase := in.GetBestCase()

	csatQuestion, err := c.CSATService.CreateQuestion(ctx, &domain.CSATQuestion{
		PoolID:    uint(poolID),
		Question:  question,
		WorstCase: worstCase,
		BestCase:  bestCase,
	})
	if err != nil {
		return
	}

	res = &csatspb.CreateQuestionResponse{
		Question: csatspb.ToCSATQuestionResponse(csatQuestion),
	}

	return
}

func (c *CSATManager) UpdateQuestion(ctx context.Context, in *csatspb.UpdateQuestionRequest) (res *csatspb.UpdateQuestionResponse, err error) {
	questionID := in.GetId()
	question := in.GetQuestion()
	worstCase := in.GetWorstCase()
	bestCase := in.GetBestCase()

	csatQuestion, err := c.CSATService.UpdateQuestion(ctx, &domain.CSATQuestion{
		ID:        uint(questionID),
		Question:  question,
		WorstCase: worstCase,
		BestCase:  bestCase,
	})
	if err != nil {
		return
	}

	res = &csatspb.UpdateQuestionResponse{
		Question: csatspb.ToCSATQuestionResponse(csatQuestion),
	}

	return
}

func (c *CSATManager) DeleteQuestion(ctx context.Context, in *csatspb.DeleteQuestionRequest) (res *csatspb.DeleteQuestionResponse, err error) {
	questionID := in.GetId()

	err = c.CSATService.DeleteQuestion(ctx, uint(questionID))
	if err != nil {
		return
	}

	res = &csatspb.DeleteQuestionResponse{}

	return
}

func (c *CSATManager) CreatePool(ctx context.Context, in *csatspb.CreatePoolRequest) (res *csatspb.CreatePoolResponse, err error) {
	poolName := in.GetName()
	isActive := in.GetIsActive()

	csatPool, err := c.CSATService.CreatePool(ctx, &domain.CSATPool{
		Name:     poolName,
		IsActive: isActive,
	})
	if err != nil {
		return
	}

	res = &csatspb.CreatePoolResponse{
		Pool: csatspb.ToCSATPoolResponse(csatPool),
	}

	return
}

func (c *CSATManager) UpdatePool(ctx context.Context, in *csatspb.UpdatePoolRequest) (res *csatspb.UpdatePoolResponse, err error) {
	poolID := in.GetId()
	poolName := in.GetName()
	isActive := in.GetIsActive()

	csatPool, err := c.CSATService.UpdatePool(ctx, &domain.CSATPool{
		ID:       uint(poolID),
		Name:     poolName,
		IsActive: isActive,
	})
	if err != nil {
		return
	}

	res = &csatspb.UpdatePoolResponse{
		Pool: csatspb.ToCSATPoolResponse(csatPool),
	}

	return
}

func (c *CSATManager) DeletePool(ctx context.Context, in *csatspb.DeletePoolRequest) (res *csatspb.DeletePoolResponse, err error) {
	poolID := in.GetId()

	err = c.CSATService.DeletePool(ctx, uint(poolID))
	if err != nil {
		return
	}

	res = &csatspb.DeletePoolResponse{}

	return
}

func (c *CSATManager) GetPools(ctx context.Context, in *csatspb.GetPoolsRequest) (res *csatspb.GetPoolsResponse, err error) {
	csatPools, err := c.CSATService.GetPools(ctx)
	if err != nil {
		return
	}

	pools := make([]*csatspb.PoolResponse, 0)
	for _, csatPool := range csatPools {
		pools = append(pools, csatspb.ToCSATPoolResponse(csatPool))
	}

	res = &csatspb.GetPoolsResponse{
		Pools: pools,
	}

	return
}

func (c *CSATManager) GetQuestionsByPoolID(ctx context.Context, in *csatspb.GetQuestionsByPoolIDRequest) (res *csatspb.GetQuestionsByPoolIDResponse, err error) {
	poolID := in.GetPoolId()

	csatQuestions, err := c.CSATService.GetQuestionsByPoolID(ctx, uint(poolID))
	if err != nil {
		return
	}

	questions := make([]*csatspb.QuestionResponse, 0)
	for _, csatQuestion := range csatQuestions {
		questions = append(questions, csatspb.ToCSATQuestionResponse(csatQuestion))
	}

	res = &csatspb.GetQuestionsByPoolIDResponse{
		Questions: questions,
	}

	return
}

func (c *CSATManager) GetUnansweredQuestionsByPoolID(ctx context.Context, in *csatspb.GetUnansweredQuestionsByPoolIDRequest) (res *csatspb.GetUnansweredQuestionsByPoolIDResponse, err error) {
	poolID := in.GetPoolId()
	userID := in.GetUserId()

	pool, csatQuestions, err := c.CSATService.GetUnansweredQuestionsByPoolID(ctx, uint(userID), uint(poolID))
	if err != nil {
		return
	}

	questions := make([]*csatspb.QuestionResponse, 0)
	for _, csatQuestion := range csatQuestions {
		questions = append(questions, csatspb.ToCSATQuestionResponse(csatQuestion))
	}

	res = &csatspb.GetUnansweredQuestionsByPoolIDResponse{
		Pool:      csatspb.ToCSATPoolResponse(pool),
		Questions: questions,
	}

	return
}

func (c *CSATManager) CreateReply(ctx context.Context, in *csatspb.CreateReplyRequest) (res *csatspb.CreateReplyResponse, err error) {
	questionID := in.GetQuestionId()
	userID := in.GetUserId()
	score := in.GetScore()

	csatReply, err := c.CSATService.CreateReply(ctx, &domain.CSATReply{
		QuestionID: uint(questionID),
		UserID:     uint(userID),
		Score:      int(score),
	})
	if err != nil {
		return
	}

	res = &csatspb.CreateReplyResponse{
		Reply: csatspb.ToCSATReplyResponse(csatReply),
	}

	return
}

func (c *CSATManager) GetStatsByPoolID(ctx context.Context, in *csatspb.GetStatsByPoolIDRequest) (res *csatspb.GetStatsByPoolIDResponse, err error) {
	poolID := in.GetPoolId()

	stats, err := c.CSATService.GetStatsByPool(ctx, uint(poolID))
	if err != nil {
		return
	}

	statsRes := make([]*csatspb.StatsResponse, 0)

	for _, stat := range stats {
		statsRes = append(statsRes, csatspb.ToCSATStatsResponse(stat))
	}

	res = &csatspb.GetStatsByPoolIDResponse{
		Stats: statsRes,
	}

	return
}
