package utils

import "golang.org/x/crypto/bcrypt"

func Hash(password string) (string, error) {
	// bcrypt.DefaultCost hiện tại là 10.
	// Bạn có thể tăng lên 12 hoặc 14 nếu server mạnh, nhưng 10 là chuẩn cân bằng.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check: Kiểm tra mật khẩu khi đăng nhập
// Trả về true nếu đúng, false nếu sai
func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
