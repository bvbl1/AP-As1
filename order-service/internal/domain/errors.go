package domain

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrNotEnoughStock = errors.New("not enough stock")
)
