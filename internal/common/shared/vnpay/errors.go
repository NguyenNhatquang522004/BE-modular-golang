package vnpay

import (
	"errors"
	"fmt"
)

// ─────────────────────────────────────────────
// Sentinel Errors
// ─────────────────────────────────────────────

var (
	// ErrInvalidSignature is returned when the HMAC-SHA512 signature does not match.
	ErrInvalidSignature = errors.New("vnpay: invalid signature")

	// ErrTransactionFailed is returned when the payment transaction is not successful.
	ErrTransactionFailed = errors.New("vnpay: transaction failed")

	// ErrTransactionPending is returned when the transaction is still being processed.
	ErrTransactionPending = errors.New("vnpay: transaction pending")

	// ErrInvalidConfig is returned when required VNPay config values are missing.
	ErrInvalidConfig = errors.New("vnpay: invalid configuration")

	// ErrInvalidOrderRef is returned when the order reference is empty or invalid.
	ErrInvalidOrderRef = errors.New("vnpay: order reference must not be empty")

	// ErrAmountTooSmall is returned when the payment amount is below the minimum.
	ErrAmountTooSmall = fmt.Errorf("vnpay: payment amount must be at least %d VND", MinimumPaymentAmount)

	// ErrMissingClientIP is returned when the client IP address is not provided.
	ErrMissingClientIP = errors.New("vnpay: client IP address is required")

	// ErrMissingOrderInfo is returned when the order info/description is empty.
	ErrMissingOrderInfo = errors.New("vnpay: order info must not be empty")

	// ErrRequestFailed is returned when the HTTP request to VNPay API fails.
	ErrRequestFailed = errors.New("vnpay: API request failed")
)

// ─────────────────────────────────────────────
// Structured Error Type
// ─────────────────────────────────────────────

// VNPayError represents an error returned by the VNPay gateway, including
// the VNPay response code and human-readable message.
type VNPayError struct {
	// Code is the vnp_ResponseCode or vnp_TransactionStatus value from VNPay.
	Code string

	// Message is the Vietnamese description of the error code.
	Message string

	// Err is the optional underlying Go error.
	Err error
}

// Error implements the error interface.
func (e *VNPayError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("vnpay error [code=%s]: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("vnpay error [code=%s]: %s", e.Code, e.Message)
}

// Unwrap allows errors.Is / errors.As to traverse the chain.
func (e *VNPayError) Unwrap() error {
	return e.Err
}

// ─────────────────────────────────────────────
// Constructors
// ─────────────────────────────────────────────

// NewVNPayError creates a *VNPayError using the standard VNPay response code table.
// If the code is not found in ResponseCodeMessages, a generic message is used.
func NewVNPayError(code string, err error) *VNPayError {
	msg, ok := ResponseCodeMessages[code]
	if !ok {
		msg = ResponseCodeMessages[ResponseCodeSystemError]
	}
	return &VNPayError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

// NewVNPayErrorWithMsg creates a *VNPayError with an explicit message.
func NewVNPayErrorWithMsg(code, message string, err error) *VNPayError {
	return &VNPayError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

// IsSuccessCode returns true if the given VNPay response code means success.
func IsSuccessCode(code string) bool {
	return code == ResponseCodeSuccess
}

// MessageForCode returns the Vietnamese description for a VNPay response code.
// If the code is unknown, the generic system-error message is returned.
func MessageForCode(code string) string {
	if msg, ok := ResponseCodeMessages[code]; ok {
		return msg
	}
	return ResponseCodeMessages[ResponseCodeSystemError]
}
