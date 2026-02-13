// Package crypto 提供加密解密功能
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/gogf/gf/v2/errors/gerror"
)

// EncryptAES 使用 AES-256-GCM 加密明文
//
// 参数:
//   - plaintext: 待加密的明文
//   - key: 32字节的加密密钥（AES-256）
//
// 返回:
//   - string: Base64编码的密文
//   - error: 错误信息
func EncryptAES(plaintext, key string) (string, error) {
	// 验证密钥长度
	if len(key) != 32 {
		return "", gerror.New("密钥长度必须为32字节(AES-256)")
	}

	// 创建 AES cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", gerror.Wrap(err, "创建AES cipher失败")
	}

	// 创建 GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", gerror.Wrap(err, "创建GCM模式失败")
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", gerror.Wrap(err, "生成nonce失败")
	}

	// 加密（nonce + ciphertext）
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64 编码
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

// DecryptAES 使用 AES-256-GCM 解密密文
//
// 参数:
//   - ciphertext: Base64编码的密文
//   - key: 32字节的解密密钥（AES-256）
//
// 返回:
//   - string: 解密后的明文
//   - error: 错误信息
func DecryptAES(ciphertext, key string) (string, error) {
	// 验证密钥长度
	if len(key) != 32 {
		return "", gerror.New("密钥长度必须为32字节(AES-256)")
	}

	// Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", gerror.Wrap(err, "Base64解码失败")
	}

	// 创建 AES cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", gerror.Wrap(err, "创建AES cipher失败")
	}

	// 创建 GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", gerror.Wrap(err, "创建GCM模式失败")
	}

	// 验证密文长度
	nonceSize := gcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", gerror.New("密文长度不足")
	}

	// 分离 nonce 和密文
	nonce, ciphertextBytes := decoded[:nonceSize], decoded[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", gerror.Wrap(err, "解密失败")
	}

	return string(plaintext), nil
}
