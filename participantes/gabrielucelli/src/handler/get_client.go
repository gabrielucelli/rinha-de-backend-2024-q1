package handler

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Client struct {
	Id           int
	AccountLimit int
	Balance      int
}

var cachedClients = make(map[int]*Client)

func (h *Handler) GetClient(ctx context.Context, clientId int) (*Client, error) {
	cachedClient, ok := cachedClients[clientId]
	if ok {
		return cachedClient, nil
	}

	rows, _ := h.databaseConn.Query(ctx, "SELECT * FROM clients WHERE id = $1", clientId)
	client, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Client])

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
