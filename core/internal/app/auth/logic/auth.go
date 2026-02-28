// Package logic 认证业务逻辑
// PIN 码认证：bcrypt 哈希存储 + JWT Token + 锁定保护
package logic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	authApi "your-finance/allfi/api/v1/auth"
	"your-finance/allfi/internal/app/auth/model"
	"your-finance/allfi/internal/app/auth/service"
	systemDao "your-finance/allfi/internal/app/system/dao"
	systemEntity "your-finance/allfi/internal/app/system/model/entity"

	"github.com/pquerna/otp/totp"
)

// pinPattern PIN 格式正则：4-20 位数字
var pinPattern = regexp.MustCompile(`^\d{4,20}$`)

// sAuth 认证服务实现
type sAuth struct {
	jwtSecret []byte // JWT 签名密钥
}

// New 创建认证服务实例
func New() service.IAuth {
	s := &sAuth{}
	s.initJWTSecret()
	return s
}

// initJWTSecret 初始化或加载 JWT 签名密钥
func (s *sAuth) initJWTSecret() {
	ctx := context.Background()

	// 从 system_config 表读取 JWT 密钥
	secret := s.getConfigValue(ctx, model.ConfigKeyJWTSecret)
	if secret == "" {
		// 首次启动，生成随机密钥
		key := make([]byte, 32)
		_, _ = rand.Read(key)
		secret = base64.StdEncoding.EncodeToString(key)
		_ = s.setConfigValue(ctx, model.ConfigKeyJWTSecret, secret)
	}

	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		// 降级使用原始字符串
		s.jwtSecret = []byte(secret)
		return
	}
	s.jwtSecret = decoded
}

// GetStatus 获取认证状态
func (s *sAuth) GetStatus(ctx context.Context) (*authApi.GetStatusRes, error) {
	pinHash := s.getConfigValue(ctx, model.ConfigKeyPINHash)
	twoFAEnabled := s.getConfigValue(ctx, model.ConfigKey2faEnabled)
	return &authApi.GetStatusRes{
		PinSet:       pinHash != "",
		TwoFAEnabled: twoFAEnabled == "true",
	}, nil
}

// Setup 首次设置 PIN
//
// 功能说明:
// 1. 检查 PIN 是否已设置（不可重复设置）
// 2. 验证 PIN 格式
// 3. bcrypt 哈希存储
// 4. 自动登录，返回 JWT Token
func (s *sAuth) Setup(ctx context.Context, pin string) (*authApi.SetupRes, error) {
	// 检查是否已设置
	existingHash := s.getConfigValue(ctx, model.ConfigKeyPINHash)
	if existingHash != "" {
		return nil, gerror.New("PIN 已设置，请使用修改 PIN 功能")
	}

	// 验证 PIN 格式
	if !pinPattern.MatchString(pin) {
		return nil, gerror.Newf("PIN 格式错误：必须为 %d-%d 位数字", model.PINMinLength, model.PINMaxLength)
	}

	// bcrypt 哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.Wrap(err, "PIN 哈希生成失败")
	}

	// 存储哈希
	if err := s.setConfigValue(ctx, model.ConfigKeyPINHash, string(hash)); err != nil {
		return nil, gerror.Wrap(err, "保存 PIN 失败")
	}

	// 生成 JWT Token（设置成功后自动登录）
	token, err := s.generateToken()
	if err != nil {
		return nil, gerror.Wrap(err, "生成 Token 失败")
	}

	g.Log().Info(ctx, "PIN 设置成功")

	return &authApi.SetupRes{Token: token}, nil
}

// Login 验证 PIN 返回 JWT Token
//
// 功能说明:
// 1. 检查 PIN 是否已设置
// 2. 检查账户是否被锁定
// 3. 验证 PIN
// 4. 成功则清除失败计数，生成 JWT Token
// 5. 失败则累加失败计数，达到上限锁定
func (s *sAuth) Login(ctx context.Context, pin string) (*authApi.LoginRes, error) {
	// 检查是否已设置 PIN
	hash := s.getConfigValue(ctx, model.ConfigKeyPINHash)
	if hash == "" {
		return nil, gerror.New("PIN 未设置，请先设置 PIN")
	}

	// 检查锁定状态
	if s.isLocked(ctx) {
		return nil, gerror.New("账户已锁定，请稍后再试")
	}

	// 验证 PIN
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin)); err != nil {
		// 记录失败
		s.recordFailure(ctx)
		return nil, gerror.New("PIN 错误")
	}

	// 验证成功，清除失败计数
	s.clearFailures(ctx)

	// 检查是否启用了 2FA
	twoFaEnabled := s.getConfigValue(ctx, model.ConfigKey2faEnabled)
	if twoFaEnabled == "true" {
		// 生成临时 JWT Token，附带 2fa_pending 声明
		token, err := s.generate2FAPendingToken()
		if err != nil {
			return nil, gerror.Wrap(err, "生成 2FA 临时 Token 失败")
		}
		g.Log().Info(ctx, "PIN 校验成功，等待 2FA 验证")
		// 前端收到 login res 后，如果有临时 token 可能需要特殊状态处理
		// 我们直接把 token 给前端，前端调用受限于中间件的实际授权（见 middleware 调整）
		return &authApi.LoginRes{Token: token}, nil
	}

	// 生成完全授权 JWT Token
	token, err := s.generateToken()
	if err != nil {
		return nil, gerror.Wrap(err, "生成 Token 失败")
	}

	g.Log().Info(ctx, "PIN 登录成功")

	return &authApi.LoginRes{Token: token}, nil
}

