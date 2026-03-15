package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateSecureStreamKey(ownerID string, sessionID string, secretKey string, ttl time.Duration) string {
	// 1. Tính toán thời gian hết hạn (Unix timestamp)
	expireAt := time.Now().Add(ttl).Unix()

	// 2. Tạo Payload cần mã hóa
	// Gắn thêm sessionID để đảm bảo mỗi phiên live có một key duy nhất,
	// tránh việc user dùng lại key cũ cho session mới.
	payload := fmt.Sprintf("%s:%s:%d", ownerID, sessionID, expireAt)

	// 3. Tạo chữ ký HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	// 4. Lắp ráp Stream Key hoàn chỉnh
	// Format: live_{ownerID}_{expireAt}_{signature[0:16]} (Cắt ngắn signature cho đỡ dài)
	shortSig := signature[:16]
	return fmt.Sprintf("live_%s_%d_%s", ownerID, expireAt, shortSig)
}

// VerifyStreamKey (Tùy chọn) dùng ở Backend hoặc Auth Webhook của Media Server
func VerifyStreamKey(streamKey string, ownerID string, sessionID string, secretKey string) bool {
	// 1. Phân tách Stream Key
	// Định dạng chuẩn lúc tạo: live_{ownerID}_{expireAt}_{shortSig}
	parts := strings.Split(streamKey, "_")

	// Bắt buộc phải có đúng 4 phần: "live", ownerID, expireAt, shortSig
	if len(parts) != 4 {
		return false
	}

	// Xác thực prefix
	if parts[0] != "live" {
		return false
	}

	extractedOwnerID := parts[1]
	expireAtStr := parts[2]
	providedShortSig := parts[3]

	// Đảm bảo Stream Key được dùng đúng cho OwnerID đang request
	if extractedOwnerID != ownerID {
		return false
	}

	// 2. Kiểm tra thời hạn (Expiration)
	expireAt, err := strconv.ParseInt(expireAtStr, 10, 64)
	if err != nil {
		return false // Lỗi parse số (Key bị can thiệp sai định dạng)
	}

	// Nếu thời gian hiện tại lớn hơn thời gian hết hạn -> Reject
	if time.Now().Unix() > expireAt {
		return false
	}

	// 3. Dựng lại Payload gốc để kiểm tra chữ ký
	// Phải dùng chung công thức: ownerID + sessionID + expireAt như lúc tạo
	payload := fmt.Sprintf("%s:%s:%d", ownerID, sessionID, expireAt)

	// Băm lại bằng HMAC-SHA256 với Secret Key
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	expectedShortSig := expectedSignature[:16]

	// 4. So sánh an toàn tuyệt đối (Constant Time Comparison)
	// BEST PRACTICE: Không dùng toán tử "==" với dữ liệu mã hóa
	return hmac.Equal([]byte(providedShortSig), []byte(expectedShortSig))
}
