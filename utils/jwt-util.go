package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/x/errors"
	"time"
)

func GenJwt(id, name, jwtSecret, secret string) (string, error) {
	encryptId, err := Encrypt([]byte(id), []byte(secret))
	if err != nil {
		return "", errors.New(5001, "Error creating token")
	}
	encryptName, err := Encrypt([]byte(name), []byte(secret))
	if err != nil {
		return "", errors.New(5001, "Error creating token")
	}
	// 创建一个新的 JWT Claims 对象
	claims := jwt.MapClaims{
		"id":   string(encryptId),
		"name": string(encryptName),
		"exp":  time.Now().Add(time.Hour * 24 * 356).Unix(), // 设置过期时间为一年
	}

	// 使用私钥进行签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(jwtSecret)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", errors.New(5001, "Error creating token")
	}

	return tokenString, nil
}

func ParseJwt(tokenString, jwtSecret, secret, targetId, targetName string) error {
	// 使用公钥（用于验证签名）解析 Token
	key := []byte(jwtSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return errors.New(5001, fmt.Sprintf("Error parsing token: %s", err))
	}

	// 验证 Token 是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimId := claims["id"].(string)
		claimName := claims["name"].(string)
		//exp := claims["exp"].(float64)
		//expirationTime := time.Unix(int64(exp), 0)
		//fmt.Println("Token expires at:", expirationTime)
		decryptId, err := Decrypt([]byte(claimId), []byte(secret))
		if err != nil {
			return errors.New(5001, "Token is not valid")
		}
		decryptName, err := Decrypt([]byte(claimName), []byte(secret))
		if err != nil {
			return errors.New(5001, "Token is not valid")
		}
		if string(decryptId) != targetId || string(decryptName) != targetName {
			return errors.New(5001, "Token is not valid")
		}
		return nil
	} else {
		return errors.New(5001, "Token is not valid")
	}
}
