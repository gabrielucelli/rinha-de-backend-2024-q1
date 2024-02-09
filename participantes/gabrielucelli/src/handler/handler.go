package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewHandler)

const defaultLimit = 20

type Handler struct {
	databaseConn *pgxpool.Pool
	validate     validator.Validate
}

func NewHandler(databaseConn *pgxpool.Pool, validate validator.Validate) Handler {
	return Handler{databaseConn, validate}
}
