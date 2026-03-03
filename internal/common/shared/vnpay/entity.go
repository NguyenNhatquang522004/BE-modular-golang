package vnpay

import "time"

// ─────────────────────────────────────────────
// TransactionStatus Enum
// ─────────────────────────────────────────────

// TransactionStatus represents the lifecycle state of a payment transaction
// as tracked in the application's own database (NOT VNPay's internal status).
type TransactionStatus string

const (
	// StatusPending means the payment URL was created but the user has not
	// completed (or cancelled) the payment yet.
	StatusPending TransactionStatus = "PENDING"

	// StatusSuccess means VNPay confirmed the payment was successful.
	StatusSuccess TransactionStatus = "SUCCESS"

	// StatusFailed means VNPay returned a non-success response code.
	StatusFailed TransactionStatus = "FAILED"

	// StatusCancelled means the user explicitly cancelled the transaction.
	StatusCancelled TransactionStatus = "CANCELLED"

	// StatusRefunded means a full or partial refund has been processed.
	StatusRefunded TransactionStatus = "REFUNDED"

	// StatusExpired means the payment URL was not used within its validity window.
	StatusExpired TransactionStatus = "EXPIRED"
)

// ─────────────────────────────────────────────
// PaymentTransaction – Domain Entity
// ─────────────────────────────────────────────

// PaymentTransaction is the central domain entity representing one VNPay
// payment attempt.  It is persisted in the Postgres table
// "vnpay_payment_transactions".
//
// Mapping to VNPay query parameters:
//
//	OrderRef        → vnp_TxnRef         (your side's unique order ID)
//	AmountVND       → vnp_Amount / 100   (VNPay sends amount * 100)
//	VNPayTransNo    → vnp_TransactionNo  (VNPay's own transaction ID)
//	BankTransNo     → vnp_BankTranNo     (Bank's transaction reference)
//	ResponseCode    → vnp_ResponseCode
//	TransactionStatus → vnp_TransactionStatus
type PaymentTransaction struct {
	// ── Primary Key ──────────────────────────────────────────────────────────
	ID string `json:"id"             gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// ── Business Keys ────────────────────────────────────────────────────────
	// OrderRef is YOUR unique order reference (vnp_TxnRef).
	// It must be unique per transaction request to VNPay.
	OrderRef string `json:"order_ref"      gorm:"uniqueIndex;not null;size:50"`

	// ── Amount ───────────────────────────────────────────────────────────────
	// AmountVND is the transaction amount in Vietnamese Dong (VND).
	// VNPay wire format is AmountVND * 100; conversion is handled by the service.
	AmountVND int64 `json:"amount_vnd"     gorm:"not null"`

	// CurrencyCode defaults to "VND".
	CurrencyCode string `json:"currency_code"  gorm:"size:10;default:'VND'"`

	// ── Order Info ───────────────────────────────────────────────────────────
	OrderInfo string `json:"order_info"     gorm:"size:255"`
	OrderType string `json:"order_type"     gorm:"size:50"`

	// BankCode is the bank selected at checkout (may be empty if VNPay chose).
	BankCode string `json:"bank_code"      gorm:"size:20"`

	// ── VNPay Response Fields ────────────────────────────────────────────────
	VNPayTransNo      string `json:"vnpay_trans_no"    gorm:"size:20"`
	BankTransNo       string `json:"bank_trans_no"     gorm:"size:50"`
	CardType          string `json:"card_type"         gorm:"size:20"`
	ResponseCode      string `json:"response_code"     gorm:"size:5"`
	TransactionStatus string `json:"transaction_status" gorm:"size:5"`

	// ── Payment Dates ────────────────────────────────────────────────────────
	PayDate    *time.Time `json:"pay_date,omitempty"`
	ExpireDate *time.Time `json:"expire_date,omitempty"`

	// ── Network ──────────────────────────────────────────────────────────────
	IPAddr string `json:"ip_addr"        gorm:"size:45"`
	Locale string `json:"locale"         gorm:"size:5;default:'vn'"`

	// ── Application Status ───────────────────────────────────────────────────
	// Status is the application-level lifecycle state (see TransactionStatus).
	Status TransactionStatus `json:"status" gorm:"size:20;default:'PENDING'"`

	// SecureHash stores the computed HMAC-SHA512 so it can be cross-checked later.
	SecureHash string `json:"secure_hash"    gorm:"type:text"`

	// ── Context / Foreign Keys (optional) ────────────────────────────────────
	// UserID links the transaction to an authenticated user.
	UserID string `json:"user_id,omitempty"         gorm:"size:36;index"`

	// ReferenceID links the transaction to the business entity being paid for
	// (e.g., an ad campaign ID, a subscription ID).
	ReferenceID string `json:"reference_id,omitempty"   gorm:"size:36;index"`

	// ReferenceType labels what ReferenceID points to (e.g., "ad_campaign").
	ReferenceType string `json:"reference_type,omitempty" gorm:"size:50"`

	// ── Timestamps ───────────────────────────────────────────────────────────
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the GORM default table name.
func (PaymentTransaction) TableName() string {
	return "vnpay_payment_transactions"
}

// ─────────────────────────────────────────────
// Domain Helpers
// ─────────────────────────────────────────────

// IsSuccess returns true when both sides (VNPay + our app) confirm success.
func (t *PaymentTransaction) IsSuccess() bool {
	return t.Status == StatusSuccess &&
		t.ResponseCode == ResponseCodeSuccess &&
		t.TransactionStatus == ResponseCodeSuccess
}

// IsPending returns true when the transaction is awaiting the user's action.
func (t *PaymentTransaction) IsPending() bool {
	return t.Status == StatusPending
}

// IsTerminal returns true when no further state transitions are possible.
func (t *PaymentTransaction) IsTerminal() bool {
	switch t.Status {
	case StatusSuccess, StatusFailed, StatusCancelled, StatusRefunded, StatusExpired:
		return true
	}
	return false
}
