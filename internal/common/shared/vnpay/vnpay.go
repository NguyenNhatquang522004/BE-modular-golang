// Package vnpay provides a complete, production-ready integration with the
// VNPay payment gateway for Vietnamese applications.
//
// Supported operations:
//   - CreatePaymentURL  – build a signed checkout URL and redirect the user.
//   - VerifyReturn      – validate the signature on the return-URL callback.
//   - ProcessIPN        – handle the server-to-server IPN notification.
//   - QueryTransaction  – query a transaction's status from VNPay's API.
//   - Refund            – request a full or partial refund.
//
// Usage:
//
//	cfg := &vnpay.Config{TmnCode: "XXX", HashSecret: "YYY", ReturnURL: "https://…"}
//	svc, err := vnpay.NewService(cfg)
//	resp, err := svc.CreatePaymentURL(ctx, req)
package vnpay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ─────────────────────────────────────────────
// Service Interface (Dependency-Inversion ready)
// ─────────────────────────────────────────────

// IVNPayService defines all operations exposed by the VNPay integration.
// Program to this interface, not the concrete *Service, so the implementation
// can be swapped or mocked in tests.
type IVNPayService interface {
	// CreatePaymentURL builds a signed VNPay checkout URL.
	// The caller should redirect the end-user's browser to the returned URL.
	CreatePaymentURL(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error)

	// VerifyReturn validates the signature on the query parameters that VNPay
	// appends to the ReturnURL after the user completes payment.
	// It returns the parsed domain params and a non-nil error if the signature
	// is invalid or the payment was unsuccessful.
	VerifyReturn(ctx context.Context, rawQuery string) (*PaymentReturnParams, error)

	// ProcessIPN handles the server-to-server Instant Payment Notification sent
	// by VNPay to your IPNURL.  Call this in your IPN handler; it verifies the
	// signature and parses the result.
	// The caller is responsible for updating the order status; this method only
	// handles signature verification and parameter parsing.
	ProcessIPN(ctx context.Context, rawQuery string) (*PaymentIPNParams, error)

	// QueryTransaction calls VNPay's Merchant API to check the current status
	// of a specific transaction.  Useful when the return URL or IPN is not
	// received reliably (e.g., user closed browser mid-flow).
	QueryTransaction(ctx context.Context, req QueryTransactionRequest) (*QueryTransactionResponse, error)

	// Refund requests a full or partial refund for a completed transaction.
	Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error)
}

// ─────────────────────────────────────────────
// Concrete Implementation
// ─────────────────────────────────────────────

// Service implements IVNPayService.
type Service struct {
	cfg        Config
	httpClient *http.Client
}

