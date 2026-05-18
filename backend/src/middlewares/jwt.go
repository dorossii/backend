package middlewares

import (
	"backend/logger"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaim struct {
	UserID   string   // ユーザーID
	Name     string   // ユーザー名
	Email    string   // メールアドレス
	Labels   []string // ラベル
	ProvCode string   // プロバイダーコード
	ProvUid  string   // プロバイダーUID
}

func ValidateToken(tokenString string) (AccessTokenClaim, error) {
	logger.Println("トークンを検証します")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

	if err != nil {
		return AccessTokenClaim{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		labels := claims["labels"].([]interface{})

		return AccessTokenClaim{
			UserID:   claims["userID"].(string),
			Name:     claims["name"].(string),
			Email:    claims["email"].(string),
			Labels:   interfaceToString(labels),
			ProvCode: claims["provCode"].(string),
			ProvUid:  claims["provUid"].(string),
		}, nil
	} else {
		logger.PrintErr(err)
	}

	return AccessTokenClaim{}, err
}

func interfaceToString(values []interface{}) []string {
	var result []string
	for _, value := range values {
		result = append(result, value.(string))
	}

	return result
}
