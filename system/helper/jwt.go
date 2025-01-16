package helper

import (
	"book-fiber/system/config"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrExpiredToken = errors.New("token has expired")
    ErrInvalidTokenType = errors.New("invalid token type")
)
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// 用于存储解析后的token信息
type TokenClaims struct {
    Slug   string  `gorm:"slug"`
    Type   string `json:"type"`
    jwt.RegisteredClaims
}

// 解析访问令牌
func ParseAccessToken(tokenString string) (*TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 验证加密方法
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(config.C.JWT.AccessTokenSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
        // 验证token类型
        if claims.Type != "access" {
            return nil, jwt.ErrTokenInvalidClaims
        }
        return claims, nil
    }

    return nil, ErrInvalidToken
}

// 解析刷新令牌
func ParseRefreshToken(tokenString string) (*TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(config.C.JWT.RefreshTokenSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
        if claims.Type != "refresh" {
            return nil, jwt.ErrTokenInvalidClaims
        }
        return claims, nil
    }

    return nil, ErrInvalidToken
}

// ParseTokenAndGetSlug 解析token并获取slug
func ParseTokenAndGetSlug(token string) (string, error) {
    claims, err := ParseAccessToken(token)
    if err != nil {
        return "", err
    }
    return claims.Slug, nil
}

func GetTokenFromHeader(c *fiber.Ctx) (string, error) {
    auth := c.Get("Authorization")
    if auth == "" {
        return "", fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
    }

    parts := strings.Split(auth, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
    }

    return parts[1], nil
}

func GenerateTokenPair(slug string, rememberMe bool) (*TokenPair, error) {
    accessExpiration := 60 * time.Minute
    if rememberMe {
        accessExpiration = 180 * 24 * time.Hour
    }
    
    accessClaims := TokenClaims{
        Slug: slug,
        Type: "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    refreshClaims := TokenClaims{
        Slug: slug,
        Type: "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(180 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    
    accessTokenString, err := accessToken.SignedString([]byte(config.C.JWT.AccessTokenSecret))
    if err != nil {
        return nil, err
    }
    
    refreshTokenString, err := refreshToken.SignedString([]byte(config.C.JWT.RefreshTokenSecret))
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessTokenString,
        RefreshToken: refreshTokenString,
    }, nil
}

func ParseTokenAndGetUUID(token string) (uuid.UUID, error) {
    claims, err := ParseAccessToken(token)
    if err != nil {
        return uuid.Nil, err
    }
    // 将字符串转回 UUID
    return uuid.Parse(claims.Slug)
}