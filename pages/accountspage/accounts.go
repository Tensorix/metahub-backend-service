package accountspage

import (
	"google.golang.org/grpc"

	accounts "github.com/Tensorix/metahub-backend-service/gen/proto/v1/accounts"
)

type server struct {
	accounts.UnimplementedAccountsServiceServer
}

func Register(s *grpc.Server) {
	accounts.RegisterAccountsServiceServer(s, &server{})
}
