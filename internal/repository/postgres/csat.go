package repository

import (
	"context"
	"socio/domain"
	"socio/errors"
	customtime "socio/pkg/time"

	"github.com/jackc/pgx/v4"
)

const (
	getQuestionByIDQuery = `
	SELECT id, pool_id, question, best_case, worst_case, created_at, updated_at
	FROM csat_question
	WHERE id = $1;
	`
	storeCSATQuestionQuery = `
	INSERT INTO csat_question (pool_id, question, best_case, worst_case)
	VALUES ($1, $2, $3, $4)
	RETURNING id, pool_id, question, best_case, worst_case, created_at, updated_at;
	`
	updateCSATQuestionQuery = `
	UPDATE csat_question
	SET pool_id = $1, question = $2, best_case = $3, worst_case = $4
	WHERE id = $5
	RETURNING id, pool_id, question, best_case, worst_case, created_at, updated_at;
	`
	deleteCSATQuestionQuery = `
	DELETE FROM csat_question
	WHERE id = $1;
	`
	storeCSATPoolQuery = `
	INSERT INTO csat_pool (name, is_active)
	VALUES ($1, $2)
	RETURNING id, name, is_active, created_at, updated_at;
	`
	updateCSATPoolQuery = `
	UPDATE csat_pool
	SET name = $1, is_active = $2
	WHERE id = $3
	RETURNING id, name, is_active, created_at, updated_at;
	`
	deleteCSATPoolQuery = `
	DELETE FROM csat_pool
	WHERE id = $1;
	`
	getActivePoolsQuery = `
	SELECT id, name, is_active, created_at, updated_at
	FROM csat_pool
	WHERE is_active = true;
	`
	getUnansweredQuestionsByPoolIDQuery = `
	SELECT q.id, q.pool_id, q.question, q.best_case, q.worst_case, q.created_at, q.updated_at
	FROM csat_question q
	LEFT JOIN csat_reply a ON q.id = a.question_id AND a.user_id = $2
	WHERE q.pool_id = $1 AND a.question_id IS NULL;
	`
	getPoolByIDQuery = `
	SELECT id, name, is_active, created_at, updated_at
	FROM csat_pool
	WHERE id = $1;
	`
)

type CSAT struct {
	db DBPool
	TP customtime.TimeProvider
}

func NewCSAT(db DBPool, tp customtime.TimeProvider) *CSAT {
	return &CSAT{
		db: db,
		TP: tp,
	}
}

func (s *CSAT) GetQuestionByID(ctx context.Context, questionID uint) (question *domain.CSATQuestion, err error) {
	question = new(domain.CSATQuestion)

	err = s.db.QueryRow(context.Background(), getQuestionByIDQuery, questionID).Scan(
		&question.ID,
		&question.PoolID,
		&question.Question,
		&question.BestCase,
		&question.WorstCase,
		&question.CreatedAt.Time,
		&question.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}

func (s *CSAT) StoreQuestion(ctx context.Context, question *domain.CSATQuestion) (newQuestion *domain.CSATQuestion, err error) {
	newQuestion = new(domain.CSATQuestion)

	err = s.db.QueryRow(context.Background(), storeCSATQuestionQuery,
		question.PoolID,
		question.Question,
		question.BestCase,
		question.WorstCase,
	).Scan(
		&newQuestion.ID,
		&newQuestion.PoolID,
		&newQuestion.Question,
		&newQuestion.BestCase,
		&newQuestion.WorstCase,
		&newQuestion.CreatedAt.Time,
		&newQuestion.UpdatedAt.Time,
	)

	if err != nil {
		return
	}

	return
}

func (s *CSAT) UpdateQuestion(ctx context.Context, question *domain.CSATQuestion) (updatedQuestion *domain.CSATQuestion, err error) {
	updatedQuestion = new(domain.CSATQuestion)

	err = s.db.QueryRow(context.Background(), updateCSATQuestionQuery,
		question.PoolID,
		question.Question,
		question.BestCase,
		question.WorstCase,
		question.ID,
	).Scan(
		&updatedQuestion.ID,
		&updatedQuestion.PoolID,
		&updatedQuestion.Question,
		&updatedQuestion.BestCase,
		&updatedQuestion.WorstCase,
		&updatedQuestion.CreatedAt.Time,
		&updatedQuestion.UpdatedAt.Time,
	)

	if err != nil {
		return
	}

	return
}

func (s *CSAT) DeleteQuestion(ctx context.Context, questionID uint) (err error) {
	_, err = s.db.Exec(context.Background(), deleteCSATQuestionQuery, questionID)
	if err != nil {
		return
	}

	return
}

func (s *CSAT) CreatePool(ctx context.Context, pool *domain.CSATPool) (newPool *domain.CSATPool, err error) {
	newPool = new(domain.CSATPool)

	err = s.db.QueryRow(context.Background(), storeCSATPoolQuery,
		pool.Name,
		pool.IsActive,
	).Scan(
		&newPool.ID,
		&newPool.Name,
		&newPool.IsActive,
		&newPool.CreatedAt.Time,
		&newPool.UpdatedAt.Time,
	)

	if err != nil {
		return
	}

	return
}

func (s *CSAT) UpdatePool(ctx context.Context, pool *domain.CSATPool) (updatedPool *domain.CSATPool, err error) {
	updatedPool = new(domain.CSATPool)

	err = s.db.QueryRow(context.Background(), updateCSATPoolQuery,
		pool.Name,
		pool.IsActive,
		pool.ID,
	).Scan(
		&updatedPool.ID,
		&updatedPool.Name,
		&updatedPool.IsActive,
		&updatedPool.CreatedAt.Time,
		&updatedPool.UpdatedAt.Time,
	)

	if err != nil {
		return
	}

	return
}

func (s *CSAT) DeletePool(ctx context.Context, poolID uint) (err error) {
	_, err = s.db.Exec(context.Background(), deleteCSATPoolQuery, poolID)
	if err != nil {
		return
	}

	return
}

func (s *CSAT) GetActivePools(ctx context.Context) (pools []*domain.CSATPool, err error) {
	rows, err := s.db.Query(context.Background(), getActivePoolsQuery)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		pool := new(domain.CSATPool)
		err = rows.Scan(
			&pool.ID,
			&pool.Name,
			&pool.IsActive,
			&pool.CreatedAt.Time,
			&pool.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		pools = append(pools, pool)
	}

	return
}

func (s *CSAT) GetPoolByID(ctx context.Context, poolID uint) (pool *domain.CSATPool, err error) {
	pool = new(domain.CSATPool)

	err = s.db.QueryRow(context.Background(), getPoolByIDQuery, poolID).Scan(
		&pool.ID,
		&pool.Name,
		&pool.IsActive,
		&pool.CreatedAt.Time,
		&pool.UpdatedAt.Time,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}

	return
}

func (s *CSAT) GetUnansweredQuestionsByPoolID(ctx context.Context, poolID uint) (questions []*domain.CSATQuestion, err error) {
	rows, err := s.db.Query(context.Background(), getUnansweredQuestionsByPoolIDQuery, poolID)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errors.ErrNotFound
			return
		}

		return
	}
	defer rows.Close()

	for rows.Next() {
		question := new(domain.CSATQuestion)
		err = rows.Scan(
			&question.ID,
			&question.PoolID,
			&question.Question,
			&question.BestCase,
			&question.WorstCase,
			&question.CreatedAt.Time,
			&question.UpdatedAt.Time,
		)
		if err != nil {
			return
		}

		questions = append(questions, question)
	}

	return
}
