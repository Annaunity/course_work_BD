package error

import (
  "encoding/json"
  "fmt"
)

type ServiceError struct {
  Err              error  `json:"-"`
  Message          string `json:"error,omitempty"`
  DeveloperMessage string `json:"developer_message,omitempty"`
  Code             string `json:"code,omitempty"`
}

func NewServiceError(message, code, developerMessage string) *ServiceError {
  return &ServiceError{
    Err:              fmt.Errorf(message),
    Code:             code,
    Message:          message,
    DeveloperMessage: developerMessage,
  }
}

// Error needed to implement Error interface
func (e *ServiceError) Error() string {
  return e.Err.Error()
}

// Unwrap needed to easy unwrap
func (e *ServiceError) Unwrap() error {
  return e.Err
}

func (e *ServiceError) Marshal() []byte {
  bytes, err := json.Marshal(e)
  if err != nil {
    return nil
  }

  return bytes
}

// Default application errors

var (
  NotFoundError     = NewServiceError("не удалось найти", "CW-0003", "")
  AccessDeniedError = NewServiceError("нет доступа", "CW-0004", "")
  AlreadyExistError = NewServiceError("пользователь уже существует", "CW-0005", "")
  WeakPasswordError = NewServiceError("пароль слишком легкий", "CW-0006", "")
  InvalidUserError  = NewServiceError("неправильный логин или пароль", "CW-0008", "")
  HouseNotFound     = NewServiceError("такого дома не существует", "CW-0009", "")
)

func InternalError(developerMessage string) *ServiceError {
  return NewServiceError("internal server error", "CW-0001", developerMessage)
}

func BadRequestError(message string) *ServiceError {
  return NewServiceError(message, "CW-0002", "wrong request data")
}
