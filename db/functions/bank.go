package functions

import (
	"context"
	"fmt"
	"shopifyx/db/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Bank struct {
	dbPool *pgxpool.Pool
}

func NewBank(dbPool *pgxpool.Pool) *Bank {
	return &Bank{
		dbPool: dbPool,
	}
}

func (b *Bank) Create(ctx context.Context, bnk entity.Bank) error {
	conn, err := b.dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed acquire connection from db pool: %v", err)
	}

	defer conn.Release()

	_, err = b.dbPool.Exec(ctx, `INSERT INTO banks (user_id, bank_name, bank_account_number, bank_account_name) values($1, $2, $3, $4)`,
		bnk.UserId, bnk.BankName, bnk.BankAccountNumber, bnk.BankAccountName,
	)

	return err
}
