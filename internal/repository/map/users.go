package repository

import (
	"socio/domain"
	"socio/errors"
	"socio/utils"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Users      *sync.Map
	NextUserId uint
	TP         utils.TimeProvider
}

func NewUsers(tp utils.TimeProvider, users *sync.Map) (s *Users) {
	s = &Users{}
	s.Users = users
	s.NextUserId = 2
	s.TP = tp

	salt1 := "salt"
	dateOfBirth, _ := time.Parse(utils.DateFormat, "1990-01-01")
	user1 := &domain.User{
		ID:        0,
		FirstName: "Petr",
		LastName:  "Mitin",
		Password:  utils.HashPassword("admin1", []byte(salt1)),
		Salt:      salt1,
		Email:     "petr09mitin@mail.ru",
		RegistrationDate: utils.CustomTime{
			Time: tp.Now(),
		},
		Avatar: "default_avatar.png",
		DateOfBirth: utils.CustomTime{
			Time: dateOfBirth,
		},
	}
	s.Users.Store(user1.ID, user1)

	salt2 := "salt"
	user2 := &domain.User{
		ID:        1,
		FirstName: "Alexey",
		LastName:  "Gorbunov",
		Password:  utils.HashPassword("admin2", []byte(salt2)),
		Salt:      salt2,
		Email:     "lexagorbunov14@gmail.com",
		RegistrationDate: utils.CustomTime{
			Time: tp.Now(),
		},
		Avatar: "leha.jpg",
		DateOfBirth: utils.CustomTime{
			Time: dateOfBirth,
		},
	}
	s.Users.Store(user2.ID, user2)

	return
}

func (s *Users) GetUserById(userID uint) (user *domain.User, err error) {
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

func (s *Users) StoreUser(user *domain.User) {
	salt := uuid.NewString()
	user.ID = s.NextUserId
	user.Password = utils.HashPassword(user.Password, []byte(salt))
	user.Salt = salt
	user.RegistrationDate = utils.CustomTime{
		Time: s.TP.Now(),
	}

	s.Users.Store(user.ID, user)

	s.NextUserId++
}

func (s *Users) RefreshSaltAndRehashPassword(user *domain.User) {
	salt := uuid.NewString()
	user.Password = utils.HashPassword(user.Password, []byte(salt))
	user.Salt = salt

	s.Users.Store(user.ID, user)
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
