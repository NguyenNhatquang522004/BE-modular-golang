package vnpay

// ─────────────────────────────────────────────
// API Metadata
// ─────────────────────────────────────────────

// APIVersion is the VNPay API version used by this integration.
const APIVersion = "2.1.0"

// SecureHashType is the hash algorithm used for signature generation.
const SecureHashType = "SHA512"

// ─────────────────────────────────────────────
// Command Codes (vnp_Command)
// ─────────────────────────────────────────────

const (
	CommandPay    = "pay"
	CommandQuery  = "querydr"
	CommandRefund = "refund"
)

// ─────────────────────────────────────────────
// Locale Codes (vnp_Locale)
// ─────────────────────────────────────────────

const (
	LocaleVN = "vn"
	LocaleEN = "en"
)

// ─────────────────────────────────────────────
// Currency Code (vnp_CurrCode)
// ─────────────────────────────────────────────

const CurrencyVND = "VND"

// ─────────────────────────────────────────────
// Order Type Codes (vnp_OrderType)
// Reference: https://sandbox.vnpayment.vn/apis/docs/loai-hang-hoa/
// ─────────────────────────────────────────────

const (
	OrderTypeBillPayment = "billpayment" // Thanh toán hóa đơn
	OrderTypeTopUp       = "topup"       // Nạp tiền điện thoại / tài khoản
	OrderTypeFashion     = "fashion"     // Thời trang
	OrderTypeFBLiquor    = "fbliquor"    // Đồ uống có cồn
	OrderTypeFBFood      = "fbbfood"     // Thức ăn nhanh
	OrderTypeFBDrink     = "fbdrink"     // Đồ uống không cồn
	OrderTypeOther       = "other"       // Các loại hàng hóa / dịch vụ khác
)

// ─────────────────────────────────────────────
// Transaction Type Codes (used in refund)
// ─────────────────────────────────────────────

const (
	TransTypeFullRefund    = "02" // Hoàn trả toàn phần
	TransTypePartialRefund = "03" // Hoàn trả một phần
)

// ─────────────────────────────────────────────
// Payment Gateway URLs
// ─────────────────────────────────────────────

const (
	// SandboxPayURL is the VNPay sandbox payment endpoint.
	SandboxPayURL = "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"

	// ProductionPayURL is the VNPay production payment endpoint.
	ProductionPayURL = "https://pay.vnpay.vn/vpcpay.html"

	// SandboxAPIURL is the VNPay sandbox Merchant API (for query/refund).
	SandboxAPIURL = "https://sandbox.vnpayment.vn/merchant_webapi/api/transaction"

	// ProductionAPIURL is the VNPay production Merchant API (for query/refund).
	ProductionAPIURL = "https://pay.vnpay.vn/merchant_webapi/api/transaction"
)

// ─────────────────────────────────────────────
// Response Codes (vnp_ResponseCode / vnp_TransactionStatus)
// Reference: https://sandbox.vnpayment.vn/apis/docs/ma-tra-ve/
// "00" means SUCCESS for both fields.
// ─────────────────────────────────────────────

const (
	ResponseCodeSuccess             = "00" // Giao dịch thành công
	ResponseCodeSuspiciousFraud     = "07" // Giao dịch bị nghi ngờ gian lận
	ResponseCodeCardNotRegistered   = "09" // Thẻ/TK chưa đăng ký InternetBanking
	ResponseCodeAuthFailed3Times    = "10" // Xác thực thẻ/TK sai quá 3 lần
	ResponseCodePaymentTimeout      = "11" // Hết hạn chờ thanh toán
	ResponseCodeCardLocked          = "12" // Thẻ/TK bị khóa
	ResponseCodeWrongOTP            = "13" // Sai mật khẩu OTP
	ResponseCodeCustomerCancelled   = "24" // Khách hàng hủy giao dịch
	ResponseCodeInsufficientBalance = "51" // Số dư không đủ
	ResponseCodeExceedDailyLimit    = "65" // Vượt hạn mức giao dịch trong ngày
	ResponseCodeBankMaintenance     = "75" // Ngân hàng đang bảo trì
	ResponseCodeWrongOTPExceeded    = "79" // Sai mật khẩu thanh toán quá số lần cho phép
	ResponseCodeSystemError         = "99" // Lỗi hệ thống không xác định
)

