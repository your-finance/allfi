package logic_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"your-finance/allfi/internal/app/auth/logic"
	"your-finance/allfi/internal/database"
)

func init() {
	ctx := context.Background()
	if err := database.Initialize(ctx); err != nil {
		g.Log().Fatalf(ctx, "Failed to initialize test DB: %v", err)
	}
}

// validateComplexPassword 验证复杂密码格式
// 要求：8-20 位，必须含大小写字母和数字，可包含特殊字符
func validateComplexPassword(password string) bool {
	length := len(password)
	if length < 8 || length > 20 {
		return false
	}

	var hasLower, hasUpper, hasDigit bool
	for _, ch := range password {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		case isAllowedSpecialChar(ch):
			// 允许的特殊字符
		default:
			// 非法字符
			return false
		}
	}

	return hasLower && hasUpper && hasDigit
}

// isAllowedSpecialChar 检查是否为允许的特殊字符
func isAllowedSpecialChar(ch rune) bool {
	switch ch {
	case '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '-', '=', '[', ']', '{', '}', '|', ';', '\'', ':', '"', ',', '.', '<', '>', '?':
		return true
	default:
		return false
	}
}

func TestAuth_SetupAndLogin(t *testing.T) {
	ctx := context.Background()

	// Create a new instance
	authService := logic.New()

	// Clear DB for a fresh test by dropping system_config
	_, _ = g.DB().Exec(ctx, "DELETE FROM system_config;")

	// 1. Initial status should be un-setup
	status, err := authService.GetStatus(ctx)
	require.NoError(t, err)
	assert.False(t, status.PinSet)

	// 2. Setup with invalid PIN should fail
	_, err = authService.Setup(ctx, "123") // too short
	assert.Error(t, err)

	// 3. Setup with valid PIN should succeed
	setupRes, err := authService.Setup(ctx, "123456")
	require.NoError(t, err)
	assert.NotEmpty(t, setupRes.Token)

	// 4. Status should now be setup
	status, err = authService.GetStatus(ctx)
	require.NoError(t, err)
	assert.True(t, status.PinSet)

	// 5. Trying to setup again should fail
	_, err = authService.Setup(ctx, "654321")
	assert.Error(t, err)

	// 6. Login with wrong PIN should fail
	_, err = authService.Login(ctx, "111111")
	assert.Error(t, err)

	// 7. Login with correct PIN should succeed
	loginRes, err := authService.Login(ctx, "123456")
	require.NoError(t, err)
	assert.NotEmpty(t, loginRes.Token)
}

func TestAuth_ChangePin(t *testing.T) {
	ctx := context.Background()
	authService := logic.New()

	// Clear DB for a fresh test by dropping system_config
	_, _ = g.DB().Exec(ctx, "DELETE FROM system_config;")

	// Setup initial PIN
	_, err := authService.Setup(ctx, "123456")
	require.NoError(t, err)

	// Change PIN with wrong current PIN
	_, err = authService.ChangePin(ctx, "654321", "111111")
	assert.Error(t, err)

	// Change PIN with invalid new PIN
	_, err = authService.ChangePin(ctx, "123456", "111")
	assert.Error(t, err)

	// Successful Change PIN
	res, err := authService.ChangePin(ctx, "123456", "654321")
	require.NoError(t, err)
	assert.True(t, res.Success)

	// Login with old PIN should fail
	_, err = authService.Login(ctx, "123456")
	assert.Error(t, err)

	// Login with new PIN should succeed
	loginRes, err := authService.Login(ctx, "654321")
	require.NoError(t, err)
	assert.NotEmpty(t, loginRes.Token)
}

