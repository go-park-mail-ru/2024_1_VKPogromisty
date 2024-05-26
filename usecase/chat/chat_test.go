package chat_test

import (
	"context"
	errorsDef "errors"
	"reflect"
	"socio/domain"
	"socio/errors"
	mock_chat "socio/mocks/usecase/chat"
	"socio/pkg/sanitizer"
	"socio/usecase/chat"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
)

type MockSyncMap struct {
	sync.Map
}

func (m *MockSyncMap) LoadAndDelete(key interface{}) (interface{}, bool) {
	// Mock the behavior of LoadAndDelete here
	// For example, always return true
	return nil, true
}

func TestUnregister(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		expectedErr error
		syncMap     *sync.Map
	}{
		{
			name:        "TestUnregister",
			userID:      1,
			expectedErr: nil,
			syncMap: func() *sync.Map {
				m := &sync.Map{}
				m.Store(uint(1), &chat.Client{})
				return m
			}(),
		},
		{
			name:        "TestUnregister",
			userID:      1,
			expectedErr: errors.ErrNotFound,
			syncMap: func() *sync.Map {
				m := &sync.Map{}
				return m
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &chat.Service{Clients: tt.syncMap}

			err := s.Unregister(tt.userID)

			if !errorsDef.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

type fields struct {
	PersonalMessagesRepo *mock_chat.MockPersonalMessagesRepository
	Sanitizer            *sanitizer.Sanitizer
}

func TestGetMessagesByDialog(t *testing.T) {
	tests := []struct {
		name            string
		userID          uint
		peerID          uint
		lastMessageID   uint
		messagesAmount  uint
		expectedErr     error
		expectedMessage []*domain.PersonalMessage
		prepare         func(f *fields)
	}{
		{
			name:           "TestGetMessagesByDialog",
			userID:         1,
			peerID:         2,
			lastMessageID:  0,
			messagesAmount: 0,
			expectedErr:    nil,
			expectedMessage: []*domain.PersonalMessage{
				{
					ID:         1,
					SenderID:   1,
					ReceiverID: 2,
					Content:    "Hello",
				},
			},
			prepare: func(f *fields) {
				f.PersonalMessagesRepo.EXPECT().GetLastMessageID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint(1), nil)
				f.PersonalMessagesRepo.EXPECT().GetMessagesByDialog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*domain.PersonalMessage{
					{
						ID:         1,
						SenderID:   1,
						ReceiverID: 2,
						Content:    "Hello",
					},
				}, nil)
			},
		},
		{
			name:            "TestGetMessagesByDialog",
			userID:          1,
			peerID:          2,
			lastMessageID:   0,
			messagesAmount:  0,
			expectedErr:     errors.ErrNotFound,
			expectedMessage: nil,
			prepare: func(f *fields) {
				f.PersonalMessagesRepo.EXPECT().GetLastMessageID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint(0), errors.ErrNotFound)
			},
		},
		{
			name:            "TestGetMessagesByDialog",
			userID:          1,
			peerID:          2,
			lastMessageID:   0,
			messagesAmount:  0,
			expectedErr:     errors.ErrNotFound,
			expectedMessage: nil,
			prepare: func(f *fields) {
				f.PersonalMessagesRepo.EXPECT().GetLastMessageID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint(1), nil)
				f.PersonalMessagesRepo.EXPECT().GetMessagesByDialog(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := &fields{
				PersonalMessagesRepo: mock_chat.NewMockPersonalMessagesRepository(ctrl),
			}

			tt.prepare(fields)

			s := chat.NewChatService(nil, nil, fields.PersonalMessagesRepo, nil, nil)

			messages, err := s.GetMessagesByDialog(context.Background(), tt.userID, tt.peerID, tt.lastMessageID, tt.messagesAmount)

			if !errorsDef.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(messages, tt.expectedMessage) {
				t.Errorf("expected messages %v, got %v", tt.expectedMessage, messages)
			}
		})
	}
}

