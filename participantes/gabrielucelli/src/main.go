package main

import (
	"context"
	"gabrielucelli/rinha-backend/src/handler"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func ProvideDatabaseConn() (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
}

func ProvideFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})
	if len(os.Getenv("DEBUG_ENABLED")) != 0 {
		app.Use(recover.New())
		app.Use(logger.New())
	}
	return app
}

func main() {
	app := fx.New(
		fx.Provide(ProvideDatabaseConn),
		fx.Provide(ProvideFiberApp),
		fx.Options(handler.Module),
		fx.Invoke(InitApp),
	)
	app.Run()
}

func InitApp(lifecycle fx.Lifecycle, app *fiber.App, handler handler.Handler) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				app.Post("/clientes/:id/transacoes", handler.CreateTransactionHandler)
				app.Get("/clientes/:id/extrato", handler.GetExtractHandler)
				go app.Listen(":" + os.Getenv("APP_PORT"))
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}
