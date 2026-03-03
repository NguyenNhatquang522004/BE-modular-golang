package vnpay

// ─────────────────────────────────────────────
// VNPay Configuration
// ─────────────────────────────────────────────

// Config holds all parameters needed to initialise the VNPay service.
// It is populated from the application-wide configs.Config struct.
type Config struct {
	// TmnCode is the Terminal Code assigned by VNPay (vnp_TmnCode).
	TmnCode string

	// HashSecret is the signing key assigned by VNPay for HMAC-SHA512.
	HashSecret string

	// PaymentURL is the VNPay checkout page URL.
	// Leave empty to auto-select based on IsSandbox.
	PaymentURL string

	// APIURL is the VNPay Merchant API URL (query / refund).
	// Leave empty to auto-select based on IsSandbox.
	APIURL string

	// ReturnURL is the URL VNPay redirects the customer back to after payment.
	ReturnURL string

	// IPNURL is the server-to-server Instant Payment Notification URL.
	// If empty, ReturnURL is used as fallback.
	IPNURL string

	// Version is the VNPay API version. Defaults to APIVersion ("2.1.0").
	Version string

	// IsSandbox controls which VNPay endpoints are used.
	// true  → sandbox  (for development / QA)
	// false → production
	IsSandbox bool
}

// Validate ensures that all required Config fields are present and sets
// sensible defaults for optional fields.  Call this before constructing
// a Service to catch mis-configurations early.
func (c *Config) Validate() error {
	if c.TmnCode == "" {
		return ErrInvalidConfig
	}
	if c.HashSecret == "" {
		return ErrInvalidConfig
	}
	if c.ReturnURL == "" {
		return ErrInvalidConfig
	}

	// Auto-select endpoint URLs based on sandbox flag.
	if c.PaymentURL == "" {
		if c.IsSandbox {
			c.PaymentURL = SandboxPayURL
		} else {
			c.PaymentURL = ProductionPayURL
		}
	}
	if c.APIURL == "" {
		if c.IsSandbox {
			c.APIURL = SandboxAPIURL
		} else {
			c.APIURL = ProductionAPIURL
		}
	}

	// IPN URL falls back to return URL if not explicitly provided.
	if c.IPNURL == "" {
		c.IPNURL = c.ReturnURL
	}

	// Default API version.
	if c.Version == "" {
		c.Version = APIVersion
	}

	return nil
}
