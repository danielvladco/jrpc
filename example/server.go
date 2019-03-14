package main

import (
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/danielvladco/jrpc/example/pb"
)

func main() {
	svc := &server{}
	err := http.ListenAndServe(":8080", pb.ServiceHTTPServer(svc))
	if err != nil {
		log.Fatal(err.Error())
	}
}

type server struct{}

func (s *server) Endpoint1(ctx context.Context, req *pb.Endpoint1Req) (*pb.Endpoint1Res, error) {
	return &pb.Endpoint1Res{
		Param1: req.Param1,
		Param2: req.Param2,
		Param3: req.Param3,
		Param4: req.Param4,
		Param5: req.Param5,
		Param6: req.Param6,
		Param7: req.Param7,
		Param8: req.Param8,
	}, nil
}

func (s *server) Endpoint2(ctx context.Context, req *pb.Endpoint2Req) (*pb.Endpoint2Res, error) {
	return &pb.Endpoint2Res{
		Param1: req.Param1,
		Param2: req.Param2,
		Param3: req.Param3,
		Param4: req.Param4,
		Param5: req.Param5,
		Param6: req.Param6,
		Param7: req.Param7,
		Param8: req.Param8,
	}, nil
}

func (s *server) Endpoint3(ctx context.Context, req *pb.Endpoint3Req) (*pb.Endpoint3Res, error) {
	return &pb.Endpoint3Res{
		Param1: req.Param1,
		Param2: req.Param2,
		Param3: req.Param3,
		Param4: req.Param4,
		Param5: req.Param5,
		Param6: req.Param6,
		Param7: req.Param7,
		Param8: req.Param8,
	}, nil
}
