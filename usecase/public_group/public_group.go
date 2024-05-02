package publicgroup

import (
	"context"
	"socio/domain"
	"socio/errors"
	"socio/pkg/sanitizer"
	"socio/pkg/static"

	"github.com/microcosm-cc/bluemonday"
)

type PublicGroupWithInfo struct {
	PublicGroup  *domain.PublicGroup `json:"publicGroup"`
	IsSubscribed bool                `json:"isSubscribed"`
}

type PublicGroupStorage interface {
	GetPublicGroupByID(ctx context.Context, groupID uint, userID uint) (publicGroupWithInfo *PublicGroupWithInfo, err error)
	SearchPublicGroupsByNameWithInfo(ctx context.Context, query string, userID uint) (publicGroups []*PublicGroupWithInfo, err error)
	StorePublicGroup(ctx context.Context, publicGroup *domain.PublicGroup) (newGroup *domain.PublicGroup, err error)
	UpdatePublicGroup(ctx context.Context, publicGroup *domain.PublicGroup) (updatedGroup *domain.PublicGroup, err error)
	DeletePublicGroup(ctx context.Context, publicGroupID uint) (err error)
	GetSubscriptionByPublicGroupIDAndSubscriberID(ctx context.Context, publicGroupID, subscriberID uint) (subscription *domain.PublicGroupSubscription, err error)
	GetPublicGroupsBySubscriberID(ctx context.Context, subscriberID uint) (groups []*domain.PublicGroup, err error)
	StorePublicGroupSubscription(ctx context.Context, publicGroupSubscription *domain.PublicGroupSubscription) (newSubscription *domain.PublicGroupSubscription, err error)
	DeletePublicGroupSubscription(ctx context.Context, subscription *domain.PublicGroupSubscription) (err error)
	GetPublicGroupSubscriptionIDs(ctx context.Context, userID uint) (subIDs []uint, err error)
}

type AvatarStorage interface {
	Store(fileName string, filePath string, contentType string) (err error)
	Delete(fileName string) (err error)
}

type Service struct {
	PublicGroupStorage PublicGroupStorage
	AvatarStorage      AvatarStorage
	Sanitizer          *sanitizer.Sanitizer
}

func NewService(publicGroupStorage PublicGroupStorage, avatarStorage AvatarStorage) *Service {
	return &Service{
		PublicGroupStorage: publicGroupStorage,
		AvatarStorage:      avatarStorage,
		Sanitizer:          sanitizer.NewSanitizer(bluemonday.UGCPolicy()),
	}
}

func (s *Service) GetByID(ctx context.Context, groupID uint, userID uint) (publicGroup *PublicGroupWithInfo, err error) {
	if groupID == 0 {
		err = errors.ErrInvalidData
		return
	}

	publicGroup, err = s.PublicGroupStorage.GetPublicGroupByID(ctx, groupID, userID)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizePublicGroup(publicGroup.PublicGroup)

	return
}

func (s *Service) SearchByName(ctx context.Context, query string, userID uint) (publicGroups []*PublicGroupWithInfo, err error) {
	publicGroups, err = s.PublicGroupStorage.SearchPublicGroupsByNameWithInfo(ctx, query, userID)
	if err != nil {
		return
	}

	for _, publicGroup := range publicGroups {
		s.Sanitizer.SanitizePublicGroup(publicGroup.PublicGroup)
	}

	return
}

func (s *Service) Create(ctx context.Context, publicGroup *domain.PublicGroup) (newGroup *domain.PublicGroup, err error) {
	if len(publicGroup.Name) == 0 {
		err = errors.ErrInvalidData
		return
	}

	if len(publicGroup.Avatar) == 0 {
		publicGroup.Avatar = static.DefaultGroupAvatarFileName
	}

	newGroup, err = s.PublicGroupStorage.StorePublicGroup(ctx, publicGroup)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizePublicGroup(newGroup)

	return
}

