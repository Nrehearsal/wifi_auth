package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var liveTime time.Duration = 5
var signatureKey []byte

type CustomClaims struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	IP       string `json:"ip"`
	Mac      string `json:"mac"`
	Level    int    `json:"level"`
	jwt.StandardClaims
}

//create a new token
func GenerateToken(uid int, username, ip, mac string, level int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(liveTime * time.Minute)

	claims := CustomClaims{
		uid,
		username,
		ip,
		mac,
		level,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "wifi_auth",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(signatureKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string) (*CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signatureKey, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
