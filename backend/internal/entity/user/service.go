package user

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
  GetUserProfile(ctx context.Context, userId string) (string, error)
  UpdateUserProfile(ctx context.Context, userId string, dto *model.UpdateUserDTO) (string, error)
  GetUserReports(ctx context.Context, userId string) (string, error)
  GetUserSupportCompanyProfile(ctx context.Context, userId string) (string, error)
  GetHouseStat(ctx context.Context, userId string) (string, error)
  GetHouseFeedbacks(ctx context.Context, userId string) (string, error)
  SendFeedback(ctx context.Context, userId string, feedback *model.CreateFeedbackDTO) (string, error)
  SendReport(ctx context.Context, userId string, report *model.CreateReportDTO) (string, error)
  IsUser(ctx context.Context, userId string) error
}

func NewService(storage Storage, logger *log.Logger) Service {
  return &service{logger: logger, storage: storage}
}

func (s *service) GetUserProfile(ctx context.Context, userId string) (string, error) {
  jsonString, err := s.storage.GetUserProfile(ctx, userId)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) UpdateUserProfile(ctx context.Context, userId string, dto *model.UpdateUserDTO) (string, error) {
  user := dto.NewUser()

  jsonString, err := s.storage.UpdateUserProfile(ctx, userId, user)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) GetUserReports(ctx context.Context, userId string) (string, error) {
  jsonString, err := s.storage.GetUserReports(ctx, userId)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) GetUserSupportCompanyProfile(ctx context.Context, userId string) (string, error) {
  jsonString, err := s.storage.GetUserSupportCompany(ctx, userId)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) GetHouseStat(ctx context.Context, userId string) (string, error) {
  jsonString, err := s.storage.GetHouseStat(ctx, userId)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) GetHouseFeedbacks(ctx context.Context, userId string) (string, error) {
  jsonString, err := s.storage.GetHouseFeedbacks(ctx, userId)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) SendFeedback(ctx context.Context, userId string, feedback *model.CreateFeedbackDTO) (string, error) {
  jsonString, err := s.storage.SendFeedback(ctx, userId, feedback)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) SendReport(ctx context.Context, userId string, report *model.CreateReportDTO) (string, error) {
  jsonString, err := s.storage.SendReport(ctx, userId, report)
  if err != nil {
    return "", err
  }

  return jsonString, nil
}

func (s *service) IsUser(ctx context.Context, userId string) error {
  err := s.storage.IsUser(ctx, userId)
  if err != nil {
    return err
  }

  return nil
}
