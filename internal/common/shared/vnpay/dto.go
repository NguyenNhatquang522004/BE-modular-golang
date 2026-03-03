package vnpay

import "time"

// ─────────────────────────────────────────────
// Inbound DTOs  (caller → VNPay Service)
// ─────────────────────────────────────────────

// CreatePaymentRequest is the input contract for generating a VNPay payment URL.
// All monetary values are in VND (Vietnamese Dong).
type CreatePaymentRequest struct {
	// OrderRef is YOUR system's unique order reference.
	// VNPay returns it unchanged as vnp_TxnRef.
	// Max 50 ASCII printable chars; must be unique per transaction attempt.
	OrderRef string `json:"order_ref" validate:"required,max=50"`

	// AmountVND is the amount to charge in whole Vietnamese Dong.
	// Must be ≥ MinimumPaymentAmount (5 000 VND).
	AmountVND int64 `json:"amount_vnd" validate:"required,min=5000"`

	// OrderInfo is a short human-readable description shown on VNPay's page.
	// Max 255 chars, no special characters like |, ;, &.
	OrderInfo string `json:"order_info" validate:"required,max=255"`

	// OrderType categorises the transaction.  Defaults to OrderTypeOther.
	OrderType string `json:"order_type"`

	// Locale controls the language of VNPay's checkout page.
	// "vn" (Vietnamese, default) or "en" (English).
	Locale string `json:"locale"`

	// BankCode optionally pre-selects a bank or payment method
	// (e.g. "NCB", "VNPAYQR").  Leave empty to show the full bank picker.
	BankCode string `json:"bank_code"`

	// ClientIPAddr is the end-user's IP address (IPv4 or IPv6).
	// Required by VNPay; used for fraud detection.
	ClientIPAddr string `json:"client_ip_addr" validate:"required"`

	// ExpireDate overrides the default payment URL expiry time.
	// Defaults to now + DefaultPaymentExpiryMinutes.
	ExpireDate *time.Time `json:"expire_date,omitempty"`

	// ── Optional Application Context ─────────────────────────────────────────

	// UserID optionally associates this transaction with an app user.
	UserID string `json:"user_id,omitempty"`

	// ReferenceID optionally links to the business object being paid for
	// (e.g., ad campaign ID, subscription plan ID).
	ReferenceID string `json:"reference_id,omitempty"`

	// ReferenceType labels what ReferenceID references (e.g., "ad_campaign").
	ReferenceType string `json:"reference_type,omitempty"`
}

// ─────────────────────────────────────────────
// Outbound DTOs  (VNPay Service → caller)
// ─────────────────────────────────────────────

// CreatePaymentResponse is returned after a payment URL has been generated.
type CreatePaymentResponse struct {
	// PaymentURL is the full URL to redirect the user to VNPay.
	PaymentURL string `json:"payment_url"`

	// OrderRef echoes the request's OrderRef for correlation.
	OrderRef string `json:"order_ref"`

	// AmountVND echoes the charged amount.
	AmountVND int64 `json:"amount_vnd"`

	// ExpireTime is when the payment URL will expire.
	ExpireTime time.Time `json:"expire_time"`
}

// ─────────────────────────────────────────────
// VNPay Return / IPN Parameter DTOs
// ─────────────────────────────────────────────

// PaymentReturnParams holds the query-string parameters VNPay appends to the
// Return URL after the user completes (or cancels) the payment.
// Bind them with `c.ShouldBindQuery(params)` in your Gin handler.
type PaymentReturnParams struct {
	// AmountRaw is vnp_Amount sent by VNPay (actual amount × 100).
	AmountRaw int64 `form:"vnp_Amount"            json:"vnp_Amount"`

	BankCode    string `form:"vnp_BankCode"          json:"vnp_BankCode"`
	BankTransNo string `form:"vnp_BankTranNo"        json:"vnp_BankTranNo"`
	CardType    string `form:"vnp_CardType"          json:"vnp_CardType"`
	OrderInfo   string `form:"vnp_OrderInfo"         json:"vnp_OrderInfo"`

	// PayDate is formatted as "yyyyMMddHHmmss" by VNPay.
	PayDate string `form:"vnp_PayDate"           json:"vnp_PayDate"`

	ResponseCode      string `form:"vnp_ResponseCode"      json:"vnp_ResponseCode"`
	TmnCode           string `form:"vnp_TmnCode"           json:"vnp_TmnCode"`
	TransactionNo     string `form:"vnp_TransactionNo"     json:"vnp_TransactionNo"`
	TransactionStatus string `form:"vnp_TransactionStatus" json:"vnp_TransactionStatus"`
	TxnRef            string `form:"vnp_TxnRef"            json:"vnp_TxnRef"`
	SecureHashType    string `form:"vnp_SecureHashType"    json:"vnp_SecureHashType"`
	SecureHash        string `form:"vnp_SecureHash"        json:"vnp_SecureHash"`
}

