package main

import (
	"context"
	"gabrielucelli/rinha-backend/src/handler"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var validate *validator.Validate
var dbConnection *pgxpool.Pool

func ProvideDatabaseConn() (*pgxpool.Pool, error) {
	con, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return con, nil
}

func ProvideFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})
	app.Use(recover.New())
	return app
}

func main() {
	fx.New(
		fx.Provide(ProvideDatabaseConn),
		fx.Provide(validator.New(validator.WithRequiredStructEnabled())),
		fx.Provide(ProvideFiberApp),
		fx.Provide(handler.Module),
	).Run()
}

func InitApp(lifecycle fx.Lifecycle, app *fiber.App, handler handler.Handler) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				app.Post("/clientes/:id/transacoes", handler.CreateTransactionHandler)
				app.Get("/clientes/:id/extrato", handler.GetExtractHandler)
				log.Print(app.Listen(":" + os.Getenv("APP_PORT")))
				return nil
			},
			OnStop: func(context.Context) error {
				return nil
			},
		},
	)
}
