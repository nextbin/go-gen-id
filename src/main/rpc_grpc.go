package main

import (
	"github.com/nextbin/go-gen-id/src/base"
	"github.com/nextbin/go-gen-id/src/service"
	"github.com/nextbin/go-gen-id/src/service/pb"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type server struct{}

func (s *server) GenId(ctx context.Context, request *service_pb.GenIdRequest) (response *service_pb.GenIdResponse, err error) {
	response = &service_pb.GenIdResponse{}
	response.Id, response.Code, err = service.GenId()
	return
}

func runRpcGrpc() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(base.RpcGrpcPort))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	service_pb.RegisterGenServer(s, &server{})
	log.WithField("port", base.RpcGrpcPort).Info("starting server now...")
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
