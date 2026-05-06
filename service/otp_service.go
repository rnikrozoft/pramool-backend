package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rnikrozoft/pramool-core/model"
	"go.uber.org/zap"
)

type OTPService interface {
	Request(ctx context.Context, msisdn string) (*model.OTPResponse, error)
	Verify(ctx context.Context, token, pin string) error
}

type otpService struct {
	logger *zap.Logger

	httpClient     *http.Client
	addressRequest string
	addressVerify  string
	apiKey         string
	apiSecret      string
}

func NewOTPService(logger *zap.Logger, addressRequest, addressVerify, apiKey, apiSecret string, httpClient *http.Client) OTPService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 45 * time.Second}
	}
	return otpService{
		logger:         logger,
		httpClient:     httpClient,
		addressRequest: addressRequest,
		addressVerify:  addressVerify,
		apiKey:         apiKey,
		apiSecret:      apiSecret,
	}
}

func (s otpService) Request(ctx context.Context, msisdn string) (*model.OTPResponse, error) {
	str := fmt.Sprintf("key=%s&secret=%s&msisdn=%s", s.apiKey, s.apiSecret, msisdn)
	payload := strings.NewReader(str)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addressRequest, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		s.logger.Warn("unexpected OTP response",
			zap.Int("status_code", res.StatusCode),
			zap.String("status", res.Status),
		)
		return nil, errors.New("cannot request otp")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var otpRes model.OTPResponse
	if err := json.Unmarshal(body, &otpRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &otpRes, nil
}

func (s otpService) Verify(ctx context.Context, token, pin string) error {
	str := fmt.Sprintf("key=%s&secret=%s&token=%s&pin=%s", s.apiKey, s.apiSecret, token, pin)
	payload := strings.NewReader(str)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.addressVerify, payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		return errors.New("cannot verify otp")
	}
	return nil
}
