package csat

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/sanitizer"

	"github.com/microcosm-cc/bluemonday"
)

type CSATStorage interface {
	GetQuestionByID(ctx context.Context, questionID uint) (question *domain.CSATQuestion, err error)
	StoreQuestion(ctx context.Context, question *domain.CSATQuestion) (newQuestion *domain.CSATQuestion, err error)
	UpdateQuestion(ctx context.Context, question *domain.CSATQuestion) (updatedQuestion *domain.CSATQuestion, err error)
	DeleteQuestion(ctx context.Context, questionID uint) (err error)
	CreatePool(ctx context.Context, pool *domain.CSATPool) (newPool *domain.CSATPool, err error)
	UpdatePool(ctx context.Context, pool *domain.CSATPool) (updatedPool *domain.CSATPool, err error)
	DeletePool(ctx context.Context, poolID uint) (err error)
	GetPoolByID(ctx context.Context, poolID uint) (pool *domain.CSATPool, err error)
	GetPools(ctx context.Context) (pools []*domain.CSATPool, err error)
	GetQuestionsByPoolID(ctx context.Context, poolID uint) (questions []*domain.CSATQuestion, err error)
	GetUnansweredQuestionsByPoolID(ctx context.Context, userID uint, poolID uint) (questions []*domain.CSATQuestion, err error)
	CreateReply(ctx context.Context, reply *domain.CSATReply) (newReply *domain.CSATReply, err error)
	GetStatsByPool(ctx context.Context, poolID uint) (stats []*domain.CSATStat, err error)
}

type Service struct {
	CSATStorage CSATStorage
	Sanitizer   *sanitizer.Sanitizer
}

func NewService(storage CSATStorage) (s *Service) {
	return &Service{
		CSATStorage: storage,
		Sanitizer:   sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
	}
}

func (s *Service) CreateQuestion(ctx context.Context, question *domain.CSATQuestion) (newQuestion *domain.CSATQuestion, err error) {
	_, err = s.CSATStorage.GetPoolByID(ctx, question.PoolID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	newQuestion, err = s.CSATStorage.StoreQuestion(ctx, question)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizeCSATQuestion(newQuestion)

	return
}

func (s *Service) UpdateQuestion(ctx context.Context, question *domain.CSATQuestion) (updatedQuestion *domain.CSATQuestion, err error) {
	oldQuestion, err := s.CSATStorage.GetQuestionByID(ctx, question.ID)
	if err != nil {
		return
	}

	if len(question.Question) > 0 {
		oldQuestion.Question = question.Question
	}

	if len(question.BestCase) > 0 {
		oldQuestion.BestCase = question.BestCase
	}

	if len(question.WorstCase) > 0 {
		oldQuestion.WorstCase = question.WorstCase
	}

	if question.PoolID > 0 {
		_, err = s.CSATStorage.GetPoolByID(ctx, question.PoolID)
		if err != nil {
			return
		}

		oldQuestion.PoolID = question.PoolID
	}

	updatedQuestion, err = s.CSATStorage.UpdateQuestion(ctx, oldQuestion)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizeCSATQuestion(updatedQuestion)

	return
}

func (s *Service) DeleteQuestion(ctx context.Context, questionID uint) (err error) {
	_, err = s.CSATStorage.GetQuestionByID(ctx, questionID)
	if err != nil {
		return
	}

	err = s.CSATStorage.DeleteQuestion(ctx, questionID)
	return
}

func (s *Service) CreatePool(ctx context.Context, pool *domain.CSATPool) (newPool *domain.CSATPool, err error) {
	newPool, err = s.CSATStorage.CreatePool(ctx, pool)
	if err != nil {
		return
	}

	return
}

func (s *Service) UpdatePool(ctx context.Context, pool *domain.CSATPool) (updatedPool *domain.CSATPool, err error) {
	oldPool, err := s.CSATStorage.GetPoolByID(ctx, pool.ID)
	if err != nil {
		return
	}

	if len(pool.Name) > 0 {
		oldPool.Name = pool.Name
	}

	if pool.IsActive != oldPool.IsActive {
		oldPool.IsActive = pool.IsActive
	}

	updatedPool, err = s.CSATStorage.UpdatePool(ctx, oldPool)
	if err != nil {
		return
	}

	return
}

func (s *Service) DeletePool(ctx context.Context, poolID uint) (err error) {
	_, err = s.CSATStorage.GetPoolByID(ctx, poolID)
	if err != nil {
		return
	}

	err = s.CSATStorage.DeletePool(ctx, poolID)
	return
}

// here
func (s *Service) GetPools(ctx context.Context) (pools []*domain.CSATPool, err error) {
	pools, err = s.CSATStorage.GetPools(ctx)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetQuestionsByPoolID(ctx context.Context, poolID uint) (questions []*domain.CSATQuestion, err error) {
	questions, err = s.CSATStorage.GetQuestionsByPoolID(ctx, poolID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetUnansweredQuestionsByPoolID(ctx context.Context, userID uint, poolID uint) (pool *domain.CSATPool, questions []*domain.CSATQuestion, err error) {
	pool, err = s.CSATStorage.GetPoolByID(ctx, poolID)
	if err != nil {
		return
	}

	questions, err = s.CSATStorage.GetUnansweredQuestionsByPoolID(ctx, userID, poolID)
	if err != nil {
		return
	}

	return
}

func (s *Service) CreateReply(ctx context.Context, reply *domain.CSATReply) (newReply *domain.CSATReply, err error) {
	_, err = s.CSATStorage.GetQuestionByID(ctx, reply.QuestionID)
	if err != nil {
		err = errors.ErrInvalidBody
		return
	}

	newReply, err = s.CSATStorage.CreateReply(ctx, reply)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetStatsByPool(ctx context.Context, poolID uint) (stats []*domain.CSATStat, err error) {
	_, err = s.CSATStorage.GetPoolByID(ctx, poolID)
	if err != nil {
		return
	}

	stats, err = s.CSATStorage.GetStatsByPool(ctx, poolID)
	if err != nil {
		return
	}

	return
}