// ResponseCodeMessages maps VNPay response codes to Vietnamese descriptions.
var ResponseCodeMessages = map[string]string{
	ResponseCodeSuccess:             "Giao dịch thành công",
	ResponseCodeSuspiciousFraud:     "Giao dịch bị nghi ngờ gian lận",
	ResponseCodeCardNotRegistered:   "Thẻ/Tài khoản chưa đăng ký dịch vụ InternetBanking",
	ResponseCodeAuthFailed3Times:    "Xác thực thông tin thẻ/Tài khoản không đúng quá 3 lần",
	ResponseCodePaymentTimeout:      "Hết hạn chờ thanh toán. Vui lòng thực hiện lại giao dịch.",
	ResponseCodeCardLocked:          "Thẻ/Tài khoản của khách hàng bị khóa",
	ResponseCodeWrongOTP:            "Nhập sai mật khẩu xác thực giao dịch (OTP)",
	ResponseCodeCustomerCancelled:   "Khách hàng hủy giao dịch",
	ResponseCodeInsufficientBalance: "Tài khoản không đủ số dư để thực hiện giao dịch",
	ResponseCodeExceedDailyLimit:    "Tài khoản đã vượt quá hạn mức giao dịch trong ngày",
	ResponseCodeBankMaintenance:     "Ngân hàng thanh toán đang bảo trì",
	ResponseCodeWrongOTPExceeded:    "Nhập sai mật khẩu thanh toán quá số lần cho phép",
	ResponseCodeSystemError:         "Lỗi không xác định",
}

// ─────────────────────────────────────────────
// IPN Response Codes (sent from your server back to VNPay)
// Reference: https://sandbox.vnpayment.vn/apis/docs/thanh-toan-pay/ipn-url.html
// ─────────────────────────────────────────────

const (
	IPNCodeConfirmSuccess   = "00" // Xác nhận thành công
	IPNCodeOrderNotFound    = "01" // Không tìm thấy đơn hàng
	IPNCodeAlreadyUpdated   = "02" // Đơn hàng đã được cập nhật trước đó
	IPNCodeInvalidAmount    = "04" // Số tiền không hợp lệ
	IPNCodeInvalidSignature = "97" // Checksum không hợp lệ
	IPNCodeUnknownError     = "99" // Lỗi không xác định
)

const (
	IPNMessageConfirmSuccess   = "Confirm Success"
	IPNMessageOrderNotFound    = "Order Not Found"
	IPNMessageAlreadyUpdated   = "Order Already Confirmed"
	IPNMessageInvalidAmount    = "Invalid Amount"
	IPNMessageInvalidSignature = "Invalid Checksum"
	IPNMessageUnknownError     = "Unknown Error"
)

// IPNResponses is a lookup for standard IPN response bodies.
var IPNResponses = map[string]map[string]string{
	IPNCodeConfirmSuccess:   {"RspCode": IPNCodeConfirmSuccess, "Message": IPNMessageConfirmSuccess},
	IPNCodeOrderNotFound:    {"RspCode": IPNCodeOrderNotFound, "Message": IPNMessageOrderNotFound},
	IPNCodeAlreadyUpdated:   {"RspCode": IPNCodeAlreadyUpdated, "Message": IPNMessageAlreadyUpdated},
	IPNCodeInvalidAmount:    {"RspCode": IPNCodeInvalidAmount, "Message": IPNMessageInvalidAmount},
	IPNCodeInvalidSignature: {"RspCode": IPNCodeInvalidSignature, "Message": IPNMessageInvalidSignature},
	IPNCodeUnknownError:     {"RspCode": IPNCodeUnknownError, "Message": IPNMessageUnknownError},
}

// ─────────────────────────────────────────────
// Card Types (vnp_CardType)
// ─────────────────────────────────────────────

const (
	CardTypeATM    = "ATM"    // Thẻ ATM / Tài khoản ngân hàng nội địa
	CardTypeCredit = "CREDIT" // Thẻ thanh toán quốc tế (Visa, Mastercard…)
	CardTypeQRCode = "QRCODE" // VNPay QR
)

// ─────────────────────────────────────────────
// Minimum Amount
// ─────────────────────────────────────────────

// MinimumPaymentAmount is the minimum transaction amount VNPay accepts (in VND).
const MinimumPaymentAmount int64 = 5_000

// DefaultPaymentExpiryMinutes is the default number of minutes before a payment URL expires.
const DefaultPaymentExpiryMinutes = 15
