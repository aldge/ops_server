package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	jwt "github.com/golang-jwt/jwt/v5"

	passport "github.com/aldge/gopkg/casdoor"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenValid            = errors.New("未知错误")
	TokenExpired          = errors.New("token已过期")
	TokenNotValidYet      = errors.New("token尚未激活")
	TokenMalformed        = errors.New("这不是一个token")
	TokenSignatureInvalid = errors.New("无效签名")
	TokenInvalid          = errors.New("无法处理此token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GVA_CONFIG.JWT.SigningKey),
	}
}

func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	bf, _ := ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	ep, _ := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	claims := request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},                   // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    global.GVA_CONFIG.JWT.Issuer,              // 签名的发行者
		},
	}
	return claims
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenValid
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

// GetAccessToken 获取 AccessToken，优先从缓存获取，如果不存在则通过 passport 获取并缓存
func GetAccessToken() (string, error) {
	// 从配置中获取用户名
	userName := global.GVA_CONFIG.Passport.ClientUser
	if userName == "" {
		userName = "admin" // 默认值
	}

	// 如果 Redis 已启用且已初始化，尝试从缓存中获取 AccessToken
	if global.GVA_CONFIG.System.UseRedis && global.GVA_REDIS != nil {
		accessToken, err := global.GVA_REDIS.Get(context.Background(), userName).Result()
		if err == nil && accessToken != "" {
			// 缓存中存在，直接返回
			return accessToken, nil
		}
	}

	// 缓存中不存在，从 passport 获取
	// 读取证书文件内容
	certPath := global.GVA_CONFIG.Passport.Certificate
	if certPath == "" {
		certPath = "./micro_service.pem" // 默认路径
	}
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate file: %v", err)
	}
	certificate := string(certBytes)

	// 从配置中获取 passport 参数
	endpoint := global.GVA_CONFIG.Passport.EndPoint
	clientID := global.GVA_CONFIG.Passport.ClientID
	clientSecret := global.GVA_CONFIG.Passport.ClientSecret
	application := global.GVA_CONFIG.Passport.Application
	if application == "" {
		application = "micro_service" // 默认值
	}

	// Initialize the SDK with your Casdoor instance configuration
	passport.InitConfig(
		endpoint,     // endpoint
		clientID,     // clientId
		clientSecret, // clientSecret
		certificate,  // certificate (x509 format) - PEM格式公钥内容（不是文件路径）
		application,  // organizationName
		application,  // applicationName
	)

	// 获取 OAuth Token
	token, err := passport.GetOAuthTokenWithClientCredentials(userName)
	if err != nil {
		return "", fmt.Errorf("failed to get OAuth token: %v", err)
	}

	// 如果 Redis 已启用且已初始化，将获取到的 AccessToken 写入缓存
	if global.GVA_CONFIG.System.UseRedis && global.GVA_REDIS != nil {
		err = SetRedisJWT(token.AccessToken, userName)
		if err != nil {
			// 即使缓存写入失败，也返回 token，但记录错误
			return token.AccessToken, fmt.Errorf("failed to cache token: %v", err)
		}
	}

	return token.AccessToken, nil
}