// NewService creates and validates a new VNPay Service.
// Returns ErrInvalidConfig if any required configuration field is missing.
func NewService(cfg Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("vnpay.NewService: %w", err)
	}
	return &Service{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// ─────────────────────────────────────────────
// CreatePaymentURL
// ─────────────────────────────────────────────

// CreatePaymentURL implements IVNPayService.
func (s *Service) CreatePaymentURL(_ context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// ── Input validation ──────────────────────────────────────────────────────
	if req.OrderRef == "" {
		return nil, ErrInvalidOrderRef
	}
	if req.AmountVND < MinimumPaymentAmount {
		return nil, ErrAmountTooSmall
	}
	if req.ClientIPAddr == "" {
		return nil, ErrMissingClientIP
	}
	if req.OrderInfo == "" {
		return nil, ErrMissingOrderInfo
	}

	// ── Defaults ──────────────────────────────────────────────────────────────
	if req.OrderType == "" {
		req.OrderType = OrderTypeOther
	}
	if req.Locale == "" {
		req.Locale = LocaleVN
	}

	now := time.Now()
	expireTime := now.Add(DefaultPaymentExpiryMinutes * time.Minute)
	if req.ExpireDate != nil {
		expireTime = *req.ExpireDate
	}

	// ── Build VNPay parameter map ─────────────────────────────────────────────
	// Wire format: amount is VND × 100 (no decimals, no separator).
	params := map[string]string{
		"vnp_Version":    s.cfg.Version,
		"vnp_Command":    CommandPay,
		"vnp_TmnCode":    s.cfg.TmnCode,
		"vnp_Locale":     req.Locale,
		"vnp_CurrCode":   CurrencyVND,
		"vnp_TxnRef":     req.OrderRef,
		"vnp_OrderInfo":  req.OrderInfo,
		"vnp_OrderType":  req.OrderType,
		"vnp_Amount":     fmt.Sprintf("%d", req.AmountVND*100),
		"vnp_ReturnUrl":  s.cfg.ReturnURL,
		"vnp_IpAddr":     req.ClientIPAddr,
		"vnp_CreateDate": now.Format("20060102150405"), // yyyyMMddHHmmss
		"vnp_ExpireDate": expireTime.Format("20060102150405"),
	}

	if req.BankCode != "" {
		params["vnp_BankCode"] = req.BankCode
	}

	// ── Sign and build URL ────────────────────────────────────────────────────
	payURL, err := BuildPaymentURL(s.cfg.PaymentURL, params, s.cfg.HashSecret)
	if err != nil {
		return nil, fmt.Errorf("vnpay.CreatePaymentURL: build URL: %w", err)
	}

	return &CreatePaymentResponse{
		PaymentURL: payURL,
		OrderRef:   req.OrderRef,
		AmountVND:  req.AmountVND,
		ExpireTime: expireTime,
	}, nil
}

// ─────────────────────────────────────────────
// VerifyReturn
// ─────────────────────────────────────────────

// VerifyReturn implements IVNPayService.
func (s *Service) VerifyReturn(_ context.Context, rawQuery string) (*PaymentReturnParams, error) {
	params, err := ParseReturnQuery(rawQuery)
	if err != nil {
		return nil, fmt.Errorf("vnpay.VerifyReturn: %w", err)
	}

	if !VerifySignature(params, s.cfg.HashSecret) {
		return nil, ErrInvalidSignature
	}

	result := &PaymentReturnParams{
		AmountRaw:         parseIntParam(params["vnp_Amount"]),
		BankCode:          params["vnp_BankCode"],
		BankTransNo:       params["vnp_BankTranNo"],
		CardType:          params["vnp_CardType"],
		OrderInfo:         params["vnp_OrderInfo"],
		PayDate:           params["vnp_PayDate"],
		ResponseCode:      params["vnp_ResponseCode"],
		TmnCode:           params["vnp_TmnCode"],
		TransactionNo:     params["vnp_TransactionNo"],
		TransactionStatus: params["vnp_TransactionStatus"],
		TxnRef:            params["vnp_TxnRef"],
		SecureHashType:    params["vnp_SecureHashType"],
		SecureHash:        params["vnp_SecureHash"],
	}

	if !result.IsSuccess() {
		return result, NewVNPayError(result.ResponseCode, ErrTransactionFailed)
	}

	return result, nil
}

// ─────────────────────────────────────────────
// ProcessIPN
// ─────────────────────────────────────────────

// ProcessIPN implements IVNPayService.
// It verifies the signature and parses the IPN parameters.
// The returned *PaymentIPNParams is populated whether or not the payment
// succeeded; the caller must inspect ResponseCode / TransactionStatus.
// If the signature is invalid, ErrInvalidSignature is returned.
func (s *Service) ProcessIPN(_ context.Context, rawQuery string) (*PaymentIPNParams, error) {
	params, err := ParseReturnQuery(rawQuery)
	if err != nil {
		return nil, fmt.Errorf("vnpay.ProcessIPN: %w", err)
	}

	if !VerifySignature(params, s.cfg.HashSecret) {
		return nil, ErrInvalidSignature
	}

	result := &PaymentIPNParams{
		AmountRaw:         parseIntParam(params["vnp_Amount"]),
		BankCode:          params["vnp_BankCode"],
		BankTransNo:       params["vnp_BankTranNo"],
		CardType:          params["vnp_CardType"],
		OrderInfo:         params["vnp_OrderInfo"],
		PayDate:           params["vnp_PayDate"],
		ResponseCode:      params["vnp_ResponseCode"],
		TmnCode:           params["vnp_TmnCode"],
		TransactionNo:     params["vnp_TransactionNo"],
		TransactionStatus: params["vnp_TransactionStatus"],
		TxnRef:            params["vnp_TxnRef"],
		SecureHashType:    params["vnp_SecureHashType"],
		SecureHash:        params["vnp_SecureHash"],
	}
	return result, nil
}

// ─────────────────────────────────────────────
// QueryTransaction
// ─────────────────────────────────────────────

// QueryTransaction implements IVNPayService.
func (s *Service) QueryTransaction(_ context.Context, req QueryTransactionRequest) (*QueryTransactionResponse, error) {
	if req.OrderRef == "" {
		return nil, ErrInvalidOrderRef
	}
	if req.ClientIPAddr == "" {
		return nil, ErrMissingClientIP
	}

	now := time.Now()
	params := map[string]string{
		"vnp_RequestId":  generateRequestID(now),
		"vnp_Version":    s.cfg.Version,
		"vnp_Command":    CommandQuery,
		"vnp_TmnCode":    s.cfg.TmnCode,
		"vnp_TxnRef":     req.OrderRef,
		"vnp_OrderInfo":  fmt.Sprintf("Truy van GD ma:%s", req.OrderRef),
		"vnp_TransDate":  req.TransactionDate.Format("20060102150405"),
		"vnp_CreateDate": now.Format("20060102150405"),
		"vnp_IpAddr":     req.ClientIPAddr,
	}

	// For the query / refund API the signature is built the same way.
	params["vnp_SecureHash"] = Sign(params, s.cfg.HashSecret)

	respBytes, err := s.postJSON(s.cfg.APIURL, params)
	if err != nil {
		return nil, fmt.Errorf("vnpay.QueryTransaction: %w", err)
	}

	var resp QueryTransactionResponse
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("vnpay.QueryTransaction: parse response: %w", err)
	}
	return &resp, nil
}

