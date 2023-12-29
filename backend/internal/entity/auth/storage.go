package auth

import (
  "context"
  serviceError "coursework/internal/error"
  "coursework/internal/model"
  "coursework/pkg/log"
  "fmt"
  "github.com/jackc/pgx/v4/pgxpool"
  "strings"
)

type database struct {
  pool   *pgxpool.Pool
  logger *log.Logger
}

type Storage interface {
  Create(ctx context.Context, user *model.User) (*model.User, error)
  Login(ctx context.Context, user *model.User) (*model.User, error)
}

func NewStorage(storage *pgxpool.Pool, logger *log.Logger) Storage {
  return &database{pool: storage, logger: logger}
}

func (d *database) Create(ctx context.Context, regUser *model.User) (*model.User, error) {
  d.logger.Printf("register new user %s", regUser.FullName)
  query := `SELECT * FROM register($1, $2, $3, $4, $5, $6, $7)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, regUser.FullName, false, regUser.Password, "", "", regUser.Street, regUser.House).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when creating user. error %v", err)
    return nil, fmt.Errorf("failed to create user")
  }

  if strings.HasPrefix(jsonRes, "House") {
    return nil, serviceError.HouseNotFound
  } else if strings.HasPrefix(jsonRes, "User") {
    return nil, serviceError.AlreadyExistError
  }

  d.logger.Printf("result of query %s", jsonRes)
  fields := strings.Split(jsonRes, " ")
  regUser.Id = fields[0]
  regUser.IsAdmin = fields[1]

  return regUser, nil
}

func (d *database) Login(ctx context.Context, loginUser *model.User) (*model.User, error) {
  d.logger.Printf("auth user %s", loginUser.FullName)
  query := `SELECT * FROM login($1, $2)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, loginUser.FullName, loginUser.Password).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when creating user. error %v", err)
    return nil, fmt.Errorf("failed to create user")
  }

  if strings.HasPrefix(jsonRes, "Invalid fullname") {
    return loginUser, serviceError.InvalidUserError
  } else if strings.HasPrefix(jsonRes, "Invalid password") {
    return loginUser, serviceError.InvalidUserError
  }

  d.logger.Printf("result of query %s", loginUser)
  fields := strings.Split(jsonRes, " ")
  loginUser.Id = fields[0]
  loginUser.IsAdmin = fields[1]

  return loginUser, nil
}
