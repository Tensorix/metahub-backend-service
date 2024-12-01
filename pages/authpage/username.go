package authpage

func GetUsername(token string) string {
	username, _ := verifyToken(token)
	return username
}
