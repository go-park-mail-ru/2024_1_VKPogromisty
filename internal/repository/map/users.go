package repository

import (
	"socio/domain"
	"socio/errors"
	"socio/pkg/hash"
	customtime "socio/pkg/time"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Users      *sync.Map
	NextUserId uint
	TP         customtime.TimeProvider
}

func NewUsers(tp customtime.TimeProvider, users *sync.Map) (s *Users) {
	s = &Users{}
	s.Users = users
	s.NextUserId = 2
	s.TP = tp

	salt1 := "salt"
	dateOfBirth, _ := time.Parse(customtime.DateFormat, "1990-01-01")
	user1 := &domain.User{
		ID:        0,
		FirstName: "Petr",
		LastName:  "Mitin",
		Password:  hash.HashPassword("admin1", []byte(salt1)),
		Salt:      salt1,
		Email:     "petr09mitin@mail.ru",
		CreatedAt: customtime.CustomTime{
			Time: tp.Now(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: tp.Now(),
		},
		Avatar: "default_avatar.png",
		DateOfBirth: customtime.CustomTime{
			Time: dateOfBirth,
		},
	}
	s.Users.Store(user1.ID, user1)

	salt2 := "salt"
	user2 := &domain.User{
		ID:        1,
		FirstName: "Alexey",
		LastName:  "Gorbunov",
		Password:  hash.HashPassword("admin2", []byte(salt2)),
		Salt:      salt2,
		Email:     "lexagorbunov14@gmail.com",
		CreatedAt: customtime.CustomTime{
			Time: tp.Now(),
		},
		UpdatedAt: customtime.CustomTime{
			Time: tp.Now(),
		},
		Avatar: "leha.jpg",
		DateOfBirth: customtime.CustomTime{
			Time: dateOfBirth,
		},
	}
	s.Users.Store(user2.ID, user2)

	return
}

func (s *Users) GetUserByEmail(email string) (user *domain.User, err error) {
	s.Users.Range(func(key, value interface{}) bool {
		currUser := value.(*domain.User)
		if currUser.Email == email {
			user = currUser
			return false
		}
		return true
	})

	if user == nil {
		err = errors.ErrNotFound
		return
	}

	return
}

func (s *Users) StoreUser(user *domain.User) (err error) {
	salt := uuid.NewString()
	user.ID = s.NextUserId
	user.Password = hash.HashPassword(user.Password, []byte(salt))
	user.Salt = salt
	user.CreatedAt = customtime.CustomTime{
		Time: s.TP.Now(),
	}
	user.UpdatedAt = customtime.CustomTime{
		Time: s.TP.Now(),
	}

	s.Users.Store(user.ID, user)

	s.NextUserId++
	return
}

func (s *Users) RefreshSaltAndRehashPassword(user *domain.User, password string) (err error) {
	salt := uuid.NewString()
	user.Password = hash.HashPassword(password, []byte(salt))
	user.Salt = salt

	s.Users.Store(user.ID, user)

	return
}

func (s *Users) GetUserByID(userID uint) (user *domain.User, err error) {
	userData, ok := s.Users.Load(userID)
	if !ok {
		err = errors.ErrNotFound
		return
	}

	user, ok = userData.(*domain.User)
	if !ok {
		err = errors.ErrInternal
		return
	}

	return
}