func TestAuth_2FA(t *testing.T) {
	ctx := context.Background()
	authService := logic.New()

	_, _ = g.DB().Exec(ctx, "DELETE FROM system_config;")

	// Initial Setup
	_, err := authService.Setup(ctx, "123456")
	require.NoError(t, err)

	// 1. Setup 2FA
	setup2FARes, err := authService.Setup2FA(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, setup2FARes.Secret)
	assert.NotEmpty(t, setup2FARes.QrUrl)

	// 2. Enable 2FA with wrong code
	_, err = authService.Enable2FA(ctx, "111111")
	assert.Error(t, err)

	// 3. Enable 2FA with correct code
	validCode, err := totp.GenerateCode(setup2FARes.Secret, time.Now())
	require.NoError(t, err)

	enableRes, err := authService.Enable2FA(ctx, validCode)
	require.NoError(t, err)
	assert.True(t, enableRes.Success)

	// 4. Login after 2FA enabled
	loginRes, err := authService.Login(ctx, "123456")
	require.NoError(t, err)
	assert.NotEmpty(t, loginRes.Token)
	// The token generated should be a 2FA pending token, but we can't easily assert claims here without decoding

	// 5. Verify 2FA
	verifyCode, err := totp.GenerateCode(setup2FARes.Secret, time.Now())
	require.NoError(t, err)

	verifyRes, err := authService.Verify2FA(ctx, verifyCode)
	require.NoError(t, err)
	assert.True(t, verifyRes.Success)
	assert.NotEmpty(t, verifyRes.Token)

	// 6. Disable 2FA
	disableCode, err := totp.GenerateCode(setup2FARes.Secret, time.Now())
	require.NoError(t, err)

	disableRes, err := authService.Disable2FA(ctx, disableCode)
	require.NoError(t, err)
	assert.True(t, disableRes.Success)
}

func TestComplexPasswordValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"有效复杂密码", "Abc12345", true},
		{"带符号的复杂密码", "Abc@12345", true},
		{"缺少大写字母", "abc12345", false},
		{"缺少小写字母", "ABC12345", false},
		{"缺少数字", "Abcdefgh", false},
		{"太短", "Abc123", false},
		{"太长", "Abc12345678901234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateComplexPassword(tt.password)
			assert.Equal(t, tt.valid, valid)
		})
	}
}

func TestAuth_ComplexPassword(t *testing.T) {
	ctx := context.Background()
	authService := logic.New()

	// Clear DB for a fresh test
	_, _ = g.DB().Exec(ctx, "DELETE FROM system_config;")

	// 1. Setup with complex password should succeed
	setupRes, err := authService.Setup(ctx, "Abc12345")
	require.NoError(t, err)
	assert.NotEmpty(t, setupRes.Token)

	// 2. Status should show complex password type
	status, err := authService.GetStatus(ctx)
	require.NoError(t, err)
	assert.True(t, status.PinSet)
	assert.Equal(t, "complex", status.PasswordType)

	// 3. Login with complex password should succeed
	loginRes, err := authService.Login(ctx, "Abc12345")
	require.NoError(t, err)
	assert.NotEmpty(t, loginRes.Token)

	// 4. Login with wrong password should fail
	_, err = authService.Login(ctx, "abc12345")
	assert.Error(t, err)
}

func TestAuth_SwitchType(t *testing.T) {
	ctx := context.Background()
	authService := logic.New()

	// Clear DB for a fresh test
	_, _ = g.DB().Exec(ctx, "DELETE FROM system_config;")

	// 1. Setup with PIN
	_, err := authService.Setup(ctx, "123456")
	require.NoError(t, err)

	// 2. Verify initial status is PIN type
	status, err := authService.GetStatus(ctx)
	require.NoError(t, err)
	assert.Equal(t, "pin", status.PasswordType)

	// 3. Switch to complex password with wrong current password should fail
	_, err = authService.SwitchType(ctx, "wrong", "complex", "NewPass123")
	assert.Error(t, err)

	// 4. Switch to complex password with invalid new password should fail
	_, err = authService.SwitchType(ctx, "123456", "complex", "weak")
	assert.Error(t, err)

	// 5. Switch to complex password with valid credentials should succeed
	switchRes, err := authService.SwitchType(ctx, "123456", "complex", "NewPass123")
	require.NoError(t, err)
	assert.True(t, switchRes.Success)

	// 6. Verify new status is complex type
	status, err = authService.GetStatus(ctx)
	require.NoError(t, err)
	assert.Equal(t, "complex", status.PasswordType)

	// 7. Login with old PIN should fail
	_, err = authService.Login(ctx, "123456")
	assert.Error(t, err)

	// 8. Login with new complex password should succeed
	loginRes, err := authService.Login(ctx, "NewPass123")
	require.NoError(t, err)
	assert.NotEmpty(t, loginRes.Token)

	// 9. Switch back to PIN should succeed
	switchRes, err = authService.SwitchType(ctx, "NewPass123", "pin", "654321")
	require.NoError(t, err)
	assert.True(t, switchRes.Success)

	// 10. Verify status is PIN type again
	status, err = authService.GetStatus(ctx)
	require.NoError(t, err)
	assert.Equal(t, "pin", status.PasswordType)

	// 11. Login with new PIN should succeed
	loginRes, err = authService.Login(ctx, "654321")
	require.NoError(t, err)
	assert.NotEmpty(t, loginRes.Token)
}