// ─────────────────────────────────────────────
// Refund
// ─────────────────────────────────────────────

// Refund implements IVNPayService.
func (s *Service) Refund(_ context.Context, req RefundRequest) (*RefundResponse, error) {
	if req.OrderRef == "" {
		return nil, ErrInvalidOrderRef
	}
	if req.AmountVND < 1000 {
		return nil, ErrAmountTooSmall
	}
	if req.ClientIPAddr == "" {
		return nil, ErrMissingClientIP
	}
	if req.TransType == "" {
		req.TransType = TransTypeFullRefund
	}

	now := time.Now()
	params := map[string]string{
		"vnp_RequestId":       generateRequestID(now),
		"vnp_Version":         s.cfg.Version,
		"vnp_Command":         CommandRefund,
		"vnp_TmnCode":         s.cfg.TmnCode,
		"vnp_TransactionType": req.TransType,
		"vnp_TxnRef":          req.OrderRef,
		"vnp_Amount":          fmt.Sprintf("%d", req.AmountVND*100),
		"vnp_OrderInfo":       fmt.Sprintf("Hoan tien GD ma:%s", req.OrderRef),
		"vnp_TransactionNo":   req.TransactionNo,
		"vnp_TransDate":       req.TransactionDate.Format("20060102150405"),
		"vnp_CreateDate":      now.Format("20060102150405"),
		"vnp_CreateBy":        req.CreatedBy,
		"vnp_IpAddr":          req.ClientIPAddr,
	}

	params["vnp_SecureHash"] = Sign(params, s.cfg.HashSecret)

	respBytes, err := s.postJSON(s.cfg.APIURL, params)
	if err != nil {
		return nil, fmt.Errorf("vnpay.Refund: %w", err)
	}

	var resp RefundResponse
	if err = json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("vnpay.Refund: parse response: %w", err)
	}
	return &resp, nil
}

// ─────────────────────────────────────────────
// Internal Helpers
// ─────────────────────────────────────────────

// postJSON sends a JSON POST request to the VNPay Merchant API and returns
// the raw response body.
func (s *Service) postJSON(endpoint string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	resp, err := s.httpClient.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: HTTP %d", ErrRequestFailed, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// generateRequestID creates a compact timestamp-based request ID.
func generateRequestID(t time.Time) string {
	return t.Format("20060102150405.000")
}

// parseIntParam safely converts a string to int64, returning 0 on failure.
func parseIntParam(s string) int64 {
	var v int64
	fmt.Sscanf(s, "%d", &v)
	return v
}
