package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"socio/errors"
	postpb "socio/internal/grpc/post/proto"
	pgpb "socio/internal/grpc/public_group/proto"
	uspb "socio/internal/grpc/user/proto"
	mock_post "socio/mocks/grpc/post_grpc"
	mock_public_group "socio/mocks/grpc/public_group_grpc"
	mock_user "socio/mocks/grpc/user_grpc"
	"socio/pkg/requestcontext"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful get group by ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusOK,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&pgpb.GetByIDResponse{
					PublicGroup: &pgpb.PublicGroupWithInfoResponse{
						PublicGroup: &pgpb.PublicGroupResponse{},
					},
				}, nil)
			},
		},
		{
			name:           "Successful get group by ID",
			ctx:            context.Background(),
			groupID:        "1",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get group by ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get group by ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get group by ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/groups/"+tt.groupID, nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			r = mux.SetURLVars(r, map[string]string{
				"groupID": tt.groupID,
			})

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleGetByID(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleSearchByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		query          string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful search by name",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			query:          "test",
			expectedStatus: http.StatusOK,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().SearchByName(gomock.Any(), gomock.Any()).Return(&pgpb.SearchByNameResponse{}, nil)
			},
		},
		{
			name:           "Successful search by name",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			query:          "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
			},
		},
		{
			name:           "Successful search by name",
			ctx:            context.Background(),
			query:          "test",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
			},
		},
		{
			name:           "Successful search by name",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			query:          "test",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().SearchByName(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/groups/search?query="+tt.query, nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleSearchByName(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		body           string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful create public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			body:           `name=test&description=test`,
			expectedStatus: http.StatusCreated,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&pgpb.CreateResponse{
					PublicGroup: &pgpb.PublicGroupResponse{
						Id: 1,
					},
				}, nil)
				userClient.EXPECT().CreatePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(&uspb.CreatePublicGroupAdminResponse{}, nil)
			},
		},
		{
			name:           "Successful create public group",
			ctx:            context.Background(),
			body:           `name=test&description=test`,
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful create public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			body:           `name=test&description=test`,
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful create public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			body:           `name=test&description=test`,
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, userClient *mock_user.MockUserClient) {
				publicGroupClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&pgpb.CreateResponse{
					PublicGroup: &pgpb.PublicGroupResponse{
						Id: 1,
					},
				}, nil)
				userClient.EXPECT().CreatePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			_ = writer.WriteField("name", "test")
			_ = writer.WriteField("description", "test")
			_ = writer.Close()

			r := httptest.NewRequest("POST", "/groups", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockPublicGroupClient, mockUserClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, mockUserClient)

			// Call the handler
			h.HandleCreate(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful update public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusOK,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&pgpb.UpdateResponse{
					PublicGroup: &pgpb.PublicGroupResponse{
						Id: 1,
					},
				}, nil)
			},
		},
		{
			name:           "Successful update public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful update public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful update public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Update(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			_ = writer.WriteField("name", "test")
			_ = writer.WriteField("description", "test")
			_ = writer.Close()
			r := httptest.NewRequest("PUT", "/groups/"+tt.groupID, body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			r = mux.SetURLVars(r, map[string]string{
				"groupID": tt.groupID,
			})

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleUpdate(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful delete public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusNoContent,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(&pgpb.DeleteResponse{}, nil)
			},
		},
		{
			name:           "Successful delete public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful delete public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful delete public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("DELETE", "/groups/"+tt.groupID, nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			r = mux.SetURLVars(r, map[string]string{
				"groupID": tt.groupID,
			})

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleDelete(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetSubscriptionByPublicGroupIDAndSubscriberID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful get subscription by public group ID and subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusOK,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetSubscriptionByPublicGroupIDAndSubscriberID(gomock.Any(), gomock.Any()).Return(&pgpb.GetSubscriptionByPublicGroupIDAndSubscriberIDResponse{
					Subscription: &pgpb.SubscriptionResponse{
						Id:            1,
						SubscriberId:  1,
						PublicGroupId: 1,
					},
				}, nil)
			},
		},
		{
			name:           "Successful get subscription by public group ID and subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get subscription by public group ID and subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get subscription by public group ID and subscriber ID",
			ctx:            context.Background(),
			groupID:        "1",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get subscription by public group ID and subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetSubscriptionByPublicGroupIDAndSubscriberID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/groups/"+tt.groupID+"/subscription", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			r = mux.SetURLVars(r, map[string]string{
				"groupID": tt.groupID,
			})

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleGetSubscriptionByPublicGroupIDAndSubscriberID(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetBySubscriberID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		userID         string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful get by subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "1",
			expectedStatus: http.StatusOK,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(&pgpb.GetBySubscriberIDResponse{
					PublicGroups: []*pgpb.PublicGroupResponse{
						{
							Id:               1,
							SubscribersCount: 1,
						},
					},
				}, nil)
			},
		},
		{
			name:           "Successful get by subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "",
			expectedStatus: http.StatusOK,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(&pgpb.GetBySubscriberIDResponse{
					PublicGroups: []*pgpb.PublicGroupResponse{
						{
							Id:               1,
							SubscribersCount: 1,
						},
					},
				}, nil)
			},
		},
		{
			name:           "Successful get by subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get by subscriber ID",
			ctx:            context.Background(),
			userID:         "1",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful get by subscriber ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			userID:         "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().GetBySubscriberID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/users/"+tt.userID+"/groups", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			r = mux.SetURLVars(r, map[string]string{
				"userID": tt.userID,
			})

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleGetBySubscriberID(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleSubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful subscribe to public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusCreated,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(&pgpb.SubscribeResponse{
					Subscription: &pgpb.SubscriptionResponse{
						Id:            1,
						SubscriberId:  1,
						PublicGroupId: 1,
					},
				}, nil)
			},
		},
		{
			name:           "Successful subscribe to public group",
			ctx:            context.Background(),
			groupID:        "1",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful subscribe to public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful subscribe to public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {

			},
		},
		{
			name:           "Successful subscribe to public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("POST", "/groups/"+tt.groupID+"/subscribe", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleSubscribe(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleUnsubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient)
	}{
		{
			name:           "Successful unsubscribe from public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusNoContent,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Unsubscribe(gomock.Any(), gomock.Any()).Return(&pgpb.UnsubscribeResponse{}, nil)
			},
		},
		{
			name:           "Successful unsubscribe from public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
			},
		},
		{
			name:           "Successful unsubscribe from public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
			},
		},
		{
			name:           "Successful unsubscribe from public group",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
				publicGroupClient.EXPECT().Unsubscribe(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful unsubscribe from public group",
			ctx:            context.Background(),
			groupID:        "1",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("DELETE", "/groups/"+tt.groupID+"/unsubscribe", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			tt.mock(mockPublicGroupClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, nil, nil)

			// Call the handler
			h.HandleUnsubscribe(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleCreateGroupPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient)
	}{
		{
			name:           "Successful create group post",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusCreated,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient) {
				postClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&postpb.CreatePostResponse{
					Post: &postpb.PostResponse{
						Id:       1,
						AuthorId: 1,
					},
				}, nil)
				postClient.EXPECT().CreateGroupPost(gomock.Any(), gomock.Any()).Return(&postpb.CreateGroupPostResponse{}, nil)
				publicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&pgpb.GetByIDResponse{
					PublicGroup: &pgpb.PublicGroupWithInfoResponse{
						PublicGroup: &pgpb.PublicGroupResponse{
							Id: 1,
						},
					},
				}, nil)
			},
		},
		{
			name:           "Successful create group post",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient) {

			},
		},
		{
			name:           "Successful create group post",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient) {

			},
		},
		{
			name:           "Successful create group post",
			ctx:            context.Background(),
			groupID:        "1",
			expectedStatus: http.StatusBadRequest,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient) {

			},
		},
		{
			name:           "Successful create group post",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient) {
				postClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
		{
			name:           "Successful create group post",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient) {
				postClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&postpb.CreatePostResponse{
					Post: &postpb.PostResponse{
						Id:       1,
						AuthorId: 1,
					},
				}, nil)
				postClient.EXPECT().CreateGroupPost(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)

			},
		},
		{
			name:           "Successful create group post",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(publicGroupClient *mock_public_group.MockPublicGroupClient, postClient *mock_post.MockPostClient) {
				postClient.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&postpb.CreatePostResponse{
					Post: &postpb.PostResponse{
						Id:       1,
						AuthorId: 1,
					},
				}, nil)
				postClient.EXPECT().CreateGroupPost(gomock.Any(), gomock.Any()).Return(&postpb.CreateGroupPostResponse{}, nil)
				publicGroupClient.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			err := writer.WriteField("content", "Test content")
			if err != nil {
				t.Fatal(err)
			}
			writer.Close()
			r := httptest.NewRequest("POST", "/groups/"+tt.groupID+"/posts", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockPublicGroupClient := mock_public_group.NewMockPublicGroupClient(ctrl)
			mockPostClient := mock_post.NewMockPostClient(ctrl)
			tt.mock(mockPublicGroupClient, mockPostClient)

			// Set up the handler
			h := NewPublicGroupHandler(mockPublicGroupClient, mockPostClient, nil)

			// Call the handler
			h.HandleCreateGroupPost(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetGroupPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		lastPostID     string
		postsAmount    string
		expectedStatus int
		mock           func(postClient *mock_post.MockPostClient)
	}{
		{
			name:           "Successful get group posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			lastPostID:     "0",
			postsAmount:    "10",
			expectedStatus: http.StatusOK,
			mock: func(postClient *mock_post.MockPostClient) {
				postClient.EXPECT().GetPostsOfGroup(gomock.Any(), gomock.Any()).Return(&postpb.GetPostsOfGroupResponse{
					Posts: []*postpb.PostResponse{
						{
							Id:       1,
							AuthorId: 1,
						},
					},
				}, nil)
			},
		},
		{
			name:           "Successful get group posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			lastPostID:     "0",
			postsAmount:    "10",
			expectedStatus: http.StatusBadRequest,
			mock: func(postClient *mock_post.MockPostClient) {

			},
		},
		{
			name:           "Successful get group posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			lastPostID:     "0",
			postsAmount:    "10",
			expectedStatus: http.StatusBadRequest,
			mock: func(postClient *mock_post.MockPostClient) {

			},
		},
		{
			name:           "Successful get group posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			lastPostID:     "asd",
			postsAmount:    "10",
			expectedStatus: http.StatusBadRequest,
			mock: func(postClient *mock_post.MockPostClient) {

			},
		},
		{
			name:           "Successful get group posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			lastPostID:     "10",
			postsAmount:    "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(postClient *mock_post.MockPostClient) {

			},
		},
		{
			name:           "Successful get group posts",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			lastPostID:     "10",
			postsAmount:    "10",
			expectedStatus: http.StatusInternalServerError,
			mock: func(postClient *mock_post.MockPostClient) {
				postClient.EXPECT().GetPostsOfGroup(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/groups/"+tt.groupID+"/posts?lastPostId="+tt.lastPostID+"&postsAmount="+tt.postsAmount, nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockPostClient := mock_post.NewMockPostClient(ctrl)
			tt.mock(mockPostClient)

			// Set up the handler
			h := NewPublicGroupHandler(nil, mockPostClient, nil)

			// Call the handler
			h.HandleGetGroupPosts(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleCreatePublicGroupAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		input          CreatePublicGroupAdminInput
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful create public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			input:          CreatePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusCreated,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().CreatePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(&uspb.CreatePublicGroupAdminResponse{}, nil)
			},
		},
		{
			name:           "Successful create public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			input:          CreatePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful create public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			input:          CreatePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful create public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			input:          CreatePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().CreatePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			inputBytes, _ := json.Marshal(tt.input)
			r := httptest.NewRequest("POST", "/groups/"+tt.groupID+"/admins", bytes.NewBuffer(inputBytes))
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewPublicGroupHandler(nil, nil, mockUserClient)

			// Call the handler
			h.HandleCreatePublicGroupAdmin(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleDeletePublicGroupAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		input          DeletePublicGroupAdminInput
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful delete public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			input:          DeletePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusNoContent,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().DeletePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(&uspb.DeletePublicGroupAdminResponse{}, nil)
			},
		},
		{
			name:           "Successful delete public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			input:          DeletePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful delete public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			input:          DeletePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful delete public group admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			input:          DeletePublicGroupAdminInput{UserID: 1},
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().DeletePublicGroupAdmin(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			inputBytes, _ := json.Marshal(tt.input)
			r := httptest.NewRequest("DELETE", "/groups/"+tt.groupID+"/admins", bytes.NewBuffer(inputBytes))
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewPublicGroupHandler(nil, nil, mockUserClient)

			// Call the handler
			h.HandleDeletePublicGroupAdmin(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleGetAdminsByPublicGroupID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful get admins by public group ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetAdminsByPublicGroupID(gomock.Any(), gomock.Any()).Return(&uspb.GetAdminsByPublicGroupIDResponse{}, nil)
			},
		},
		{
			name:           "Successful get admins by public group ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful get admins by public group ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {

			},
		},
		{
			name:           "Successful get admins by public group ID",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().GetAdminsByPublicGroupID(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/groups/"+tt.groupID+"/admins", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewPublicGroupHandler(nil, nil, mockUserClient)

			// Call the handler
			h.HandleGetAdminsByPublicGroupID(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestHandleCheckIfUserIsAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		ctx            context.Context
		groupID        string
		expectedStatus int
		mock           func(userClient *mock_user.MockUserClient)
	}{
		{
			name:           "Successful check if user is admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusOK,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().CheckIfUserIsAdmin(gomock.Any(), gomock.Any()).Return(&uspb.CheckIfUserIsAdminResponse{IsAdmin: true}, nil)
			},
		},
		{
			name:           "Successful check if user is admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "",
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful check if user is admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "asd",
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful check if user is admin",
			ctx:            context.Background(),
			groupID:        "1",
			expectedStatus: http.StatusBadRequest,
			mock: func(userClient *mock_user.MockUserClient) {
			},
		},
		{
			name:           "Successful check if user is admin",
			ctx:            context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			groupID:        "1",
			expectedStatus: http.StatusInternalServerError,
			mock: func(userClient *mock_user.MockUserClient) {
				userClient.EXPECT().CheckIfUserIsAdmin(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal.GRPCStatus().Err(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the request
			r := httptest.NewRequest("GET", "/groups/"+tt.groupID+"/admins/check", nil)
			r = r.WithContext(tt.ctx)

			// Set up the response recorder
			rr := httptest.NewRecorder()

			// Set up the mux vars
			muxVars := make(map[string]string)
			muxVars["groupID"] = tt.groupID
			r = mux.SetURLVars(r, muxVars)

			mockUserClient := mock_user.NewMockUserClient(ctrl)
			tt.mock(mockUserClient)

			// Set up the handler
			h := NewPublicGroupHandler(nil, nil, mockUserClient)

			// Call the handler
			h.HandleCheckIfUserIsAdmin(rr, r)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
