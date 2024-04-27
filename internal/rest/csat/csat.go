package rest

import (
	defJSON "encoding/json"
	"net/http"
	"socio/errors"
	csatpb "socio/internal/grpc/csat/proto"
	"socio/pkg/json"
)

type CreateQuestionInput struct {
	PoolID    uint   `json:"poolId"`
	Question  string `json:"question"`
	BestCase  string `json:"bestCase"`
	WorstCase string `json:"worstCase"`
}

type UpdateQuestionInput struct {
	ID        uint   `json:"id"`
	Question  string `json:"question"`
	BestCase  string `json:"bestCase"`
	WorstCase string `json:"worstCase"`
}

type DeleteQuestionInput struct {
	ID uint `json:"id"`
}

type CreatePoolInput struct {
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type UpdatePoolInput struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type DeletePoolInput struct {
	ID uint `json:"id"`
}

type CSATHandler struct {
	CSATClient csatpb.CSATClient
}

func NewCSATHandler(csatClient csatpb.CSATClient) (c *CSATHandler) {
	return &CSATHandler{
		CSATClient: csatClient,
	}
}

func (c *CSATHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input CreateQuestionInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	res, err := c.CSATClient.CreateQuestion(r.Context(), &csatpb.CreateQuestionRequest{
		PoolId:    uint64(input.PoolID),
		Question:  input.Question,
		BestCase:  input.BestCase,
		WorstCase: input.WorstCase,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, csatpb.ToCSATQuestion(res.Question), http.StatusCreated)
}

func (c *CSATHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input UpdateQuestionInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	res, err := c.CSATClient.UpdateQuestion(r.Context(), &csatpb.UpdateQuestionRequest{
		Id:        uint64(input.ID),
		Question:  input.Question,
		BestCase:  input.BestCase,
		WorstCase: input.WorstCase,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, csatpb.ToCSATQuestion(res.Question), http.StatusOK)
}

func (c *CSATHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input DeleteQuestionInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	_, err = c.CSATClient.DeleteQuestion(r.Context(), &csatpb.DeleteQuestionRequest{
		Id: uint64(input.ID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CSATHandler) CreatePool(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input CreatePoolInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	res, err := c.CSATClient.CreatePool(r.Context(), &csatpb.CreatePoolRequest{
		Name:     input.Name,
		IsActive: input.IsActive,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, csatpb.ToCSATPool(res.Pool), http.StatusCreated)
}

func (c *CSATHandler) UpdatePool(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input UpdatePoolInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	res, err := c.CSATClient.UpdatePool(r.Context(), &csatpb.UpdatePoolRequest{
		Id:       uint64(input.ID),
		Name:     input.Name,
		IsActive: input.IsActive,
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	json.ServeJSONBody(r.Context(), w, csatpb.ToCSATPool(res.Pool), http.StatusOK)
}

func (c *CSATHandler) DeletePool(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input DeletePoolInput

	decoder := defJSON.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		json.ServeJSONError(r.Context(), w, errors.ErrJSONUnmarshalling)
		return
	}

	_, err = c.CSATClient.DeletePool(r.Context(), &csatpb.DeletePoolRequest{
		Id: uint64(input.ID),
	})
	if err != nil {
		json.ServeGRPCStatus(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
