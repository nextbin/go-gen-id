package rpc_grpc

import (
	"github.com/nextbin/go-gen-id/src/server/component/gen"
	"github.com/nextbin/go-gen-id/src/server/rpc/rpc_grpc/pb"
	"golang.org/x/net/context"
)

type GenIdService struct{}

func (s *GenIdService) GenId(ctx context.Context, request *pb.GenIdRequest) (response *pb.GenIdResponse, err error) {
	response = &pb.GenIdResponse{}
	response.Id, response.Code, err = gen.GenId()
	return
}
