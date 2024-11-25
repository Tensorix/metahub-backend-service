package authpage

import (
	"os"
	"time"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	secretKey []byte
	expireAt  = 24 * time.Hour
	db *gorm.DB
)

type User struct{
	Username string
	Pwd string
}

type login struct {
	auth.UnimplementedAuthServiceServer
}

func Register(s *grpc.Server,gormdb *gorm.DB) {
	var err error
	secretKey, err = os.ReadFile("secret.img")
	if err != nil {
		panic("please create secret.img")
	}
	auth.RegisterAuthServiceServer(s, &login{})
	db = gormdb
}
