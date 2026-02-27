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
