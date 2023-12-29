package jwt

import (
  "coursework/internal/config"
  "coursework/internal/model"
  "coursework/pkg/cache"
  "coursework/pkg/log"
  "encoding/json"
  "github.com/cristalhq/jwt/v3"
  "github.com/google/uuid"
  "time"
)

type UserClaims struct {
  jwt.RegisteredClaims
  Email string `json:"email"`
}

type RT struct {
  RefreshToken string `json:"refresh_token"`
}

type helper struct {
  Logger  *log.Logger
  RTCache cache.Repository
}

type Helper interface {
  GenerateAccessToken(user *model.User) ([]byte, error)
  UpdateRefreshToken(rt *RT) ([]byte, error)
}

func NewHelper(RTCache cache.Repository, logger *log.Logger) Helper {
  return &helper{RTCache: RTCache, Logger: logger}
}

func (h *helper) GenerateAccessToken(user *model.User) ([]byte, error) {
  key := []byte(config.GetConfig().JWT.Secret)
  signer, err := jwt.NewSignerHS(jwt.HS256, key)
  if err != nil {
    return nil, err
  }

  builder := jwt.NewBuilder(signer)

  claims := UserClaims{
    RegisteredClaims: jwt.RegisteredClaims{
      ID:        user.Id,
      Audience:  []string{"users"},
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
    },
    Email: user.FullName,
  }

  token, err := builder.Build(claims)
  if err != nil {
    h.Logger.Println(err)
    return nil, err
  }

  h.Logger.Printf("create refresh token for %s", user.Id)
  refreshTokenUUID := uuid.New()
  userBytes, _ := json.Marshal(user)
  err = h.RTCache.Set([]byte(refreshTokenUUID.String()), userBytes, 0)
  if err != nil {
    h.Logger.Println(err)
    return nil, err
  }

  jsonBytes, err := json.Marshal(map[string]string{
    "token":         token.String(),
    "refresh_token": refreshTokenUUID.String(),
    "user_id":       user.Id,
    "is_admin":      user.IsAdmin,
  })
  if err != nil {
    return nil, err
  }

  h.Logger.Printf("token for %s was created. token: %s. refresh: %s", user.Id, token.String(), refreshTokenUUID.String())

  return jsonBytes, nil
}

func (h *helper) UpdateRefreshToken(rt *RT) ([]byte, error) {
  defer h.RTCache.Del([]byte(rt.RefreshToken))

  userBytes, err := h.RTCache.Get([]byte(rt.RefreshToken))
  if err != nil {
    return nil, err
  }

  var user model.User
  h.Logger.Println(user.Email, user.Id, user.FullName, user.Password)
  err = json.Unmarshal(userBytes, &user)
  if err != nil {
    return nil, err
  }

  return h.GenerateAccessToken(&user)
}
