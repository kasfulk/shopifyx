package functions

import (
	"context"
	"errors"
	"fmt"
	"shopifyx/db/interfaces"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	dbPool *pgxpool.Pool
}

func NewProductFn(dbPool *pgxpool.Pool) *Product {
	return &Product{
		dbPool: dbPool,
	}
}

func (p *Product) Buy(ctx context.Context, payment interfaces.Payment) (interfaces.Payment, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return interfaces.Payment{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	bankAccount := interfaces.BankPayment{}

	err = conn.QueryRow(ctx, "select user_id, bank_name, bank_account_name, bank_account_number from banks where id = $1", payment.BankAccountId).Scan(
		&bankAccount.UserId, &bankAccount.BankName, &bankAccount.BankAccountName, &bankAccount.BankAccountNumber,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return interfaces.Payment{}, ErrNoRow
	}

	if err != nil {
		return interfaces.Payment{}, fmt.Errorf("failed get bank when do payment: %v", err)
	}

	user := interfaces.UserPayment{}

	err = conn.QueryRow(ctx, "select id, username, name from users where id = $1", bankAccount.UserId).Scan(
		&user.UserId, &user.BuyerUsername, &user.BuyerName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return interfaces.Payment{}, ErrNoRow
	}

	if err != nil {
		return interfaces.Payment{}, fmt.Errorf("failed get user when do payment: %v", err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return interfaces.Payment{}, fmt.Errorf("failed start transaction: %v", err)
	}

	product := interfaces.ProductPayment{}

	err = tx.QueryRow(ctx, "select id, name, image_url, stock, price from products where id = $1 for update", payment.ProductId).Scan(
		&product.Id, &product.Name, &product.ImageUrl, &product.Qty, &product.Price,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		tx.Rollback(ctx)
		return interfaces.Payment{}, ErrNoRow
	}

	if err != nil {
		tx.Rollback(ctx)
		return interfaces.Payment{}, fmt.Errorf("failed get product when do payment: %v", err)
	}

	if product.Qty < payment.Qty {
		tx.Rollback(ctx)
		return interfaces.Payment{}, ErrInsuficientQty
	}

	_, err = tx.Exec(ctx, "update products set stock = stock - $1 where id = $2", payment.Qty, payment.ProductId)
	if err != nil {
		tx.Rollback(ctx)
		return interfaces.Payment{}, fmt.Errorf("failed update product stock: %v", err)
	}

	err = tx.QueryRow(ctx, `INSERT INTO payments (product_id, product_name, product_image_url, product_qty, product_price, user_id, buyer_username, buyer_name, bank_name, bank_account_name, bank_account_number, payment_proof_image_url) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	) RETURNING id, created_at, updated_at`,
		product.Id, product.Name, product.ImageUrl, payment.Qty, product.Price, user.UserId, user.BuyerUsername, user.BuyerName, bankAccount.BankName, bankAccount.BankAccountName, bankAccount.BankAccountNumber, payment.PaymentProofImageUrl,
	).Scan(&payment.Id, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		tx.Rollback(ctx)
		return interfaces.Payment{}, fmt.Errorf("failed create payment: %v", err)
	}

	tx.Commit(ctx)

	return payment, nil
}
