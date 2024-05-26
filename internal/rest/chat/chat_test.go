package rest_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"socio/domain"
	"socio/errors"
	rest "socio/internal/rest/chat"
	mock_rest "socio/mocks/rest/chat"
	"socio/pkg/requestcontext"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetDialogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name        string
		userID      uint
		ctx         context.Context
		wantDialogs []*domain.Dialog
		wantErr     error
		setup       func()
	}{
		{
			name:   "test case 1 - successful retrieval",
			userID: 1,
			wantDialogs: []*domain.Dialog{
				{
					User1: &domain.User{ID: 1},
					User2: &domain.User{ID: 2},
					LastMessage: &domain.PersonalMessage{
						ID: 1,
					},
				},
			},
			wantErr: nil,
			ctx:     context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().GetDialogsByUserID(gomock.Any(), gomock.Any()).Return([]*domain.Dialog{
					{
						User1: &domain.User{ID: 1},
						User2: &domain.User{ID: 2},
						LastMessage: &domain.PersonalMessage{
							ID: 1,
						},
					},
				}, nil)
			},
		},
		{
			name:   "test case 2",
			userID: 1,
			wantDialogs: []*domain.Dialog{
				{
					User1: &domain.User{ID: 1},
					User2: &domain.User{ID: 2},
					LastMessage: &domain.PersonalMessage{
						ID: 1,
					},
				},
			},
			wantErr: errors.ErrInternal,
			ctx:     context.Background(),
			setup: func() {
			},
		},
		{
			name:   "test case 3",
			userID: 1,
			wantDialogs: []*domain.Dialog{
				{
					User1: &domain.User{ID: 1},
					User2: &domain.User{ID: 2},
					LastMessage: &domain.PersonalMessage{
						ID: 1,
					},
				},
			},
			wantErr: errors.ErrNotFound,
			ctx:     context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().GetDialogsByUserID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("GET", "/dialogs", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleGetDialogs)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHandleGetMessagesByDialog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name         string
		query        string
		ctx          context.Context
		wantMessages []*domain.PersonalMessage
		wantErr      error
		setup        func()
	}{
		{
			name:  "test case 1 - successful retrieval",
			query: "/dialogs?peerId=2&lastMessageId=&messagesAmount=",
			wantMessages: []*domain.PersonalMessage{
				{
					ID: 1,
				},
			},
			wantErr: nil,
			ctx:     context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().GetMessagesByDialog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:  "test case 2",
			query: "/dialogs?peerId=2&lastMessageId=1&messagesAmount=1",
			wantMessages: []*domain.PersonalMessage{
				{
					ID: 1,
				},
			},
			wantErr: errors.ErrForbidden,
			ctx:     context.Background(),
			setup: func() {
			},
		},
		{
			name:         "test case 3",
			query:        "/dialogs?peerId=asd&lastMessageId=1&messagesAmount=1",
			wantMessages: nil,
			wantErr:      errors.ErrInternal,
			ctx:          context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
		},
		{
			name:         "test case 4",
			query:        "/dialogs?peerId=1&lastMessageId=asd&messagesAmount=1",
			wantMessages: nil,
			wantErr:      errors.ErrInternal,
			ctx:          context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
		},
		{
			name:         "test case 5",
			query:        "/dialogs?peerId=1&lastMessageId=1&messagesAmount=asd4",
			wantMessages: nil,
			wantErr:      errors.ErrInternal,
			ctx:          context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
		},
		{
			name:         "test case 6",
			query:        "/dialogs?peerId=&lastMessageId=1&messagesAmount=1",
			wantMessages: nil,
			wantErr:      errors.ErrInternal,
			ctx:          context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
		},
		{
			name:  "test case 7",
			query: "/dialogs?peerId=2&lastMessageId=&messagesAmount=",
			wantMessages: []*domain.PersonalMessage{
				{
					ID: 1,
				},
			},
			wantErr: errors.ErrInternal,
			ctx:     context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().GetMessagesByDialog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("GET", tt.query, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleGetMessagesByDialog)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHandleGetStickersByAuthorID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name         string
		authorID     uint
		ctx          context.Context
		wantStickers []*domain.Sticker
		wantErr      error
		setup        func()
		muxVars      map[string]string
	}{
		{
			name:     "test case 1 - successful retrieval",
			authorID: 1,
			wantStickers: []*domain.Sticker{
				{
					ID: 1,
				},
			},
			wantErr: nil,
			ctx:     context.Background(),
			setup: func() {
				mockService.EXPECT().GetStickersByAuthorID(gomock.Any(), gomock.Any()).Return([]*domain.Sticker{
					{
						ID: 1,
					},
				}, nil)
			},
			muxVars: map[string]string{
				"authorID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:     "test case 2",
			authorID: 1,
			wantStickers: []*domain.Sticker{
				{
					ID: 1,
				},
			},
			wantErr: errors.ErrInvalidBody,
			ctx:     context.Background(),
			setup: func() {

			},
			muxVars: map[string]string{},
		},
		{
			name:     "test case 3",
			authorID: 1,
			wantStickers: []*domain.Sticker{
				{
					ID: 1,
				},
			},
			wantErr: errors.ErrInternal,
			ctx:     context.Background(),
			setup: func() {

			},
			muxVars: map[string]string{
				"authorID": "oppa",
			},
		},
		{
			name:     "test case 1 - successful retrieval",
			authorID: 1,
			wantStickers: []*domain.Sticker{
				{
					ID: 1,
				},
			},
			wantErr: errors.ErrInternal,
			ctx:     context.Background(),
			setup: func() {
				mockService.EXPECT().GetStickersByAuthorID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrInternal)
			},
			muxVars: map[string]string{
				"authorID": fmt.Sprintf("%d", 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("GET", "/{authorID}/stickers/", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			req = mux.SetURLVars(req, tt.muxVars)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleGetStickersByAuthorID)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHandleGetAllStickers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name         string
		ctx          context.Context
		wantStickers []*domain.Sticker
		wantErr      error
		setup        func()
	}{
		{
			name: "test case 1 - successful retrieval",
			wantStickers: []*domain.Sticker{
				{
					ID: 1,
					// Fill other fields as needed
				},
			},
			wantErr: nil,
			ctx:     context.Background(),
			setup: func() {
				mockService.EXPECT().GetAllStickers(gomock.Any()).Return([]*domain.Sticker{
					{
						ID: 1,
						// Fill other fields as needed
					},
				}, nil)
			},
		},
		{
			name: "test case 2",
			wantStickers: []*domain.Sticker{
				{
					ID: 1,
					// Fill other fields as needed
				},
			},
			wantErr: errors.ErrInternal,
			ctx:     context.Background(),
			setup: func() {
				mockService.EXPECT().GetAllStickers(gomock.Any()).Return(nil, errors.ErrInternal)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("GET", "/stickers", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleGetAllStickers)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHandleDeleteSticker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name      string
		stickerID uint
		userID    uint
		ctx       context.Context
		wantErr   error
		setup     func()
		muxVars   map[string]string
	}{
		{
			name:      "test case 1 - successful deletion",
			stickerID: 1,
			userID:    1,
			wantErr:   nil,
			ctx:       context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().DeleteSticker(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			muxVars: map[string]string{
				"stickerID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:      "test case 2",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrInternal,
			ctx:       context.Background(),
			setup: func() {
			},
			muxVars: map[string]string{
				"stickerID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:      "test case 3",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrInvalidSlug,
			ctx:       context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
			muxVars: map[string]string{},
		},
		{
			name:      "test case 4",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrInvalidSlug,
			ctx:       context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
			muxVars: map[string]string{
				"stickerID": "oppa",
			},
		},
		{
			name:      "test case 5",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrInternal,
			ctx:       context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().DeleteSticker(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.ErrInternal)
			},
			muxVars: map[string]string{
				"stickerID": fmt.Sprintf("%d", 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/stickers/%d", tt.stickerID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			req = mux.SetURLVars(req, tt.muxVars)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleDeleteSticker)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusNoContent, rr.Code)
			} else {
				assert.Equal(t, http.StatusNoContent, rr.Code)
			}
		})
	}
}

func TestHandleGetUnsentMessageAttachments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name       string
		senderID   uint
		receiverID uint
		ctx        context.Context
		wantErr    error
		setup      func()
		muxVars    map[string]string
	}{
		{
			name:       "test case 1 - successful retrieval",
			senderID:   1,
			receiverID: 1,
			wantErr:    nil,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().GetUnsentMessageAttachments(gomock.Any(), gomock.Any()).Return([]string{}, nil)
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:       "test case 2",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInvalidBody,
			ctx:        context.Background(),
			setup: func() {
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:       "test case 3",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInvalidBody,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
			muxVars: map[string]string{},
		},
		{
			name:       "test case 4",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInvalidBody,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
			muxVars: map[string]string{
				"receiverID": "oppa",
			},
		},
		{
			name:       "test case 5",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInternal,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().GetUnsentMessageAttachments(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrInternal,
				)
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("GET", fmt.Sprintf("/users/%d/unsent-attachments", tt.receiverID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			req = mux.SetURLVars(req, tt.muxVars)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleGetUnsentMessageAttachments)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusOK, rr.Code)
			} else {
				assert.Equal(t, http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHandleDeleteUnsentMessageAttachments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name       string
		senderID   uint
		receiverID uint
		ctx        context.Context
		wantErr    error
		setup      func()
		muxVars    map[string]string
	}{
		{
			name:       "test case 1 - successful deletion",
			senderID:   1,
			receiverID: 1,
			wantErr:    nil,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().DeleteUnsentMessageAttachments(gomock.Any(), gomock.Any()).Return(nil)
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:       "test case 2",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInvalidBody,
			ctx:        context.Background(),
			setup: func() {
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:       "test case 3",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInvalidBody,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
			muxVars: map[string]string{},
		},
		{
			name:       "test case 4",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInvalidBody,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
			},
			muxVars: map[string]string{
				"receiverID": "oppa",
			},
		},
		{
			name:       "test case 5",
			senderID:   1,
			receiverID: 1,
			wantErr:    errors.ErrInternal,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().DeleteUnsentMessageAttachments(gomock.Any(), gomock.Any()).Return(
					errors.ErrInternal,
				)
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%d/unsent-attachments", tt.receiverID), nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			req = mux.SetURLVars(req, tt.muxVars)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleDeleteUnsentMessageAttachments)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusNoContent, rr.Code)
			} else {
				assert.Equal(t, http.StatusNoContent, rr.Code)
			}
		})
	}
}

func TestHandleDeleteUnsentMessageAttachment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_rest.NewMockChatService(ctrl)

	tests := []struct {
		name       string
		senderID   uint
		receiverID uint
		fileName   string
		ctx        context.Context
		wantErr    error
		setup      func()
		muxVars    map[string]string
	}{
		{
			name:       "test case 1 - successful deletion",
			senderID:   1,
			receiverID: 1,
			fileName:   "testfile",
			wantErr:    nil,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().DeleteUnsentMessageAttachment(gomock.Any(), &domain.UnsentMessageAttachment{
					SenderID:   1,
					ReceiverID: 1,
					FileName:   "testfile",
				}).Return(nil)
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
				"fileName":   "testfile",
			},
		},
		{
			name:       "test case 2",
			senderID:   1,
			receiverID: 1,
			fileName:   "testfile",
			wantErr:    errors.ErrInternal,
			ctx:        context.Background(),
			setup: func() {

			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
				"fileName":   "testfile",
			},
		},
		{
			name:       "test case 3",
			senderID:   1,
			receiverID: 1,
			fileName:   "testfile",
			wantErr:    errors.ErrInternal,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {

			},
			muxVars: map[string]string{
				"fileName": "testfile",
			},
		},
		{
			name:       "test case 4",
			senderID:   1,
			receiverID: 1,
			fileName:   "testfile",
			wantErr:    errors.ErrInternal,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {

			},
			muxVars: map[string]string{
				"receiverID": "asd",
				"fileName":   "testfile",
			},
		},
		{
			name:       "test case 5",
			senderID:   1,
			receiverID: 1,
			fileName:   "testfile",
			wantErr:    errors.ErrInternal,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {

			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
			},
		},
		{
			name:       "test case 6",
			senderID:   1,
			receiverID: 1,
			fileName:   "testfile",
			wantErr:    errors.ErrInternal,
			ctx:        context.WithValue(context.Background(), requestcontext.UserIDKey, uint(1)),
			setup: func() {
				mockService.EXPECT().DeleteUnsentMessageAttachment(gomock.Any(), &domain.UnsentMessageAttachment{
					SenderID:   1,
					ReceiverID: 1,
					FileName:   "testfile",
				}).Return(errors.ErrInternal)
			},
			muxVars: map[string]string{
				"receiverID": fmt.Sprintf("%d", 1),
				"fileName":   "testfile",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			c := rest.NewChatServer(mockService)

			req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%d/unsent-attachments/%s", tt.receiverID, tt.fileName), nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(tt.ctx)

			req = mux.SetURLVars(req, tt.muxVars)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(c.HandleDeleteUnsentMessageAttachment)

			handler.ServeHTTP(rr, req)

			if tt.wantErr != nil {
				assert.NotEqual(t, http.StatusNoContent, rr.Code)
			} else {
				assert.Equal(t, http.StatusNoContent, rr.Code)
			}
		})
	}
}
