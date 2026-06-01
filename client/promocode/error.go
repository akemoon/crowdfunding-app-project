package promocode

import "errors"

var (
	ErrPromoCodeNotFound  = errors.New("promo code not found")
	ErrInvalidCodeFormat  = errors.New("invalid promo code format")
	ErrInvalidEffectValue = errors.New("invalid effect value")
)
