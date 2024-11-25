package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Tensorix/metahub-backend-service/pages/authpage"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	port = flag.Int("port", 9090, "The server port")
)

func main() {

	// GORM init
	db, err := gorm.Open(sqlite.Open("mbs.sqlite"),&gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err.Error())
	}
	// Register start
	authpage.Register(s,db)
	// Register end

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
