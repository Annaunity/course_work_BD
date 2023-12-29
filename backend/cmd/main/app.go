package main

import (
  "context"
  "coursework/internal/config"
  "coursework/internal/entity/admin"
  "coursework/internal/entity/auth"
  "coursework/internal/entity/user"
  "coursework/pkg/cache/freecache"
  "coursework/pkg/jwt"
  "coursework/pkg/log"
  "coursework/pkg/pgsql"
  "coursework/pkg/shutdown"
  "errors"
  "fmt"
  "github.com/jackc/pgx/v4/pgxpool"
  "github.com/julienschmidt/httprouter"
  "net"
  "net/http"
  "os"
  "syscall"
  "time"
)

func main() {
  log.Init()
  logger := log.GetLogger()
  logger.Println("logger initialized")

  conf := config.GetConfig()
  logger.Println("config initialized")

  router := httprouter.New()
  router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Access-Control-Request-Method") != "" {
      header := w.Header()
      header.Set("Access-Control-Allow-Headers", header.Get("Allow"))
      header.Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers, authorization")
      header.Set("Access-Control-Allow-Origin", "*")
    }
    w.WriteHeader(http.StatusNoContent)
  })
  logger.Println("router initialized")

  postgresClient, err := pgsql.NewClient(
    context.Background(),
    conf.Postgres.Host,
    conf.Postgres.Port,
    conf.Postgres.Username,
    conf.Postgres.Password,
    conf.Postgres.Database,
  )
  if err != nil {
    logger.Fatalln(err)
  }
  logger.Println("postgres client initialized")

  refreshTokenCache := freecache.NewCacheRepo(104857600) // 100MB
  logger.Println("cache initialized")

  jwtHelper := jwt.NewHelper(refreshTokenCache, logger)
  logger.Println("jwt initialized")

  authStorage := auth.NewStorage(postgresClient, logger)
  authService := auth.NewService(authStorage, logger)
  authHandler := auth.Handler{Logger: logger, JWTHelper: jwtHelper, Service: authService}
  authHandler.Register(router)
  logger.Println("auth handler initialized")

  userStorage := user.NewStorage(postgresClient, logger)
  userService := user.NewService(userStorage, logger)
  userHandler := user.Handler{Logger: logger, Service: userService}
  userHandler.Register(router)
  logger.Println("user handler initialized")

  adminStorage := admin.NewStorage(postgresClient, logger)
  adminService := admin.NewService(adminStorage, logger)
  adminHandler := admin.Handler{Logger: logger, Service: adminService}
  adminHandler.Register(router)
  logger.Println("admin handler initialized")

  start(router, logger, conf, postgresClient)
}

func start(router *httprouter.Router, logger *log.Logger, conf *config.Config, pgClient *pgxpool.Pool) {
  logger.Println("server starting")
  logger.Printf("bing application to host %s and port %s", conf.Listen.BindIp, conf.Listen.Port)
  listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", conf.Listen.BindIp, conf.Listen.Port))

  if err != nil {
    logger.Fatal(err)
  }

  server := &http.Server{
    Handler:      router,
    WriteTimeout: 30 * time.Second,
    ReadTimeout:  30 * time.Second,
  }

  go shutdown.Hook([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, server)

  logger.Printf("server started and listening on %s:%s", conf.Listen.BindIp, conf.Listen.Port)

  if err := server.Serve(listener); err != nil {
    pgClient.Close()
    switch {
    case errors.Is(err, http.ErrServerClosed):
      logger.Println("server close")
    default:
      logger.Fatalln(err)
    }
  }
}
