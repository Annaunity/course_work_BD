package auth

import (
  "context"
  "coursework/internal/model"
  "coursework/pkg/log"
)

type service struct {
  storage Storage
  logger  *log.Logger
}

type Service interface {
  Create(ctx context.Context, dto *model.CreateUserDTO) (*model.User, error)
  Login(ctx context.Context, dto *model.AuthUserDTO) (*model.User, error)
}

func NewService(storage Storage, logger *log.Logger) Service {
  return &service{logger: logger, storage: storage}
}

func (s *service) Create(ctx context.Context, dto *model.CreateUserDTO) (*model.User, error) {
  s.logger.Printf("try to register user with email: %s in database", dto.FullName)

  user := dto.NewUser()

  user, err := s.storage.Create(ctx, user)
  if err != nil {
    return nil, err
  }

  return user, nil
}

func (s *service) Login(ctx context.Context, dto *model.AuthUserDTO) (*model.User, error) {
  s.logger.Printf("try to login user with email: %s in database", dto.FullName)

  user := dto.NewUser()

  user, err := s.storage.Login(ctx, user)
  if err != nil {
    return nil, err
  }

  return user, nil
}
