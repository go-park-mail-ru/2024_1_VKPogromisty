package domain

import customtime "socio/pkg/time"

type Admin struct {
	ID        uint                  `json:"id"`
	UserID    uint                  `json:"userId"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type CSATPool struct {
	ID        uint                  `json:"id"`
	Name      string                `json:"name"`
	IsActive  bool                  `json:"isActive"`
	Questions []*CSATQuestion       `json:"questions"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type CSATQuestion struct {
	ID        uint                  `json:"id"`
	PoolID    uint                  `json:"poolId"`
	Question  string                `json:"question"`
	WorstCase string                `json:"worstCase"`
	BestCase  string                `json:"bestCase"`
	CreatedAt customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}

type CSATReply struct {
	ID         uint
	QuestionID uint
	Score      int
	CreatedAt  customtime.CustomTime `json:"createdAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
	UpdatedAt  customtime.CustomTime `json:"updatedAt,omitempty" swaggertype:"string" example:"2021-01-01T00:00:00Z" format:"date-time"`
}
