package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewHandler)

type Handler struct {
	databaseConn *pgxpool.Pool
	validate     *validator.Validate
}

func NewHandler(databaseConn *pgxpool.Pool) Handler {
	return Handler{
		databaseConn: databaseConn,
		validate:     validator.New(validator.WithRequiredStructEnabled()),
	}
}