// ChangePin 修改 PIN
//
// 功能说明:
// 1. 验证旧 PIN
// 2. 验证新 PIN 格式
// 3. 生成新 PIN 哈希并存储
func (s *sAuth) ChangePin(ctx context.Context, currentPin string, newPin string) (*authApi.ChangePinRes, error) {
	// 获取存储的哈希
	hash := s.getConfigValue(ctx, model.ConfigKeyPINHash)
	if hash == "" {
		return nil, gerror.New("PIN 未设置，请先设置 PIN")
	}

	// 验证旧 PIN
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(currentPin)); err != nil {
		return nil, gerror.New("旧 PIN 错误")
	}

	// 验证新 PIN 格式
	if !pinPattern.MatchString(newPin) {
		return nil, gerror.Newf("新 PIN 格式错误：必须为 %d-%d 位数字", model.PINMinLength, model.PINMaxLength)
	}

	// 生成新哈希
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPin), bcrypt.DefaultCost)
	if err != nil {
		return nil, gerror.Wrap(err, "PIN 哈希生成失败")
	}

	// 存储新哈希
	if err := s.setConfigValue(ctx, model.ConfigKeyPINHash, string(newHash)); err != nil {
		return nil, gerror.Wrap(err, "保存新 PIN 失败")
	}

	g.Log().Info(ctx, "PIN 修改成功")

	return &authApi.ChangePinRes{Success: true}, nil
}

