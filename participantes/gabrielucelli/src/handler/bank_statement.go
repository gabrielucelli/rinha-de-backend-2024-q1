package handler

import (
	"context"
	"errors"
	"gabrielucelli/rinha-backend/src/model"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) GetExtractHandler(ctx *fiber.Ctx) error {
	clientId, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.SendStatus(422)
	}

	bankExtract, err := h.getBankExtract(ctx.Context(), clientId)
	if err != nil {
		if err.Error() == "USER_DONT_EXISTS" {
			return ctx.SendStatus(404)
		} else {
			return ctx.SendStatus(500)
		}
	}

	return ctx.JSON(bankExtract)
}

func (h *Handler) getBankExtract(ctx context.Context, clientId int) (model.BankStatementResponse, error) {

	group := errgroup.Group{}
	var balance model.Balance
	var lastTransactions []model.Transaction

	group.Go(func() error {
		balanceReceiver, err := h.getBalance(ctx, clientId)
		balance = balanceReceiver
		return err
	})

	group.Go(func() error {
		lastTransactionsReceiver, err := h.getLastTransactions(ctx, clientId)
		lastTransactions = lastTransactionsReceiver
		return err
	})

	err := group.Wait()
	if err != nil {
		return model.BankStatementResponse{}, err
	}

	result := model.BankStatementResponse{
		Balance:          balance,
		LastTransactions: lastTransactions,
	}

	return result, nil
}

func (h *Handler) getBalance(ctx context.Context, clientId int) (model.Balance, error) {
	rows, err := h.databaseConn.Query(ctx, "SELECT balance, account_limit, now() FROM clients WHERE id = $1", clientId)
	if err != nil {
		return model.Balance{}, err
	}
	balance, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[model.Balance])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return model.Balance{}, errors.New("USER_DONT_EXISTS")
	}
	if err != nil {
		return model.Balance{}, err
	}
	return balance, nil
}

func (h *Handler) getLastTransactions(ctx context.Context, clientId int) ([]model.Transaction, error) {
	rows, err := h.databaseConn.Query(ctx, "SELECT value, operation, description, created_at FROM transactions WHERE client_id = $1 ORDER BY id DESC LIMIT 10", clientId)

	if err != nil {
		return nil, err
	}
	lastTransactions, err := pgx.CollectRows(rows, pgx.RowToStructByPos[model.Transaction])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("USER_DONT_EXISTS")
	}
	if err != nil {
		return nil, err
	}
	return lastTransactions, err
}
