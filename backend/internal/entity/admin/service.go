package admin

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
  GetNextReport(ctx context.Context, userId string) (string, error)
  CloseReport(ctx context.Context, userId string, dto *model.CloseReportDTO) (string, error)
  IsAdmin(ctx context.Context, userId string) error
}

func NewService(storage Storage, logger *log.Logger) Service {
  return &service{logger: logger, storage: storage}
}

func (s service) GetNextReport(ctx context.Context, userId string) (string, error) {
  jsonString, err := s.storage.GetNextReport(ctx, userId)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s service) CloseReport(ctx context.Context, userId string, dto *model.CloseReportDTO) (string, error) {
  jsonString, err := s.storage.CloseReport(ctx, userId, dto)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s service) IsAdmin(ctx context.Context, userId string) error {
  err := s.storage.IsAdmin(ctx, userId)
  if err != nil {
    return err
  }

  return nil
}
