package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/accountspage"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
	"github.com/Tensorix/metahub-backend-service/pages/friendpage"
	"github.com/Tensorix/metahub-backend-service/pages/notifypage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	port = flag.Int("port", 9090, "The server port")
)

type Account struct {
	Id         int32
	AccountTag string
	UserId     uint
	IP         string
	Port       int
}

func main() {
	flag.Parse()
	fmt.Printf(("Starting server on port %d\n"), *port)
	// GORM init
	db, err := gorm.Open(sqlite.Open("mbs.sqlite"), &gorm.Config{TranslateError: true, Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	onebot.DB = db

	// New grpc server
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	// Create bots
	registerBot()
	// Register start
	reflection.Register(s)
	authpage.Register(s)
	notifypage.Register(s)
	friendpage.Register(s)
	accountspage.Register(s)
	// Register end

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func registerBot() {
	var accounts []Account
	if err := onebot.DB.Find(&accounts).Error; err != nil {
		log.Println(err)
	}

	for _, account := range accounts {
		var user authpage.User
		if err := onebot.DB.First(&user, "id = ?", account.UserId).Error; err != nil {
			log.Println(err)
		}
		onebot.NewOnebot(user.Username, account.AccountTag, account.IP, account.Port, user.Id, account.Id)
	}
}
