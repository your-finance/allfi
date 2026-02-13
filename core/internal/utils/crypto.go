// Package utils 加密工具
// 提供 AES-256-GCM 加密/解密功能，用于保护 API Key 等敏感数据
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Crypto 加密工具结构体
type Crypto struct {
	masterKey []byte
}

// NewCrypto 创建新的加密工具实例
// masterKey: 32 字节的主密钥（AES-256）
func NewCrypto(masterKey string) (*Crypto, error) {
	// 解码 Base64 格式的主密钥
	keyBytes, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		// 如果不是 Base64 格式，直接使用字符串
		keyBytes = []byte(masterKey)
	}

	// 验证密钥长度（AES-256 需要 32 字节）
	if len(keyBytes) != 32 {
		return nil, errors.New("主密钥长度必须为 32 字节")
	}

	return &Crypto{masterKey: keyBytes}, nil
}

// Encrypt 使用 AES-256-GCM 加密数据
// plaintext: 明文字符串
// 返回: Base64 编码的密文（包含 nonce）
func (c *Crypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// 创建 AES cipher
	block, err := aes.NewCipher(c.masterKey)
	if err != nil {
		return "", err
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密数据（nonce 附加在密文前面）
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// 返回 Base64 编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用 AES-256-GCM 解密数据
// ciphertext: Base64 编码的密文
// 返回: 明文字符串
func (c *Crypto) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// 解码 Base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建 AES cipher
	block, err := aes.NewCipher(c.masterKey)
	if err != nil {
		return "", err
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 验证数据长度
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文数据长度不足")
	}

	// 分离 nonce 和密文
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateMasterKey 生成新的 32 字节主密钥
// 返回: Base64 编码的主密钥
func GenerateMasterKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// ValidateMasterKey 验证主密钥格式是否正确
func ValidateMasterKey(masterKey string) bool {
	if masterKey == "" {
		return false
	}

	// 尝试解码 Base64
	keyBytes, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		// 如果不是 Base64，直接检查长度
		return len(masterKey) == 32
	}

	return len(keyBytes) == 32
}

// 全局加密工具实例（延迟初始化）
var globalCrypto *Crypto

// InitGlobalCrypto 初始化全局加密工具
func InitGlobalCrypto(masterKey string) error {
	crypto, err := NewCrypto(masterKey)
	if err != nil {
		return err
	}
	globalCrypto = crypto
	return nil
}

// EncryptAPIKey 加密 API Key（使用全局加密工具）
func EncryptAPIKey(apiKey string) (string, error) {
	if globalCrypto == nil {
		return "", errors.New("加密工具未初始化")
	}
	return globalCrypto.Encrypt(apiKey)
}

// DecryptAPIKey 解密 API Key（使用全局加密工具）
func DecryptAPIKey(encryptedKey string) (string, error) {
	if globalCrypto == nil {
		return "", errors.New("加密工具未初始化")
	}
	return globalCrypto.Decrypt(encryptedKey)
}
