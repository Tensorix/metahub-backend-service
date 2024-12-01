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
	secretKey    []byte
	expireAt     = 24 * time.Hour
	db           *gorm.DB
	Bots         []*onebot.Onebot
	registedList []string
	max_bot      = 10000
)

type User struct {
	Id       uint
	Username string
	Pwd      string
}

type Account struct {
	Id         uint
	AccountTag string
	UserId     uint
	SrvId      uint
}

type Srv struct {
	Id        uint
	ImgName   string
	Container string
	IpAddr    string
	Port      int
}

type server struct {
	auth.UnimplementedAuthServiceServer
}

func Register(s *grpc.Server, gormdb *gorm.DB) {
	var err error
	secretKey, err = os.ReadFile("secret.img")
	if err != nil {
		panic("please create secret.img")
	}
	auth.RegisterAuthServiceServer(s, &server{})
	db = gormdb
	Bots = make([]*onebot.Onebot, 10000)
}
