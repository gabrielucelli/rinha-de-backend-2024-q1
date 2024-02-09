package model

import "time"

// POST /clientes/[id]/transacoes

// request
type CreateTransactionRequest struct {
	Amount      int    `json:"valor" validate:"required,gt=0"`
	Description string `json:"descricao" validate:"required,min=1,max=10"`
	Type        string `json:"tipo" validate:"required,oneof=c d"`
}

// response
type CreateTransactionResponse struct {
	AccountLimit   int `json:"limite"`
	AccountBalance int `json:"saldo"`
}

// GET /clientes/[id]/extrato

// response
type BankStatementResponse struct {
	Balance          Balance       `json:"saldo"`
	LastTransactions []Transaction `json:"ultimas_transacoes"`
}

type Balance struct {
	AccountBalance int       `json:"total"`
	AccountLimit   int       `json:"limite"`
	BalanceDate    time.Time `json:"data_extrato"`
}

type Transaction struct {
	Amount      int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}
