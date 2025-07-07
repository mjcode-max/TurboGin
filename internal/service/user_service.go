package service

import (
	"github.com/mjcode-max/TurboGin/internal/dao"
	"github.com/mjcode-max/TurboGin/internal/model"
	"github.com/mjcode-max/TurboGin/pkg/logger"
	"github.com/mjcode-max/TurboGin/pkg/redis"
)

type IUserService interface {
	GetUser(id uint) (*model.User, error)
	CreateUser(user *model.User) error
}

type UserService struct {
	userDao dao.IUserDAO
	log     *logger.Logger
	client  *redis.Client
}

func NewUserService(userDao dao.IUserDAO, log *logger.Logger, client *redis.Client) IUserService {
	return &UserService{userDao: userDao, log: log, client: client}
}

func (s *UserService) GetUser(id uint) (*model.User, error) {
	return s.userDao.GetByID(id)
}

func (s *UserService) CreateUser(user *model.User) error {
	return s.userDao.Create(user)
}