func TestGetDialogsByUserID(t *testing.T) {
	tests := []struct {
		name            string
		userID          uint
		peerID          uint
		lastMessageID   uint
		expectedErr     error
		expectedDialogs []*domain.Dialog
		prepare         func(f *fields)
	}{
		{
			name:          "GetDialogsByUserID",
			userID:        1,
			peerID:        2,
			lastMessageID: 0,
			expectedErr:   nil,
			expectedDialogs: []*domain.Dialog{
				{
					User1: &domain.User{
						ID: 1,
					},
					User2: &domain.User{
						ID: 2,
					},
					LastMessage: &domain.PersonalMessage{
						ID:         1,
						SenderID:   1,
						ReceiverID: 2,
						Content:    "Hello",
					},
				},
			},
			prepare: func(f *fields) {
				f.PersonalMessagesRepo.EXPECT().GetDialogsByUserID(gomock.Any(), gomock.Any()).Return([]*domain.Dialog{
					{
						User1: &domain.User{
							ID: 1,
						},
						User2: &domain.User{
							ID: 2,
						},
						LastMessage: &domain.PersonalMessage{
							ID:         1,
							SenderID:   1,
							ReceiverID: 2,
							Content:    "Hello",
						},
					},
				}, nil)
			},
		},
		{
			name:            "GetDialogsByUserID",
			userID:          1,
			peerID:          2,
			lastMessageID:   0,
			expectedErr:     errors.ErrNotFound,
			expectedDialogs: nil,
			prepare: func(f *fields) {
				f.PersonalMessagesRepo.EXPECT().GetDialogsByUserID(gomock.Any(), gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fields := &fields{
				PersonalMessagesRepo: mock_chat.NewMockPersonalMessagesRepository(ctrl),
				Sanitizer:            sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
			}

			tt.prepare(fields)

			s := chat.NewChatService(nil, nil, fields.PersonalMessagesRepo, nil, nil)

			dialogs, err := s.GetDialogsByUserID(context.Background(), tt.userID)

			if !errorsDef.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if !reflect.DeepEqual(dialogs, tt.expectedDialogs) {
				t.Errorf("expected messages %v, got %v", tt.expectedDialogs, dialogs)
			}
		})
	}
}

func TestGetClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		userID  uint
		want    *chat.Client
		wantErr error
		setup   func(s *chat.Service)
	}{
		{
			name:   "test case 1 - client found",
			userID: 1,
			want: &chat.Client{
				UserID: 1,
			},
			wantErr: nil,
			setup: func(s *chat.Service) {
				s.Clients.Store(uint(1), &chat.Client{
					UserID: 1,
				})
			},
		},
		{
			name:    "test case 2 - client not found",
			userID:  2,
			want:    nil,
			wantErr: errors.ErrNotFound,
			setup:   func(s *chat.Service) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := chat.NewChatService(nil, nil, nil, nil, nil)

			tt.setup(s)

			got, err := s.GetClient(context.Background(), tt.userID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetStickersByAuthorID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessagesRepo := mock_chat.NewMockPersonalMessagesRepository(ctrl)

	tests := []struct {
		name     string
		authorID uint
		want     []*domain.Sticker
		wantErr  error
		setup    func()
	}{
		{
			name:     "test case 1 - successful retrieval",
			authorID: 1,
			want: []*domain.Sticker{
				{
					ID:       1,
					AuthorID: 1,
					Name:     "sticker1",
					FileName: "sticker1.png",
				},
			},
			wantErr: nil,
			setup: func() {
				mockMessagesRepo.EXPECT().GetStickersByAuthorID(gomock.Any(), uint(1)).Return([]*domain.Sticker{
					{
						ID:       1,
						AuthorID: 1,
						Name:     "sticker1",
						FileName: "sticker1.png",
					},
				}, nil)
			},
		},
		{
			name:     "test case 2",
			authorID: 1,
			want:     []*domain.Sticker{},
			wantErr:  errors.ErrNotFound,
			setup: func() {
				mockMessagesRepo.EXPECT().GetStickersByAuthorID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := chat.NewChatService(nil, nil, mockMessagesRepo, nil, nil)

			got, err := s.GetStickersByAuthorID(context.Background(), tt.authorID)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetAllStickers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessagesRepo := mock_chat.NewMockPersonalMessagesRepository(ctrl)

	tests := []struct {
		name    string
		want    []*domain.Sticker
		wantErr error
		setup   func()
	}{
		{
			name: "test case 1 - successful retrieval",
			want: []*domain.Sticker{
				{
					ID:       1,
					AuthorID: 1,
					Name:     "sticker1",
					FileName: "sticker1.png",
				},
			},
			wantErr: nil,
			setup: func() {
				mockMessagesRepo.EXPECT().GetAllStickers(gomock.Any()).Return([]*domain.Sticker{
					{
						ID:       1,
						AuthorID: 1,
						Name:     "sticker1",
						FileName: "sticker1.png",
					},
				}, nil)
			},
		},
		{
			name: "test case 1 - successful retrieval",
			want: []*domain.Sticker{
				{
					ID:       1,
					AuthorID: 1,
					Name:     "sticker1",
					FileName: "sticker1.png",
				},
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockMessagesRepo.EXPECT().GetAllStickers(gomock.Any()).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := chat.NewChatService(nil, nil, mockMessagesRepo, nil, nil)

			got, err := s.GetAllStickers(context.Background())
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteSticker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessagesRepo := mock_chat.NewMockPersonalMessagesRepository(ctrl)
	mockStickerStorage := mock_chat.NewMockStickerStorage(ctrl)

	tests := []struct {
		name      string
		stickerID uint
		userID    uint
		wantErr   error
		setup     func()
	}{
		{
			name:      "test case 1 - successful deletion",
			stickerID: 1,
			userID:    1,
			wantErr:   nil,
			setup: func() {
				mockMessagesRepo.EXPECT().GetStickerByID(gomock.Any(), uint(1)).Return(&domain.Sticker{AuthorID: 1, FileName: "sticker1.png"}, nil)
				mockStickerStorage.EXPECT().Delete("sticker1.png").Return(nil)
				mockMessagesRepo.EXPECT().DeleteSticker(gomock.Any(), uint(1)).Return(nil)
			},
		},
		{
			name:      "test case 2",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrNotFound,
			setup: func() {
				mockMessagesRepo.EXPECT().GetStickerByID(gomock.Any(), uint(1)).Return(&domain.Sticker{AuthorID: 1, FileName: "sticker1.png"}, nil)
				mockStickerStorage.EXPECT().Delete("sticker1.png").Return(nil)
				mockMessagesRepo.EXPECT().DeleteSticker(gomock.Any(), uint(1)).Return(errors.ErrNotFound)
			},
		},
		{
			name:      "test case 2",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrNotFound,
			setup: func() {
				mockMessagesRepo.EXPECT().GetStickerByID(gomock.Any(), uint(1)).Return(&domain.Sticker{AuthorID: 1, FileName: "sticker1.png"}, nil)
				mockStickerStorage.EXPECT().Delete("sticker1.png").Return(errors.ErrNotFound)
			},
		},
		{
			name:      "test case 3",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrForbidden,
			setup: func() {
				mockMessagesRepo.EXPECT().GetStickerByID(gomock.Any(), uint(1)).Return(&domain.Sticker{AuthorID: 2, FileName: "sticker1.png"}, nil)
			},
		},
		{
			name:      "test case 4",
			stickerID: 1,
			userID:    1,
			wantErr:   errors.ErrNotFound,
			setup: func() {
				mockMessagesRepo.EXPECT().GetStickerByID(gomock.Any(), uint(1)).Return(nil, errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := chat.NewChatService(nil, nil, mockMessagesRepo, mockStickerStorage, nil)

			err := s.DeleteSticker(context.Background(), tt.stickerID, tt.userID)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
