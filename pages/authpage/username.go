package authpage

import "github.com/golang-jwt/jwt/v5"

func GetUsername(tokenString string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return ""
	}
	username, err := token.Claims.GetIssuer()
	if err != nil {
		return ""
	}
	return username
}
