package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Tensorix/metahub-backend-service/pages/authpage"
	"github.com/Tensorix/metahub-backend-service/pages/notifypage"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	port = flag.Int("port", 9090, "The server port")
)

func main() {

	// GORM init
	db, err := gorm.Open(sqlite.Open("mbs.sqlite"), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}

	// New grpc server
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	// Create bots

	// Register start
	authpage.Register(s, db)
	notifypage.Register(s)
	// Register end

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
