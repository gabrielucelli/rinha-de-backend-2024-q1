package handler

import (
	"context"
	"gabrielucelli/rinha-backend/src/model"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) GetExtractHandler(ctx *fiber.Ctx) error {
	clientId, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(422)
	}

	client, err := h.GetClient(ctx.Context(), clientId)
	if err != nil {
		return ctx.SendStatus(500)
	}
	if client == nil {
		return ctx.SendStatus(404)
	}

	bankExtract, err := h.getBankExtract(ctx.Context(), clientId)
	if err != nil {
		log.Fatal(err.Error())
		return ctx.SendStatus(500)
	}

	return ctx.JSON(bankExtract)
}

func (h *Handler) getBankExtract(ctx context.Context, clientId int) (model.BankStatementResponse, error) {
	rows, _ := h.databaseConn.Query(ctx, "SELECT balance, account_limit, now() FROM clients WHERE id = $1", clientId)
	balance, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.Balance])
	if err != nil {
		return model.BankStatementResponse{}, err
	}

	rows, _ = h.databaseConn.Query(ctx, "SELECT value, operation, description, created_at FROM transactions WHERE client_id = $1 ORDER BY id DESC LIMIT 10", clientId)
	lastTransactions, err := pgx.CollectRows(rows, pgx.RowToStructByPos[model.Transaction])
	if err != nil {
		return model.BankStatementResponse{}, err
	}

	result := model.BankStatementResponse{
		Balance:          balance,
		LastTransactions: lastTransactions,
	}

	return result, nil
}
