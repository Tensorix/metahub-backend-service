package authpage

import (
	"os"
	"time"

	auth "github.com/Tensorix/metahub-backend-service/gen/proto/v1/auth"
	"google.golang.org/grpc"
)

var (
	secretKey []byte
	expireAt  = 24 * time.Hour
)

type User struct {
	Id       int32
	Username string
	Pwd      string
}

type server struct {
	auth.UnimplementedAuthServiceServer
}

func Register(s *grpc.Server) {
	var err error
	secretKey, err = os.ReadFile("secret.img")
	if err != nil {
		panic("please create secret.img")
	}
	auth.RegisterAuthServiceServer(s, &server{})
}
