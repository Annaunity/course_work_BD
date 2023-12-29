package middleware

import (
  "context"
  "coursework/internal/config"
  serviceError "coursework/internal/error"
  serviceJwt "coursework/pkg/jwt"
  "coursework/pkg/log"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/cristalhq/jwt/v3"
  "net/http"
  "strings"
  "time"
)

type serviceHandler func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(handler serviceHandler) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Access-Control-Allow-Origin", "*")

    var serviceErr *serviceError.ServiceError
    err := handler(writer, request)

    if err == nil {
      return
    }

    // unexpected error (418 code + internal error)
    statusCode := http.StatusTeapot
    handlerError := serviceError.InternalError(err.Error())

    // expected error
    if errors.As(err, &serviceErr) {
      log.GetLogger().Logger.Printf(serviceErr.Code)
      switch {
      case errors.Is(err, serviceError.NotFoundError):
        statusCode = http.StatusNotFound
        handlerError = serviceError.NotFoundError
      case errors.Is(err, serviceError.AlreadyExistError):
        statusCode = http.StatusOK
        handlerError = serviceError.AlreadyExistError
      case errors.Is(err, serviceError.WeakPasswordError):
        statusCode = http.StatusOK
        handlerError = serviceError.WeakPasswordError
      case errors.Is(err, serviceError.AlreadyExistError):
        statusCode = http.StatusOK
        handlerError = serviceError.AlreadyExistError
      case errors.Is(err, serviceError.InvalidUserError):
        statusCode = http.StatusOK
        handlerError = serviceError.InvalidUserError
      case errors.Is(err, serviceError.HouseNotFound):
        statusCode = http.StatusOK
        handlerError = serviceError.HouseNotFound
      default:
        statusCode = http.StatusBadRequest
        handlerError = err.(*serviceError.ServiceError)
      }
    }

    // send to client
    writer.WriteHeader(statusCode)
    writer.Write(handlerError.Marshal())
  }
}

func AuthHandler(handler http.HandlerFunc) http.HandlerFunc {
  return func(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Access-Control-Allow-Origin", "*")

    logger := log.GetLogger()
    authHeader := strings.Split(request.Header.Get("Authorization"), "Bearer ")
    if len(authHeader) != 2 {
      unauthorized(writer, fmt.Errorf("no access token"))
      return
    }

    logger.Println("create jwt verifier")
    jwtToken := authHeader[1]
    key := []byte(config.GetConfig().JWT.Secret)
    verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
    if err != nil {
      unauthorized(writer, err)
    }

    logger.Println("parse and verify token")
    token, err := jwt.ParseAndVerifyString(jwtToken, verifier)
    if err != nil {
      unauthorized(writer, err)
      return
    }

    logger.Println("parse user claims")
    var userClaims serviceJwt.UserClaims
    err = json.Unmarshal(token.RawClaims(), &userClaims)
    if err != nil {
      unauthorized(writer, err)
      return
    }

    if valid := userClaims.IsValidAt(time.Now()); !valid {
      logger.Println("token has been expired")
      unauthorized(writer, err)
      return
    }

    ctx := context.WithValue(request.Context(), "user_id", userClaims.ID)
    handler(writer, request.WithContext(ctx))
  }
}

func unauthorized(w http.ResponseWriter, err error) {
  log.GetLogger().Println(err)
  w.WriteHeader(http.StatusUnauthorized)
  w.Write([]byte("unauthorized"))
}
