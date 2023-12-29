package admin

import (
  serviceError "coursework/internal/error"
  "coursework/internal/middleware"
  "coursework/internal/model"
  "coursework/pkg/log"
  "encoding/json"
  "github.com/julienschmidt/httprouter"
  "net/http"
)

const (
  reportURL = "/api/admin/:id/report"
)

type Handler struct {
  Logger  *log.Logger
  Service Service
}

func (handler *Handler) Register(router *httprouter.Router) {
  router.HandlerFunc(http.MethodGet, reportURL, middleware.AuthHandler(middleware.ErrorHandler(handler.GetNextReport)))
  router.HandlerFunc(http.MethodPost, reportURL, middleware.AuthHandler(middleware.ErrorHandler(handler.CloseReport)))
}

func (handler *Handler) GetNextReport(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to get next report from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlAdminId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlAdminId != userId {
    return serviceError.AccessDeniedError
  }

  if err := handler.Service.IsAdmin(r.Context(), urlAdminId); err != nil {
    return err
  }

  jsonString, err := handler.Service.GetNextReport(r.Context(), urlAdminId)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) CloseReport(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to close report from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlAdminId := params.ByName("id")
  if urlAdminId != urlAdminId {
    return serviceError.AccessDeniedError
  }

  if err := handler.Service.IsAdmin(r.Context(), urlAdminId); err != nil {
    return err
  }

  var reportDTO model.CloseReportDTO
  if err := json.NewDecoder(r.Body).Decode(&reportDTO); err != nil {
    handler.Logger.Printf("invalid JSON scheme. error: %v", err)
    return serviceError.BadRequestError("invalid JSON scheme")
  }

  if reportDTO.Rating < 0 || reportDTO.Rating > 5 {
    handler.Logger.Printf("invalid rating. error")
    return serviceError.BadRequestError("rating must be in [0.0;5.0]")
  }

  jsonString, err := handler.Service.CloseReport(r.Context(), urlAdminId, &reportDTO)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}