// AmountVND converts the raw wire-format amount (×100) back to VND.
func (p *PaymentReturnParams) AmountVND() int64 {
	return p.AmountRaw / 100
}

// IsSuccess returns true when both response fields indicate a successful payment.
func (p *PaymentReturnParams) IsSuccess() bool {
	return p.ResponseCode == ResponseCodeSuccess &&
		p.TransactionStatus == ResponseCodeSuccess
}

// PaymentIPNParams is a type alias because the IPN callback uses the same
// parameters as the return URL callback.
type PaymentIPNParams = PaymentReturnParams

// IPNResponse is the JSON body your server must return to VNPay's IPN caller.
// VNPay expects exactly {"RspCode": "XX", "Message": "..."}.
type IPNResponse struct {
	RspCode string `json:"RspCode"`
	Message string `json:"Message"`
}

// NewIPNResponse constructs an IPNResponse from a predefined IPN code.
// If the code is unknown, the "unknown error" response is used.
func NewIPNResponse(code string) IPNResponse {
	r, ok := IPNResponses[code]
	if !ok {
		r = IPNResponses[IPNCodeUnknownError]
	}
	return IPNResponse{RspCode: r["RspCode"], Message: r["Message"]}
}

// ─────────────────────────────────────────────
// Query Transaction DTOs
// ─────────────────────────────────────────────

// QueryTransactionRequest is the input for querying a transaction from VNPay.
type QueryTransactionRequest struct {
	// OrderRef is the original vnp_TxnRef of the transaction to query.
	OrderRef string `json:"order_ref" validate:"required"`

	// TransactionDate is the date the original transaction was created.
	// VNPay uses this to locate the transaction (format: yyyyMMddHHmmss).
	TransactionDate time.Time `json:"transaction_date" validate:"required"`

	// ClientIPAddr is the IP of the request initiator.
	ClientIPAddr string `json:"client_ip_addr" validate:"required"`
}

// QueryTransactionResponse is the parsed response from VNPay's query API.
type QueryTransactionResponse struct {
	ResponseCode      string `json:"vnp_ResponseCode"`
	Message           string `json:"vnp_Message"`
	TmnCode           string `json:"vnp_TmnCode"`
	TxnRef            string `json:"vnp_TxnRef"`
	AmountRaw         int64  `json:"vnp_Amount"`
	OrderInfo         string `json:"vnp_OrderInfo"`
	VNPayTransNo      string `json:"vnp_TransactionNo"`
	BankCode          string `json:"vnp_BankCode"`
	PayDate           string `json:"vnp_PayDate"`
	TransactionType   string `json:"vnp_TransactionType"`
	TransactionStatus string `json:"vnp_TransactionStatus"`
}

// AmountVND returns the transaction amount in VND (wire format ÷ 100).
func (r *QueryTransactionResponse) AmountVND() int64 {
	return r.AmountRaw / 100
}

// ─────────────────────────────────────────────
// Refund DTOs
// ─────────────────────────────────────────────

// RefundRequest is the input DTO for requesting a refund from VNPay.
type RefundRequest struct {
	// OrderRef is the original order reference (vnp_TxnRef).
	OrderRef string `json:"order_ref" validate:"required"`

	// TransactionNo is VNPay's own transaction number from the original payment.
	TransactionNo string `json:"transaction_no" validate:"required"`

	// AmountVND is the amount to refund in VND.
	// For full refunds use the original amount; for partial refunds any
	// lesser amount up to the original.
	AmountVND int64 `json:"amount_vnd" validate:"required,min=1000"`

	// TransactionDate is the date the original payment was settled.
	TransactionDate time.Time `json:"transaction_date" validate:"required"`

	// TransType is "02" (full) or "03" (partial).  Defaults to "02".
	TransType string `json:"trans_type"`

	// RefundOrderRef is a NEW unique reference for this refund request.
	RefundOrderRef string `json:"refund_order_ref" validate:"required,max=50"`

	// CreatedBy is the username / operator initiating the refund.
	CreatedBy string `json:"created_by" validate:"required"`

	// ClientIPAddr is the IP of the system/operator making the refund call.
	ClientIPAddr string `json:"client_ip_addr" validate:"required"`
}

// RefundResponse holds VNPay's response to a refund request.
type RefundResponse struct {
	ResponseCode    string `json:"vnp_ResponseCode"`
	Message         string `json:"vnp_Message"`
	TxnRef          string `json:"vnp_TxnRef"`
	AmountRaw       int64  `json:"vnp_Amount"`
	TransactionNo   string `json:"vnp_TransactionNo"`
	TransactionType string `json:"vnp_TransactionType"`
}

// AmountVND returns the refunded amount in VND.
func (r *RefundResponse) AmountVND() int64 {
	return r.AmountRaw / 100
}
