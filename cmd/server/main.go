package main

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    chimiddleware "github.com/go-chi/chi/v5/middleware"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/ruturaj/taskflow/internal/handler"
    "github.com/ruturaj/taskflow/internal/middleware"
    "github.com/ruturaj/taskflow/internal/repository"
    "github.com/ruturaj/taskflow/internal/service"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    dbURL := mustEnv("DB_URL")
    jwtSecret := mustEnv("JWT_SECRET")
    port := getEnv("PORT", "8080")

    // Run migrations
    m, err := migrate.New("file://migrations", dbURL)
    if err != nil {
        slog.Error("failed to init migrations", "err", err)
        os.Exit(1)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        slog.Error("migration failed", "err", err)
        os.Exit(1)
    }
    slog.Info("migrations applied")

    // Run seed
    // (seed via separate SQL executed in entrypoint script or init container)

    // DB connection pool
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        slog.Error("failed to connect to db", "err", err)
        os.Exit(1)
    }
    defer pool.Close()

    // Wire up layers
    userRepo    := repository.NewUserRepository(pool)
    projectRepo := repository.NewProjectRepository(pool)
    taskRepo    := repository.NewTaskRepository(pool)

    authSvc    := service.NewAuthService(userRepo, jwtSecret)
    projectSvc := service.NewProjectService(projectRepo, taskRepo)
    taskSvc    := service.NewTaskService(taskRepo, projectRepo)

    authH    := handler.NewAuthHandler(authSvc)
    projectH := handler.NewProjectHandler(projectSvc)
    taskH    := handler.NewTaskHandler(taskSvc)

    authMiddleware := middleware.Auth(authSvc)

    // Router
    r := chi.NewRouter()
    r.Use(chimiddleware.RequestID)
    r.Use(chimiddleware.RealIP)
    r.Use(chimiddleware.Logger)       // request logging
    r.Use(chimiddleware.Recoverer)    // panic recovery → 500 instead of crash

    r.Post("/auth/register", authH.Register)
    r.Post("/auth/login", authH.Login)

    r.Group(func(r chi.Router) {
        r.Use(authMiddleware)

        r.Get("/projects", projectH.List)
        r.Post("/projects", projectH.Create)
        r.Get("/projects/{id}", projectH.Get)
        r.Patch("/projects/{id}", projectH.Update)
        r.Delete("/projects/{id}", projectH.Delete)
        r.Get("/projects/{id}/stats", projectH.Stats) // bonus

        r.Get("/projects/{id}/tasks", taskH.List)
        r.Post("/projects/{id}/tasks", taskH.Create)
        r.Patch("/tasks/{id}", taskH.Update)
        r.Delete("/tasks/{id}", taskH.Delete)
    })

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%s", port),
        Handler:      r,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Graceful shutdown
    go func() {
        slog.Info("server starting", "port", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("server error", "err", err)
            os.Exit(1)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    slog.Info("shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("shutdown error", "err", err)
    }
    slog.Info("server stopped")
}

func mustEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        slog.Error("required env var not set", "key", key)
        os.Exit(1)
    }
    return v
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}