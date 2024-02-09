package handler

import (
	"context"
	"errors"
	"gabrielucelli/rinha-backend/src/model"

	"github.com/jackc/pgx/v5"
)

var cachedClients = make(map[int]*model.Client)

func (h *Handler) GetClient(ctx context.Context, clientId int) (*model.Client, error) {

	cachedClient, ok := cachedClients[clientId]
	if ok {
		return cachedClient, nil
	}

	rows, err := h.databaseConn.Query(ctx, "SELECT * FROM clients WHERE id = $1", clientId)
	if err != nil {
		return nil, err
	}

	client, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Client])

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		cachedClients[client.Id] = nil
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	cachedClients[client.Id] = &client
	return &client, nil
}
