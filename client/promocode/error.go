package promocode

import "errors"

var (
	ErrPromoCodeNotFound  = errors.New("promo code not found")
	ErrInvalidEffectValue = errors.New("invalid effect value")
)