func (s *Service) Update(ctx context.Context, publicGroup *domain.PublicGroup) (updatedGroup *domain.PublicGroup, err error) {
	if publicGroup.ID == 0 {
		err = errors.ErrInvalidData
		return
	}

	oldGroup, err := s.PublicGroupStorage.GetPublicGroupByID(ctx, publicGroup.ID, 0)
	if err != nil {
		return
	}

	if len(publicGroup.Name) > 0 {
		oldGroup.PublicGroup.Name = publicGroup.Name
	}

	oldGroup.PublicGroup.Description = publicGroup.Description

	if len(publicGroup.Avatar) > 0 {
		err = s.AvatarStorage.Delete(oldGroup.PublicGroup.Avatar)
		if err != nil {
			return
		}

		oldGroup.PublicGroup.Avatar = publicGroup.Avatar
	}

	updatedGroup, err = s.PublicGroupStorage.UpdatePublicGroup(ctx, oldGroup.PublicGroup)
	if err != nil {
		return
	}

	s.Sanitizer.SanitizePublicGroup(updatedGroup)

	return
}

func (s *Service) Delete(ctx context.Context, publicGroupID uint) (err error) {
	if publicGroupID == 0 {
		err = errors.ErrInvalidData
		return
	}

	err = s.PublicGroupStorage.DeletePublicGroup(ctx, publicGroupID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetSubscriptionByPublicGroupIDAndSubscriberID(ctx context.Context, publicGroupID, subscriberID uint) (subscription *domain.PublicGroupSubscription, err error) {
	if publicGroupID == 0 || subscriberID == 0 {
		err = errors.ErrInvalidData
		return
	}

	subscription, err = s.PublicGroupStorage.GetSubscriptionByPublicGroupIDAndSubscriberID(ctx, publicGroupID, subscriberID)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetBySubscriberID(ctx context.Context, subscriberID uint) (groups []*domain.PublicGroup, err error) {
	if subscriberID == 0 {
		err = errors.ErrInvalidData
		return
	}

	groups, err = s.PublicGroupStorage.GetPublicGroupsBySubscriberID(ctx, subscriberID)
	if err != nil {
		return
	}

	for _, group := range groups {
		s.Sanitizer.SanitizePublicGroup(group)
	}

	return
}

func (s *Service) Subscribe(ctx context.Context, publicGroupSubscription *domain.PublicGroupSubscription) (newSubscription *domain.PublicGroupSubscription, err error) {
	if publicGroupSubscription.PublicGroupID == 0 || publicGroupSubscription.SubscriberID == 0 {
		err = errors.ErrInvalidData
		return
	}

	_, err = s.PublicGroupStorage.GetPublicGroupByID(ctx, publicGroupSubscription.PublicGroupID, 0)
	if err != nil {
		return
	}

	newSubscription, err = s.PublicGroupStorage.StorePublicGroupSubscription(ctx, publicGroupSubscription)
	if err != nil {
		return
	}

	return
}

func (s *Service) Unsubscribe(ctx context.Context, publicGroupSubscription *domain.PublicGroupSubscription) (err error) {
	if publicGroupSubscription.PublicGroupID == 0 || publicGroupSubscription.SubscriberID == 0 {
		err = errors.ErrInvalidData
		return
	}

	err = s.PublicGroupStorage.DeletePublicGroupSubscription(ctx, publicGroupSubscription)
	if err != nil {
		return
	}

	return
}

func (s *Service) UploadAvatar(fileName string, filePath string, contentType string) (err error) {
	err = s.AvatarStorage.Store(fileName, filePath, contentType)
	if err != nil {
		return
	}

	return
}

func (s *Service) GetSubscriptionIDs(ctx context.Context, userID uint) (subIDs []uint, err error) {
	if userID == 0 {
		err = errors.ErrInvalidData
		return
	}

	subIDs, err = s.PublicGroupStorage.GetPublicGroupSubscriptionIDs(ctx, userID)
	if err != nil {
		return
	}

	return
}
