package main

import (
	"context"
	"github.com/nextbin/go-gen-id/src/server/rpc/rpc_grpc/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"testing"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
}

func TestRpcGrpc(t *testing.T) {
	addr := "localhost:12001"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGenClient(conn)
	res, err := client.GenId(context.Background(), &pb.GenIdRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"code": res.Code, "id": res.Id}).Info("request rpc-grpc success")
}
