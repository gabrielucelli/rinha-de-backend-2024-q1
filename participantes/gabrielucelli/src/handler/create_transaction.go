package handler

import (
	"context"
	"errors"
	"gabrielucelli/rinha-backend/src/model"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) CreateTransactionHandler(ctx *fiber.Ctx) error {

	clientId, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(422)
	}

	var createTransactionRequest model.CreateTransactionRequest
	err = ctx.BodyParser(&createTransactionRequest)
	if err != nil {
		return ctx.SendStatus(422)
	}

	err = h.validate.Struct(createTransactionRequest)
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

	transactionResult, err := h.createTransaction(ctx.Context(), clientId, createTransactionRequest)
	if err != nil {
		if err.Error() == "LIMIT_EXCEEDED" {
			return ctx.SendStatus(422)
		} else {
			return ctx.SendStatus(500)
		}
	}

	return ctx.JSON(transactionResult)
}

func (h *Handler) createTransaction(ctx context.Context, clientId int, transactionData model.CreateTransactionRequest) (model.CreateTransactionResponse, error) {
	tx, err := h.databaseConn.Begin(ctx)
	if err != nil {
		return model.CreateTransactionResponse{}, err
	}

	defer tx.Rollback(ctx)

	var accountLimit int
	var accountBalance int

	err = tx.QueryRow(ctx, "SELECT account_limit, balance FROM clients WHERE id = $1 FOR UPDATE", clientId).Scan(&accountLimit, &accountBalance)
	if err != nil {
		return model.CreateTransactionResponse{}, err
	}

	var newAccountBalance int

	if transactionData.Type == "d" {
		newAccountBalance = accountBalance - transactionData.Amount
	} else {
		newAccountBalance = accountBalance + transactionData.Amount
	}

	if (accountLimit + newAccountBalance) < 0 {
		return model.CreateTransactionResponse{}, errors.New("LIMIT_EXCEEDED")
	}

	batch := &pgx.Batch{}
	batch.Queue("INSERT INTO transactions(client_id,value,operation,description) values ($1, $2, $3, $4)", clientId, transactionData.Amount, transactionData.Type, transactionData.Description)
	batch.Queue("UPDATE clients SET balance = $1 WHERE id = $2", newAccountBalance, clientId)
	br := tx.SendBatch(ctx, batch)
	_, err = br.Exec()

	if err != nil {
		return model.CreateTransactionResponse{}, err
	}

	err = br.Close()
	if err != nil {
		return model.CreateTransactionResponse{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.CreateTransactionResponse{}, err
	}

	result := model.CreateTransactionResponse{
		AccountBalance: newAccountBalance,
		AccountLimit:   accountLimit,
	}

	return result, nil
}
