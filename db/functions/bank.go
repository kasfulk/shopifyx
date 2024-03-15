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

	_, err = conn.Exec(ctx, `INSERT INTO banks (user_id, bank_name, bank_account_number, bank_account_name) values($1, $2, $3, $4)`,
		bnk.UserId, bnk.BankName, bnk.BankAccountNumber, bnk.BankAccountName,
	)

	return err
}

func (b *Bank) Get(ctx context.Context, userId string) ([]entity.Bank, error) {
	result := []entity.Bank{}

	conn, err := b.dbPool.Acquire(ctx)
	if err != nil {
		return result, fmt.Errorf("failed acquire connection from db pool: %v", err)
	}

	defer conn.Release()

	rows, err := conn.Query(ctx, `select id, bank_name, bank_account_number, bank_account_name from banks where user_id = $1`, userId)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var bnk entity.Bank

		err := rows.Scan(&bnk.Id, &bnk.BankName, &bnk.BankAccountNumber, &bnk.BankAccountName)
		if err != nil {
			return []entity.Bank{}, err
		}

		result = append(result, bnk)
	}

	return result, nil
}
