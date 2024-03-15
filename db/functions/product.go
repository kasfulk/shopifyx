package functions

import (
	"context"
	"errors"
	"fmt"
	"shopifyx/db/entity"
	"strings"

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

func (p *Product) FindAll(ctx context.Context, filter entity.FilterGetProducts, userID int) ([]entity.Product, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `SELECT id, user_id, name, price, image_url, stock, condition, tags, is_purchaseable FROM products`

	whereSQL := []string{}
	if filter.UserOnly {
		whereSQL = append(whereSQL, " user_id = "+fmt.Sprintf("%d", userID))
	}
	if filter.Limit > 0 {
		sql += " LIMIT " + fmt.Sprintf("%d", filter.Limit)
	}
	if filter.Offset > 0 {
		sql += " OFFSET " + fmt.Sprintf("%d", filter.Offset)
	}
	if filter.Tags != nil && len(filter.Tags) > 0 {
		tagsSQL := []string{}
		for _, tag := range filter.Tags {
			tagsSQL = append(tagsSQL, fmt.Sprintf("'%s'", tag))
		}

		whereSQL = append(whereSQL, " tags IN ("+strings.Join(tagsSQL, ",")+")")
	}

	if len(whereSQL) > 0 {
		sql += " WHERE " + strings.Join(whereSQL, " AND ")
	}

	fmt.Println(sql)

	rows, err := conn.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed get products: %v", err)
	}

	defer rows.Close()

	products := []entity.Product{}

	for rows.Next() {
		product := entity.Product{}
		err := rows.Scan(&product.ID, &product.UserID, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, &product.Tags, &product.IsPurchaseable)
		if err != nil {
			return nil, fmt.Errorf("failed scan products: %v", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *Product) Buy(ctx context.Context, payment entity.Payment) (entity.Payment, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Payment{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	bankAccount := entity.BankPayment{}

	err = conn.QueryRow(ctx, "select user_id, bank_name, bank_account_name, bank_account_number from banks where id = $1", payment.BankAccountId).Scan(
		&bankAccount.UserId, &bankAccount.BankName, &bankAccount.BankAccountName, &bankAccount.BankAccountNumber,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Payment{}, ErrNoRow
	}

	if err != nil {
		return entity.Payment{}, fmt.Errorf("failed get bank when do payment: %v", err)
	}

	user := entity.UserPayment{}

	err = conn.QueryRow(ctx, "select id, username, name from users where id = $1", bankAccount.UserId).Scan(
		&user.UserId, &user.BuyerUsername, &user.BuyerName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Payment{}, ErrNoRow
	}

	if err != nil {
		return entity.Payment{}, fmt.Errorf("failed get user when do payment: %v", err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return entity.Payment{}, fmt.Errorf("failed start transaction: %v", err)
	}

	product := entity.ProductPayment{}

	err = tx.QueryRow(ctx, "select id, name, image_url, stock, price from products where id = $1 for update", payment.ProductId).Scan(
		&product.Id, &product.Name, &product.ImageUrl, &product.Qty, &product.Price,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		tx.Rollback(ctx)
		return entity.Payment{}, ErrNoRow
	}

	if err != nil {
		tx.Rollback(ctx)
		return entity.Payment{}, fmt.Errorf("failed get product when do payment: %v", err)
	}

	if product.Qty < payment.Qty {
		tx.Rollback(ctx)
		return entity.Payment{}, ErrInsuficientQty
	}

	_, err = tx.Exec(ctx, "update products set stock = stock - $1 where id = $2", payment.Qty, payment.ProductId)
	if err != nil {
		tx.Rollback(ctx)
		return entity.Payment{}, fmt.Errorf("failed update product stock: %v", err)
	}

	err = tx.QueryRow(ctx, `INSERT INTO payments (product_id, product_name, product_image_url, product_qty, product_price, user_id, buyer_username, buyer_name, bank_name, bank_account_name, bank_account_number, payment_proof_image_url) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	) RETURNING id, created_at, updated_at`,
		product.Id, product.Name, product.ImageUrl, payment.Qty, product.Price, user.UserId, user.BuyerUsername, user.BuyerName, bankAccount.BankName, bankAccount.BankAccountName, bankAccount.BankAccountNumber, payment.PaymentProofImageUrl,
	).Scan(&payment.Id, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		tx.Rollback(ctx)
		return entity.Payment{}, fmt.Errorf("failed create payment: %v", err)
	}

	tx.Commit(ctx)

	return payment, nil
}

func (p *Product) Add(ctx context.Context, product entity.Product) (entity.Product, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `
		insert into products (user_id, name, price, image_url, stock, condition, tags, is_purchaseable) 
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = conn.Exec(ctx, sql,
		product.UserID,
		product.Name,
		product.Price,
		product.ImageUrl,
		product.Stock,
		product.Condition,
		product.Tags,
		product.IsPurchaseable)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return entity.Product{}, ErrProductNameDuplicate
		}
		return entity.Product{}, fmt.Errorf("failed insert product: %v", err)
	}

	var result entity.Product

	err = conn.QueryRow(ctx, `SELECT id FROM products WHERE name = $1`, product.Name).Scan(&result.ID)

	if err != nil {
		return entity.Product{}, err
	}
	product.ID = result.ID

	return product, nil
}

func (p *Product) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `
		update products set name = $1, price = $2, image_url = $3, stock = $4, condition = $5, tags = $6, is_purchaseable = $7 
		where id = $8 and user_id = $9
	`

	_, err = conn.Exec(ctx, sql,
		product.Name,
		product.Price,
		product.ImageUrl,
		product.Stock,
		product.Condition,
		product.Tags,
		product.IsPurchaseable,
		product.ID,
		product.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return entity.Product{}, ErrNoRow
		} else if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return entity.Product{}, ErrProductNameDuplicate
		}
		return entity.Product{}, fmt.Errorf("failed update product: %v", err)
	}

	return product, nil
}

func (p *Product) FindByID(ctx context.Context, productID int) (entity.Product, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	var product entity.Product

	err = conn.QueryRow(ctx, `SELECT id, user_id, name, price, image_url, stock, condition, tags, is_purchaseable FROM products WHERE id = $1`, productID).Scan(
		&product.ID, &product.UserID, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, &product.Tags, &product.IsPurchaseable,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Product{}, ErrNoRow
	}

	sql := `
		update products set stock = $1 where id = $2 AND user_id = $3
	`

	_, err = conn.Exec(ctx, sql, product.Stock, product.ID, product.UserID)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed update product stock: %v", err)
	}
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed get product: %v", err)
	}

	return product, nil
}

func (p *Product) FindByIDUser(ctx context.Context, productID int, userID int) (entity.Product, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	var product entity.Product

	err = conn.QueryRow(ctx, `SELECT id, user_id, name, price, image_url, stock, condition, tags, is_purchaseable FROM products WHERE id = $1 AND UserID = $2`, productID, userID).Scan(
		&product.ID, &product.UserID, &product.Name, &product.Price, &product.ImageUrl, &product.Stock, &product.Condition, &product.Tags, &product.IsPurchaseable,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Product{}, ErrNoRow
	}

	sql := `
		update products set stock = $1 where id = $2 AND user_id = $3
	`

	_, err = conn.Exec(ctx, sql, product.Stock, product.ID, product.UserID)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed update product stock: %v", err)
	}
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed get product: %v", err)
	}

	return product, nil
}

func (p *Product) UpdateStock(ctx context.Context, product entity.Product, userID int) (entity.Product, error) {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed acquire db connection from pool: %v", err)
	}
	defer conn.Release()

	sqlCheck := `SELECT id FROM products WHERE id = $1 AND user_id = $2`

	err = conn.QueryRow(ctx, sqlCheck, product.ID, userID).Scan(&product.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Product{}, ErrNoRow
	}

	sql := `
		update products set stock = $1 where id = $2 AND user_id = $3
	`

	_, err = conn.Exec(ctx, sql, product.Stock, product.ID, userID)
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed update product stock: %v", err)
	}
	return product, nil
}

func (p *Product) DeleteByID(ctx context.Context, productID int) error {
	conn, err := p.dbPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed acquire db connection from pool: %v", err)
	}

	defer conn.Release()

	sql := `delete from products where id = $1`
	_, err = conn.Exec(ctx, sql, productID)
	if err != nil {
		return fmt.Errorf("failed delete product: %v", err)
	}

	return nil
}
