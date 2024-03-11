package functions

import (
	"context"
	"errors"
	"shopifyx/configs"
	"shopifyx/internal/database/interfaces"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	config configs.Config
	dbPool *pgxpool.Pool
}

func NewUser(dbPool *pgxpool.Pool, config configs.Config) *User {
	return &User{
		dbPool: dbPool,
		config: config,
	}
}

func (u *User) Register(ctx context.Context, usr interfaces.User) error {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), u.config.BcryptSalt)
	if err != nil {
		return err
	}

	sql := `
		INSERT INTO users (name, username, password) VALUES ($1, $2, $3)
	`

	_, err = conn.Exec(ctx, sql, usr.Name, usr.Username, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Login(ctx context.Context, username, password string) (interfaces.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return interfaces.User{}, err
	}
	defer conn.Release()

	var result interfaces.User

	err = conn.QueryRow(ctx, `SELECT id, name, username, password FROM users WHERE username = $1`, username).Scan(
		&result.Id, &result.Name, &result.Username, &result.Password,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, errors.New("user not found")
	}
	if err != nil {
		return result, err
	}

	// Compare the provided password with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		return result, errors.New("invalid password")
	}

	return result, nil
}
