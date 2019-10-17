package rpc_grpc

import (
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/server/rpc/rpc_grpc/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type RpcGrpcServer struct {
	Port int
}

func (s *RpcGrpcServer) Run() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGenServer(grpcServer, &GenIdService{})
	log.WithField("port", base.RpcGrpcPort).Info("starting server now...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
