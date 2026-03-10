package promocode

import "errors"

var (
	ErrPromoCodeNotFound  = errors.New("promo code not found")
	ErrPromoCodeUsed      = errors.New("promo code already used")
	ErrPromoCodeForbidden = errors.New("promo code forbidden")
	ErrInvalidEffectValue = errors.New("invalid effect value")
)
