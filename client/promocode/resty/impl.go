package resty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akemoon/crowdfunding-app-project/client/promocode"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const userIDHeader = "X-User-ID"

type PromoClient struct {
	client *resty.Client
}

func NewPromoClient(baseURL string) *PromoClient {
	return &PromoClient{
		client: resty.New().SetBaseURL(baseURL),
	}
}

type usePromoCodeReq struct {
	ExpectedType string `json:"expectedType"`
}

type usePromoCodeResp struct {
	EffectValue int `json:"effectValue"`
}

func (c *PromoClient) UsePromoCode(ctx context.Context, userID uuid.UUID, code string, expectedType string) (int, error) {
	var data usePromoCodeResp

	resp, err := c.client.R().
		SetContext(ctx).
		SetPathParam("code", code).
		SetHeader(userIDHeader, userID.String()).
		SetBody(usePromoCodeReq{ExpectedType: expectedType}).
		SetResult(&data).
		Post("/promocodes/{code}/use")
	if err != nil {
		return 0, fmt.Errorf("use promo code: %w", err)
	}

	if resp.StatusCode() == http.StatusBadRequest {
		return 0, promocode.ErrInvalidCodeFormat
	}

	if resp.StatusCode() == http.StatusNotFound {
		return 0, promocode.ErrPromoCodeNotFound
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("use promo code: unexpected status: %d", resp.StatusCode())
	}

	if data.EffectValue <= 0 {
		return 0, promocode.ErrInvalidEffectValue
	}

	return data.EffectValue, nil
}
