// pkg/encryption/encryption.go
package encryption

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// 错误定义
var (
	ErrHashFailed      = errors.New("密码加密失败")
	ErrInvalidHash     = errors.New("无效的密码哈希")
	ErrPasswordTooLong = errors.New("密码长度超过72字节限制")
)

// HashPassword 使用bcrypt加密密码
func HashPassword(password string) (string, error) {
	// bcrypt限制密码长度为72字节
	if len(password) > 72 {
		return "", ErrPasswordTooLong
	}

	// 生成盐值并加密，使用默认成本(10)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrHashFailed, err)
	}

	return string(hashedBytes), nil
}

// CheckPasswordHash 验证密码与哈希是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("生成随机字符串失败: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateSalt 生成盐值
func GenerateSalt() (string, error) {
	return GenerateRandomString(32) // 32字节 = 64字符十六进制
}

// IsHashed 检查字符串是否为bcrypt哈希格式
func IsHashed(password string) bool {
	// bcrypt哈希格式: $2a$, $2b$, $2y$ 开头，后面跟着成本和盐值
	return strings.HasPrefix(password, "$2a$") ||
		strings.HasPrefix(password, "$2b$") ||
		strings.HasPrefix(password, "$2y$")
}

// GetHashCost 获取哈希的计算成本
func GetHashCost(hashedPassword string) (int, error) {
	if !IsHashed(hashedPassword) {
		return 0, ErrInvalidHash
	}

	// bcrypt哈希格式: $2a$10$salt.hash
	parts := strings.Split(hashedPassword, "$")
	if len(parts) < 4 {
		return 0, ErrInvalidHash
	}

	var cost int
	_, err := fmt.Sscanf(parts[2], "%d", &cost)
	if err != nil {
		return 0, fmt.Errorf("解析哈希成本失败: %v", err)
	}

	return cost, nil
}

// NeedsRehash 检查哈希是否需要重新加密（成本过低）
func NeedsRehash(hashedPassword string) bool {
	cost, err := GetHashCost(hashedPassword)
	if err != nil {
		return true // 如果解析失败，建议重新加密
	}

	// 如果成本低于默认成本，需要重新加密
	return cost < bcrypt.DefaultCost
}

// RehashPassword 重新加密密码（如果成本过低）
func RehashPassword(password, currentHash string) (string, error) {
	if !NeedsRehash(currentHash) {
		return currentHash, nil // 不需要重新加密
	}

	return HashPassword(password)
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少8位")
	}

	if len(password) > 72 {
		return errors.New("密码长度不能超过72位")
	}

	// 检查包含数字
	hasDigit := false
	// 检查包含字母
	hasLetter := false
	// 检查包含特殊字符
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= '0' && char <= '9':
			hasDigit = true
		case (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z'):
			hasLetter = true
		case char >= 33 && char <= 126: // 可打印ASCII字符
			hasSpecial = true
		}
	}

	var requirements []string
	if !hasDigit {
		requirements = append(requirements, "数字")
	}
	if !hasLetter {
		requirements = append(requirements, "字母")
	}
	if !hasSpecial {
		requirements = append(requirements, "特殊字符")
	}

	if len(requirements) > 0 {
		return fmt.Errorf("密码必须包含: %s", strings.Join(requirements, ", "))
	}

	return nil
}
