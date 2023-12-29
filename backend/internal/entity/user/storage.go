package user

import (
  "context"
  serviceError "coursework/internal/error"
  "coursework/internal/model"
  "coursework/pkg/log"
  "errors"
  "fmt"
  "github.com/jackc/pgconn"
  "github.com/jackc/pgx/v4/pgxpool"
  "strconv"
  "strings"
)

type database struct {
  pool   *pgxpool.Pool
  logger *log.Logger
}

type Storage interface {
  GetUserProfile(ctx context.Context, userId string) (string, error)
  UpdateUserProfile(ctx context.Context, userId string, user *model.User) (string, error)
  GetUserReports(ctx context.Context, userId string) (string, error)
  GetUserSupportCompany(ctx context.Context, userId string) (string, error)
  GetHouseStat(ctx context.Context, userId string) (string, error)
  GetHouseFeedbacks(ctx context.Context, userId string) (string, error)
  SendFeedback(ctx context.Context, userId string, dto *model.CreateFeedbackDTO) (string, error)
  SendReport(ctx context.Context, userId string, dto *model.CreateReportDTO) (string, error)
  IsUser(ctx context.Context, userId string) error
}

func NewStorage(storage *pgxpool.Pool, logger *log.Logger) Storage {
  return &database{pool: storage, logger: logger}
}

func (d *database) GetUserProfile(ctx context.Context, userId string) (string, error) {
  d.logger.Printf("getting user info")
  query := `SELECT * FROM get_user_profile($1)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when finding user info. error %v", err)
    return "", serviceError.NotFoundError
  }

  return jsonRes, nil
}

func (d *database) UpdateUserProfile(ctx context.Context, userId string, user *model.User) (string, error) {
  d.logger.Printf("updating user")
  intId, _ := strconv.Atoi(userId)
  query := `SELECT * FROM update_user_profile($1, $2, $3, $4, $5, $6)`

  d.logger.Printf("SQL Query: %s", query)
  _, err := d.pool.Exec(ctx, query, intId, user.FullName, user.Email, user.Phone, user.Street, user.House)
  if err != nil {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
      pgErr = err.(*pgconn.PgError)
      newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
      d.logger.Println(newErr)
    }

    d.logger.Printf("error when updating user. error %v", err)
    textErr := fmt.Sprintf("%v", err)
    textErr = strings.Replace(textErr, "ERROR: ", "", 1)
    if strings.HasPrefix(textErr, "Such") {
      return "", serviceError.HouseNotFound
    }

    return "", fmt.Errorf("failed to update user. error: %v", err)
  }

  return d.GetUserProfile(ctx, userId)
}

func (d *database) GetUserReports(ctx context.Context, userId string) (string, error) {
  d.logger.Printf("getting user reports")
  query := `SELECT * FROM get_user_reports($1)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when finding user reports. error %v", err)
    return "[]", nil
  }

  return jsonRes, nil
}

func (d *database) GetUserSupportCompany(ctx context.Context, userId string) (string, error) {
  d.logger.Printf("getting user support company info")
  query := `SELECT * FROM get_user_company($1)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when finding user support company. error %v", err)
    return "", serviceError.NotFoundError
  }

  return jsonRes, nil
}

func (d *database) GetHouseStat(ctx context.Context, userId string) (string, error) {
  d.logger.Printf("getting user house info")
  query := `SELECT * FROM get_house_stat($1)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when finding user house. error %v", err)
    return "", serviceError.NotFoundError
  }

  return jsonRes, nil
}

func (d *database) GetHouseFeedbacks(ctx context.Context, userId string) (string, error) {
  d.logger.Printf("getting user house info")
  query := `SELECT * FROM get_house_feedbacks($1)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when finding user house. error %v", err)
    return "[]", nil
  }

  return jsonRes, nil
}

func (d *database) SendFeedback(ctx context.Context, userId string, dto *model.CreateFeedbackDTO) (string, error) {
  d.logger.Printf("adding feedback")
  query := `SELECT * FROM add_feedback($1, $2)`

  d.logger.Printf("SQL Query: %s", query)
  _, err := d.pool.Exec(ctx, query, userId, dto.Text)
  if err != nil {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
      pgErr = err.(*pgconn.PgError)
      newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
      d.logger.Println(newErr)
    }
    d.logger.Printf("error when creating feedback. error %v", err)
    return "", fmt.Errorf("failed to create feedback. error: %v", err)
  }

  return d.GetHouseFeedbacks(ctx, userId)
}

func (d *database) SendReport(ctx context.Context, userId string, dto *model.CreateReportDTO) (string, error) {
  d.logger.Printf("adding report")
  query := `SELECT * FROM add_complaint($1, $2, $3)`

  d.logger.Printf("SQL Query: %s", query)
  _, err := d.pool.Exec(ctx, query, userId, dto.Title, dto.Text)
  if err != nil {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
      pgErr = err.(*pgconn.PgError)
      newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
      d.logger.Println(newErr)
    }
    d.logger.Printf("error when creating report. error %v", err)
    return "", fmt.Errorf("failed to create report. error: %v", err)
  }

  return `{"status": "success"}`, nil
}

func (d *database) IsUser(ctx context.Context, userId string) error {
  d.logger.Printf("getting user house info")
  query := `SELECT * FROM get_is_admin($1)`

  d.logger.Printf("SQL Query: %s", query)
  var isAdmin string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&isAdmin); err != nil {
    d.logger.Printf("error when finding user. error %v", err)
    return serviceError.NotFoundError
  }

  if isAdmin == "TRUE" {
    return serviceError.AccessDeniedError
  }

  return nil
}
