package crypto

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestEncryptDecryptAES(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试正常加密解密
		key := "12345678901234567890123456789012" // 32字节
		plaintext := "this is a secret API key"

		// 加密
		encrypted, err := EncryptAES(plaintext, key)
		t.AssertNil(err)
		t.AssertNE(encrypted, "")
		t.AssertNE(encrypted, plaintext)

		// 解密
		decrypted, err := DecryptAES(encrypted, key)
		t.AssertNil(err)
		t.AssertEQ(decrypted, plaintext)
	})
}

func TestEncryptAES_InvalidKeyLength(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试错误的密钥长度
		key := "short_key"
		plaintext := "test"

		_, err := EncryptAES(plaintext, key)
		t.AssertNE(err, nil)
	})
}

func TestDecryptAES_InvalidCiphertext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试无效的密文
		key := "12345678901234567890123456789012"
		invalidCiphertext := "invalid_base64_string"

		_, err := DecryptAES(invalidCiphertext, key)
		t.AssertNE(err, nil)
	})
}

func TestDecryptAES_WrongKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试使用错误密钥解密
		key1 := "12345678901234567890123456789012"
		key2 := "abcdefghijklmnopqrstuvwxyz123456"
		plaintext := "secret"

		// 用 key1 加密
		encrypted, err := EncryptAES(plaintext, key1)
		t.AssertNil(err)

		// 用 key2 解密（应该失败）
		_, err = DecryptAES(encrypted, key2)
		t.AssertNE(err, nil)
	})
}
