package enum

// =============================================================================
// 1. THEME & DISPLAY
// =============================================================================

//go:generate enumer -type=ThemeMode -json -transform=snake -trimprefix=Theme
type ThemeMode int

const (
	ThemeSystem ThemeMode = iota
	ThemeLight
	ThemeDark
	ThemeHighContrast
)

//go:generate enumer -type=FontSize -json -transform=snake -trimprefix=Font
type FontSize int

const (
	FontMedium FontSize = iota
	FontLarge
)


// =============================================================================
// 3. LANGUAGE & REGION
// =============================================================================

//go:generate enumer -type=LangCode -json -transform=snake -trimprefix=Lang
type LangCode int

const (
	LangVi LangCode = iota
	LangEn
	LangJp
)

//go:generate enumer -type=Timezone -json -trimprefix=Timezone
type Timezone int

const (
	TimezoneHCM Timezone = iota // Sẽ lưu là "hcm" -> Cần Custom String() nếu muốn lưu "Asia/Ho_Chi_Minh"
	TimezoneNY
	TimezoneLDN
)
// Lưu ý: Với Timezone, nếu bạn muốn lưu chính xác chuỗi "Asia/Ho_Chi_Minh"
// thay vì "hcm" (do enumer tự transform), bạn nên giữ nguyên type Timezone string
// hoặc phải viết hàm String() thủ công.
// Ở ví dụ này, để an toàn và chuẩn xác cho Timezone, tôi khuyên nên giữ nguyên String cho Timezone.

// =============================================================================
// 4. NOTIFICATIONS
// =============================================================================

