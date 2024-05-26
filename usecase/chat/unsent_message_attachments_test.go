package chat_test

import (
	"context"
	"socio/domain"
	"socio/errors"
	mock_chat "socio/mocks/usecase/chat"
	"socio/usecase/chat"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUnsentMessageAttachments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUnsentMessageAttachmentsStorage := mock_chat.NewMockUnsentMessageAttachmentsStorage(ctrl)

	tests := []struct {
		name    string
		attach  *domain.UnsentMessageAttachment
		want    []string
		wantErr error
		setup   func()
	}{
		{
			name: "test case 1 - successful retrieval",
			attach: &domain.UnsentMessageAttachment{
				SenderID:   1,
				ReceiverID: 2,
				FileName:   "file1.png",
			},
			want:    []string{"file1.png"},
			wantErr: nil,
			setup: func() {
				mockUnsentMessageAttachmentsStorage.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]string{"file1.png"}, nil)
			},
		},
		{
			name: "test case 2",
			attach: &domain.UnsentMessageAttachment{
				SenderID:   1,
				ReceiverID: 2,
				FileName:   "file1.png",
			},
			want:    []string{"file1.png"},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockUnsentMessageAttachmentsStorage.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrNotFound,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := chat.NewChatService(nil, mockUnsentMessageAttachmentsStorage, nil, nil, nil)

			got, err := s.GetUnsentMessageAttachments(context.Background(), tt.attach)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeleteUnsentMessageAttachments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUnsentMessageAttachmentsStorage := mock_chat.NewMockUnsentMessageAttachmentsStorage(ctrl)
	mockMessageAttachmentStorage := mock_chat.NewMockMessageAttachmentStorage(ctrl)

	tests := []struct {
		name    string
		attach  *domain.UnsentMessageAttachment
		wantErr error
		setup   func()
	}{
		{
			name: "test case 1 - successful deletion",
			attach: &domain.UnsentMessageAttachment{
				SenderID:   1,
				ReceiverID: 2,
				FileName:   "file1.png",
			},
			wantErr: nil,
			setup: func() {
				mockUnsentMessageAttachmentsStorage.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]string{"file1.png"}, nil)
				mockMessageAttachmentStorage.EXPECT().Delete("file1.png").Return(nil)
				mockUnsentMessageAttachmentsStorage.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "test case 2",
			attach: &domain.UnsentMessageAttachment{
				SenderID:   1,
				ReceiverID: 2,
				FileName:   "file1.png",
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockUnsentMessageAttachmentsStorage.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]string{"file1.png"}, nil)
				mockMessageAttachmentStorage.EXPECT().Delete("file1.png").Return(nil)
				mockUnsentMessageAttachmentsStorage.EXPECT().DeleteAll(gomock.Any(), gomock.Any()).Return(errors.ErrNotFound)
			},
		},
		{
			name: "test case 3",
			attach: &domain.UnsentMessageAttachment{
				SenderID:   1,
				ReceiverID: 2,
				FileName:   "file1.png",
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockUnsentMessageAttachmentsStorage.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return([]string{"file1.png"}, nil)
				mockMessageAttachmentStorage.EXPECT().Delete("file1.png").Return(errors.ErrNotFound)
			},
		},
		{
			name: "test case 4",
			attach: &domain.UnsentMessageAttachment{
				SenderID:   1,
				ReceiverID: 2,
				FileName:   "file1.png",
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockUnsentMessageAttachmentsStorage.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(
					nil, errors.ErrNotFound,
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := chat.NewChatService(nil, mockUnsentMessageAttachmentsStorage, nil, nil, mockMessageAttachmentStorage)

			err := s.DeleteUnsentMessageAttachments(context.Background(), tt.attach)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestDeleteUnsentMessageAttachment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUnsentMessageAttachmentsStorage := mock_chat.NewMockUnsentMessageAttachmentsStorage(ctrl)
	mockMessageAttachmentStorage := mock_chat.NewMockMessageAttachmentStorage(ctrl)

	tests := []struct {
		name    string
		attach  *domain.UnsentMessageAttachment
		wantErr error
		setup   func()
	}{
		{
			name: "test case 1 - successful deletion",
			attach: &domain.UnsentMessageAttachment{
				FileName: "file1.png",
			},
			wantErr: nil,
			setup: func() {
				mockMessageAttachmentStorage.EXPECT().Delete("file1.png").Return(nil)
				mockUnsentMessageAttachmentsStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name: "test case 2",
			attach: &domain.UnsentMessageAttachment{
				FileName: "file1.png",
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockMessageAttachmentStorage.EXPECT().Delete("file1.png").Return(nil)
				mockUnsentMessageAttachmentsStorage.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.ErrNotFound)
			},
		},
		{
			name: "test case 3",
			attach: &domain.UnsentMessageAttachment{
				FileName: "file1.png",
			},
			wantErr: errors.ErrNotFound,
			setup: func() {
				mockMessageAttachmentStorage.EXPECT().Delete("file1.png").Return(errors.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			s := chat.NewChatService(nil, mockUnsentMessageAttachmentsStorage, nil, nil, mockMessageAttachmentStorage)

			err := s.DeleteUnsentMessageAttachment(context.Background(), tt.attach)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