// generateToken 生成完全授权 JWT Token
func (s *sAuth) generateToken() (string, error) {
	claims := jwt.MapClaims{
		"sub": "allfi-user",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(model.TokenExpireHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// generate2FAPendingToken 生成 2FA 临时 JWT Token，仅用于发起 Verify2FA 请求
func (s *sAuth) generate2FAPendingToken() (string, error) {
	claims := jwt.MapClaims{
		"sub":         "allfi-user",
		"2fa_pending": true,
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(time.Duration(model.TwoFATokenExpireMinutes) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// Setup2FA 获取 2FA 配置（密钥与 QR Url）
func (s *sAuth) Setup2FA(ctx context.Context) (*authApi.Setup2FARes, error) {
	// 如果已启用，则不允许重新 setup，需要先 disable
	if s.getConfigValue(ctx, model.ConfigKey2faEnabled) == "true" {
		return nil, gerror.New("2FA 已启用，请先禁用后再重新配置")
	}

	// 生成新的 TOTP 密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "AllFi",
		AccountName: "admin",
		SecretSize:  20,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "生成 2FA 密钥失败")
	}

	// 将密钥临时存入 DB 供 Enable2FA 验证（此时状态尚未启用）
	secret := key.Secret()
	if err := s.setConfigValue(ctx, model.ConfigKey2faSecret, secret); err != nil {
		return nil, gerror.Wrap(err, "临时保存 2FA 密钥失败")
	}

	g.Log().Info(ctx, "生成 2FA TOTP 密钥成功")

	return &authApi.Setup2FARes{
		Secret: secret,
		QrUrl:  key.URL(),
	}, nil
}

// Enable2FA 验证并启用 2FA
func (s *sAuth) Enable2FA(ctx context.Context, code string) (*authApi.Enable2FARes, error) {
	secret := s.getConfigValue(ctx, model.ConfigKey2faSecret)
	if secret == "" {
		return nil, gerror.New("尚未配置 2FA，请先调用 Setup")
	}

	// 验证 TOTP Code
	valid := totp.Validate(code, secret)
	if !valid {
		return nil, gerror.New("2FA 验证码错误")
	}

	// 验证成功，将其标记为启用
	if err := s.setConfigValue(ctx, model.ConfigKey2faEnabled, "true"); err != nil {
		return nil, gerror.Wrap(err, "启用 2FA 失败")
	}

	g.Log().Info(ctx, "2FA 启用成功")
	return &authApi.Enable2FARes{Success: true}, nil
}

// Disable2FA 验证并禁用 2FA
func (s *sAuth) Disable2FA(ctx context.Context, code string) (*authApi.Disable2FARes, error) {
	if s.getConfigValue(ctx, model.ConfigKey2faEnabled) != "true" {
		return nil, gerror.New("未启用 2FA，无需禁用")
	}

	secret := s.getConfigValue(ctx, model.ConfigKey2faSecret)
	if secret == "" {
		return nil, gerror.New("未找到 2FA 密钥")
	}

	// 验证 TOTP Code
	valid := totp.Validate(code, secret)
	if !valid {
		s.recordFailure(ctx) // 防爆破记录
		return nil, gerror.New("2FA 验证码错误")
	}

	// 验证成功，清除启用标记和密钥
	_ = s.setConfigValue(ctx, model.ConfigKey2faEnabled, "false")
	_ = s.setConfigValue(ctx, model.ConfigKey2faSecret, "")

	g.Log().Info(ctx, "2FA 禁用成功")
	return &authApi.Disable2FARes{Success: true}, nil
}

// Verify2FA 登录流程中的 2FA 验证
func (s *sAuth) Verify2FA(ctx context.Context, code string) (*authApi.Verify2FARes, error) {
	// （通常在前置中间件中会验证是否有 token，或者在这里我们可以假设调用此接口表示正在登录过程中）
	// 注意：因为是两步验证，可能只有 2fa_pending 的 token 时才允许调这个接口。在 middleware 里处理比较好，
	// 此处仅进行 Code 校验和签发新 Token。

	if s.getConfigValue(ctx, model.ConfigKey2faEnabled) != "true" {
		return nil, gerror.New("2FA 未启用，无需验证")
	}

	if s.isLocked(ctx) {
		return nil, gerror.New("账户已锁定，请稍后再试")
	}

	secret := s.getConfigValue(ctx, model.ConfigKey2faSecret)
	if secret == "" {
		return nil, gerror.New("未找到 2FA 密钥，配置异常")
	}

	// 验证 TOTP Code
	valid := totp.Validate(code, secret)
	if !valid {
		s.recordFailure(ctx) // 输错也算认证失败，纳入冷却机制
		return nil, gerror.New("2FA 验证码错误")
	}

	// 验证成功，清除错误计数
	s.clearFailures(ctx)

	// 签发完全授权的 JWT Token
	token, err := s.generateToken()
	if err != nil {
		return nil, gerror.Wrap(err, "生成授权 Token 失败")
	}

	g.Log().Info(ctx, "2FA 验证成功，登录完成")
	return &authApi.Verify2FARes{
		Success: true,
		Token:   token,
	}, nil
}

// isLocked 检查账户是否被锁定
func (s *sAuth) isLocked(ctx context.Context) bool {
	lockUntil := s.getConfigValue(ctx, model.ConfigKeyLockUntil)
	if lockUntil == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339, lockUntil)
	if err != nil {
		return false
	}
	return time.Now().Before(t)
}

// recordFailure 记录认证失败
func (s *sAuth) recordFailure(ctx context.Context) {
	countStr := s.getConfigValue(ctx, model.ConfigKeyFailCount)
	count, _ := strconv.Atoi(countStr)
	count++
	_ = s.setConfigValue(ctx, model.ConfigKeyFailCount, strconv.Itoa(count))

	// 达到上限则锁定
	if count >= model.MaxFailCount {
		lockTime := time.Now().Add(time.Duration(model.LockDurationMinutes) * time.Minute)
		_ = s.setConfigValue(ctx, model.ConfigKeyLockUntil, lockTime.Format(time.RFC3339))
		g.Log().Warning(ctx, "认证失败次数达到上限，账户已锁定",
			"failCount", count,
			"lockUntil", lockTime.Format(time.RFC3339),
		)
	}
}

// clearFailures 清除失败计数
func (s *sAuth) clearFailures(ctx context.Context) {
	_ = s.setConfigValue(ctx, model.ConfigKeyFailCount, "0")
	_ = s.setConfigValue(ctx, model.ConfigKeyLockUntil, "")
}

// getConfigValue 从 system_config 表读取配置值
func (s *sAuth) getConfigValue(ctx context.Context, key string) string {
	var config systemEntity.SystemConfig
	err := systemDao.SystemConfig.Ctx(ctx).
		Where(systemDao.SystemConfig.Columns().ConfigKey, key).
		Scan(&config)
	if err != nil {
		return ""
	}
	return config.ConfigValue
}

// setConfigValue 写入 system_config 表
func (s *sAuth) setConfigValue(ctx context.Context, key string, value string) error {
	// 检查是否存在
	count, err := systemDao.SystemConfig.Ctx(ctx).
		Where(systemDao.SystemConfig.Columns().ConfigKey, key).
		Count()
	if err != nil {
		return gerror.Wrap(err, "查询配置失败")
	}

	if count > 0 {
		// 更新
		_, err = systemDao.SystemConfig.Ctx(ctx).
			Where(systemDao.SystemConfig.Columns().ConfigKey, key).
			Data(g.Map{
				systemDao.SystemConfig.Columns().ConfigValue: value,
			}).
			Update()
	} else {
		// 插入
		_, err = systemDao.SystemConfig.Ctx(ctx).Insert(g.Map{
			systemDao.SystemConfig.Columns().ConfigKey:   key,
			systemDao.SystemConfig.Columns().ConfigValue: value,
			systemDao.SystemConfig.Columns().Description: fmt.Sprintf("认证配置: %s", key),
		})
	}
	return err
}
