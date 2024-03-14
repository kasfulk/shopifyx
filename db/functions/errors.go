package functions

import "errors"

var (
	ErrNoRow          = errors.New("data not found")
	ErrInsuficientQty = errors.New("insuficient quantity")
)

// product section

var ErrProductNameDuplicate = errors.New("product name already exists")
