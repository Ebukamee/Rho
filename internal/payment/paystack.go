package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type PaystackProvider struct {
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

func NewPaystackProvider(secretKey string) *PaystackProvider {
	return &PaystackProvider{
		secretKey:  secretKey,
		baseURL:    "https://api.paystack.co",
		httpClient: &http.Client{},
	}
}

type paystackInitializeRequest struct {
	Email    string `json:"email"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type paystackInitializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type paystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status string `json:"status"`
	} `json:"data"`
}

func (p *PaystackProvider) Initialize(
	ctx context.Context,
	payment *Payment,
	email string,
) (*PaymentInitialization, error) {

	if strings.TrimSpace(p.secretKey) == "" {
		return nil, errors.New("paystack secret key is not configured")
	}

	requestBody := paystackInitializeRequest{
		Email:    email,
		Amount:   payment.Amount,
		Currency: payment.Currency,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal paystack request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseURL+"/transaction/initialize",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create paystack request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send paystack request: %w", err)
	}
	defer resp.Body.Close()

	var result paystackInitializeResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode paystack response: %w", err)
	}

	if !result.Status {
		return nil, fmt.Errorf(
			"paystack initialization failed: %s",
			result.Message,
		)
	}

	return &PaymentInitialization{
		Provider:         "paystack",
		ProviderRef:      result.Data.Reference,
		AuthorizationURL: result.Data.AuthorizationURL,
	}, nil
}

func (p *PaystackProvider) Verify(
	ctx context.Context,
	providerRef string,
) (*PaymentVerification, error) {

	if strings.TrimSpace(p.secretKey) == "" {
		return nil, errors.New("paystack secret key is not configured")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		p.baseURL+"/transaction/verify/"+providerRef,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create paystack verification request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.secretKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send paystack verification request: %w", err)
	}
	defer resp.Body.Close()

	var result paystackVerifyResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode paystack verification response: %w", err)
	}

	if !result.Status {
		return nil, fmt.Errorf(
			"paystack verification failed: %s",
			result.Message,
		)
	}

	var status PaymentStatus

	switch result.Data.Status {
	case "success":
		status = PaymentSucceeded

	case "failed":
		status = PaymentFailed

	default:
		status = PaymentPending
	}

	return &PaymentVerification{
		ProviderRef: providerRef,
		Status:      status,
	}, nil
}
