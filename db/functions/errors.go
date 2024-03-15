package functions

import "errors"

var (
	ErrNoRow                = errors.New("data not found")
	ErrInsuficientQty       = errors.New("insuficient quantity")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrProductNameDuplicate = errors.New("product name already exists")
)
