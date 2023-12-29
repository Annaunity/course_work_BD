package user

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
  profileURL  = "/api/user/:id/profile"
  reportURL   = "/api/user/:id/report"
  companyURL  = "/api/user/:id/support"
  houseURL    = "/api/user/:id/house"
  feedbackURL = "/api/user/:id/feedback"
)

type Handler struct {
  Logger  *log.Logger
  Service Service
}

func (handler *Handler) Register(router *httprouter.Router) {
  router.HandlerFunc(http.MethodGet, profileURL, middleware.AuthHandler(middleware.ErrorHandler(handler.GetUserProfile)))
  router.HandlerFunc(http.MethodPatch, profileURL, middleware.AuthHandler(middleware.ErrorHandler(handler.UpdateUserProfile)))
  router.HandlerFunc(http.MethodGet, companyURL, middleware.AuthHandler(middleware.ErrorHandler(handler.GetUserSupportCompany)))
  router.HandlerFunc(http.MethodGet, houseURL, middleware.AuthHandler(middleware.ErrorHandler(handler.GetHouseStat)))
  router.HandlerFunc(http.MethodGet, feedbackURL, middleware.AuthHandler(middleware.ErrorHandler(handler.GetHouseFeedback)))
  router.HandlerFunc(http.MethodPost, feedbackURL, middleware.AuthHandler(middleware.ErrorHandler(handler.SendFeedback)))
  router.HandlerFunc(http.MethodGet, reportURL, middleware.AuthHandler(middleware.ErrorHandler(handler.GetUserReports)))
  router.HandlerFunc(http.MethodPost, reportURL, middleware.AuthHandler(middleware.ErrorHandler(handler.SendReport)))
}

func (handler *Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to get user profile from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")

  jsonString, err := handler.Service.GetUserProfile(r.Context(), urlUserId)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to update user profile from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlUserId != userId {
    return serviceError.AccessDeniedError
  }

  err := handler.Service.IsUser(r.Context(), userId)
  if err != nil {
    return err
  }

  var updateUser model.UpdateUserDTO
  if err = json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
    handler.Logger.Printf("invalid JSON scheme. error: %v", err)
    return serviceError.BadRequestError("invalid JSON scheme")
  }
  handler.Logger.Println(updateUser.Street, updateUser.House)
  jsonString, err := handler.Service.UpdateUserProfile(r.Context(), userId, &updateUser)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) GetUserReports(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to get user reports from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlUserId != userId {
    return serviceError.AccessDeniedError
  }

  err := handler.Service.IsUser(r.Context(), userId)
  if err != nil {
    return err
  }

  jsonString, err := handler.Service.GetUserReports(r.Context(), urlUserId)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) GetUserSupportCompany(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to get company statistic from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlUserId != userId {
    return serviceError.AccessDeniedError
  }

  err := handler.Service.IsUser(r.Context(), userId)
  if err != nil {
    return err
  }

  jsonString, err := handler.Service.GetUserSupportCompanyProfile(r.Context(), urlUserId)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) GetHouseStat(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to get house statistic from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlUserId != userId {
    return serviceError.AccessDeniedError
  }

  err := handler.Service.IsUser(r.Context(), userId)
  if err != nil {
    return err
  }

  jsonString, err := handler.Service.GetHouseStat(r.Context(), urlUserId)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) GetHouseFeedback(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to get house feedbacks from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlUserId != userId {
    return serviceError.AccessDeniedError
  }

  err := handler.Service.IsUser(r.Context(), userId)
  if err != nil {
    return err
  }

  jsonString, err := handler.Service.GetHouseFeedbacks(r.Context(), urlUserId)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) SendFeedback(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to update user profile from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlUserId != userId {
    return serviceError.AccessDeniedError
  }

  err := handler.Service.IsUser(r.Context(), userId)
  if err != nil {
    return err
  }

  var feedback model.CreateFeedbackDTO
  if err = json.NewDecoder(r.Body).Decode(&feedback); err != nil {
    handler.Logger.Printf("invalid JSON scheme. error: %v", err)
    return serviceError.BadRequestError("invalid JSON scheme")
  }

  jsonString, err := handler.Service.SendFeedback(r.Context(), userId, &feedback)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}

func (handler *Handler) SendReport(w http.ResponseWriter, r *http.Request) error {
  handler.Logger.Printf("request to update user profile from %s", r.RemoteAddr)
  w.Header().Set("Content-Type", "application/json")

  handler.Logger.Printf("getting id from context")
  params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
  urlUserId := params.ByName("id")
  userId := r.Context().Value("user_id").(string)

  if urlUserId != userId {
    return serviceError.AccessDeniedError
  }

  err := handler.Service.IsUser(r.Context(), userId)
  if err != nil {
    return err
  }

  var report model.CreateReportDTO
  if err = json.NewDecoder(r.Body).Decode(&report); err != nil {
    handler.Logger.Printf("invalid JSON scheme. error: %v", err)
    return serviceError.BadRequestError("invalid JSON scheme")
  }

  jsonString, err := handler.Service.SendReport(r.Context(), userId, &report)
  if err != nil {
    return err
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(jsonString))

  return nil
}
