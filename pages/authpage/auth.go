package authpage

import (
	"os"
	"time"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"github.com/Tensorix/metahub-backend-service/onebot"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var (
	secretKey []byte
	expireAt  = 24 * time.Hour
	db        *gorm.DB
	bots      []onebot.Onebot
)

type User struct {
	Id       uint
	Username string
	Pwd      string
}

type server struct {
	auth.UnimplementedAuthServiceServer
}

func Register(s *grpc.Server, gormdb *gorm.DB, _bots []onebot.Onebot) {
	var err error
	secretKey, err = os.ReadFile("secret.img")
	if err != nil {
		panic("please create secret.img")
	}
	auth.RegisterAuthServiceServer(s, &server{})
	db = gormdb
	bots = _bots
}
