package admin

import (
  "context"
  serviceError "coursework/internal/error"
  "coursework/internal/model"
  "coursework/pkg/log"
  "github.com/jackc/pgx/v4/pgxpool"
)

type database struct {
  pool   *pgxpool.Pool
  logger *log.Logger
}

type Storage interface {
  GetNextReport(ctx context.Context, userId string) (string, error)
  CloseReport(ctx context.Context, userId string, user *model.CloseReportDTO) (string, error)
  IsAdmin(ctx context.Context, userId string) error
}

func NewStorage(storage *pgxpool.Pool, logger *log.Logger) Storage {
  return &database{pool: storage, logger: logger}
}

func (d database) GetNextReport(ctx context.Context, userId string) (string, error) {
  d.logger.Printf("getting next report")
  query := `SELECT * FROM get_next_report_card($1)`

  d.logger.Printf("SQL Query: %s", query)
  var jsonRes string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&jsonRes); err != nil {
    d.logger.Printf("error when getting next report. error %v", err)
    return "", serviceError.NotFoundError
  }

  return jsonRes, nil
}

func (d database) CloseReport(ctx context.Context, userId string, dto *model.CloseReportDTO) (string, error) {
  d.logger.Printf("closing report")
  query := `SELECT * FROM close_complaint($1, $2, $3)`

  d.logger.Printf("SQL Query: %s", query)
  if err := d.pool.QueryRow(ctx, query, dto.ComplaintId, dto.Verdict, dto.Rating); err != nil {
    d.logger.Printf("error when closing report. error %v", err)
    return "", serviceError.NotFoundError
  }

  return d.GetNextReport(ctx, userId)
}

func (d database) IsAdmin(ctx context.Context, userId string) error {
  d.logger.Printf("check that admin")
  query := `SELECT * FROM get_is_admin($1)`

  d.logger.Printf("SQL Query: %s", query)
  var isAdmin string
  if err := d.pool.QueryRow(ctx, query, userId).Scan(&isAdmin); err != nil {
    d.logger.Printf("error when finding user. error %v", err)
    return serviceError.NotFoundError
  }

  if isAdmin == "FALSE" {
    return serviceError.AccessDeniedError
  }

  return nil
}
