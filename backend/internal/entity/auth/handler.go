package auth

import (
  serviceError "coursework/internal/error"
  "coursework/internal/middleware"
  "coursework/internal/model"
  "coursework/pkg/jwt"
  "coursework/pkg/log"
  "encoding/json"
  "github.com/julienschmidt/httprouter"
  "net/http"
)

const (
  authURL     = "/api/auth"
  registerURL = "/api/register"
)

type Handler struct {
  Logger    *log.Logger
  Service   Service
  JWTHelper jwt.Helper
}

func (handler *Handler) Register(router *httprouter.Router) {
  router.HandlerFunc(http.MethodPost, registerURL, middleware.ErrorHandler(handler.Create))
  router.HandlerFunc(http.MethodPost, authURL, middleware.ErrorHandler(handler.Auth))
  router.HandlerFunc(http.MethodPatch, authURL, middleware.ErrorHandler(handler.Auth))
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to sign up from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  var dto model.CreateUserDTO
  if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
    handler.Logger.Printf("invalid JSON scheme. error: %v", err)
    return serviceError.BadRequestError("invalid JSON scheme")
  }

  user, err := handler.Service.Create(r.Context(), &dto)
  if err != nil {
    return err
  }

  handler.Logger.Printf("user %s was created", user.FullName)

  token, err := handler.JWTHelper.GenerateAccessToken(user)
  if err != nil {
    handler.Logger.Printf("failed to generate access token of user %s. error: %v", user.FullName, err)
    return serviceError.BadRequestError("failed to generate access token")
  }

  w.WriteHeader(http.StatusCreated)
  w.Write(token)

  return nil
}

func (handler *Handler) Auth(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to sign in from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  var token []byte
  var err error

  switch r.Method {
  case http.MethodPost:
    var dto model.AuthUserDTO
    if err = json.NewDecoder(r.Body).Decode(&dto); err != nil {
      handler.Logger.Printf("invalid JSON scheme. error: %v", err)
      return serviceError.BadRequestError("invalid JSON scheme")
    }

    user, err := handler.Service.Login(r.Context(), &dto)
    if err != nil {
      return err
    }

    handler.Logger.Printf("user %s was auth", user.FullName)
    token, err = handler.JWTHelper.GenerateAccessToken(user)
    if err != nil {
      handler.Logger.Printf("failed to generate access token of user %s. error: %v", user.FullName, err)
      return serviceError.BadRequestError("failed to generate access token")
    }
  case http.MethodPatch:
    var rt jwt.RT
    if err = json.NewDecoder(r.Body).Decode(&rt); err != nil {
      handler.Logger.Printf("invalid JSON scheme. error: %v", err)
      return serviceError.BadRequestError("invalid JSON scheme")
    }

    token, err = handler.JWTHelper.UpdateRefreshToken(&rt)
    if err != nil {
      handler.Logger.Printf("failed to update access token of user. error: %v", err)
      return serviceError.BadRequestError("failed to update access token")
    }
  }

  w.WriteHeader(http.StatusOK)
  w.Write(token)

  return nil
}
