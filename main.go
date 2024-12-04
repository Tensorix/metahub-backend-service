package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Tensorix/metahub-backend-service/onebot"
	"github.com/Tensorix/metahub-backend-service/pages/authpage"
	"github.com/Tensorix/metahub-backend-service/pages/friendpage"
	"github.com/Tensorix/metahub-backend-service/pages/notifypage"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	port = flag.Int("port", 9090, "The server port")
)

type Account struct {
	Id         uint
	UID        uint
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
	authpage.Register(s)
	notifypage.Register(s)
	friendpage.Register(s)
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
		var srv Srv
		var user authpage.User
		if err := onebot.DB.First(&srv, "id = ?", account.SrvId).Error; err != nil {
			log.Println(err)
		}
		if err := onebot.DB.First(&user, "id = ?", account.UserId).Error; err != nil {
			log.Println(err)
		}
		bot := onebot.NewOnebot(user.Username, account.AccountTag, srv.IpAddr, srv.Port)
		bot.UID = account.UID
		bot.UserID = user.Id
		bot.AccountID = account.Id
		bot.SrvID = account.SrvId
		bot.Run()
		onebot.Bots = append(onebot.Bots, bot)
		log.Println("New bot:", bot.Username, bot.AccountTag, bot.IP, bot.Port)
	}
}
