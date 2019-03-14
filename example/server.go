package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if code := errDemo(req.Err); code != codes.OK {
		return nil, status.Error(code, "test error")
	}

	return &pb.Endpoint1Res{
		Bytes:   req.Bytes,
		Bool:    req.Bool,
		Int32:   req.Int32,
		Int64:   req.Int64,
		Msg:     req.Msg,
		String_: req.String_,
		Uint32:  req.Uint32,
	}, nil
}

func errDemo(enum pb.Error) codes.Code {
	return map[pb.Error]codes.Code{
		pb.Error_OK:                  codes.OK,
		pb.Error_CANCELLED:           codes.Canceled,
		pb.Error_UNKNOWN:             codes.Unknown,
		pb.Error_INVALID_ARGUMENT:    codes.InvalidArgument,
		pb.Error_DEADLINE_EXCEEDED:   codes.DeadlineExceeded,
		pb.Error_NOT_FOUND:           codes.NotFound,
		pb.Error_ALREADY_EXISTS:      codes.AlreadyExists,
		pb.Error_PERMISSION_DENIED:   codes.PermissionDenied,
		pb.Error_UNAUTHENTICATED:     codes.Unauthenticated,
		pb.Error_RESOURCE_EXHAUSTED:  codes.ResourceExhausted,
		pb.Error_FAILED_PRECONDITION: codes.FailedPrecondition,
		pb.Error_ABORTED:             codes.Aborted,
		pb.Error_OUT_OF_RANGE:        codes.OutOfRange,
		pb.Error_UNIMPLEMENTED:       codes.Unimplemented,
		pb.Error_INTERNAL:            codes.Internal,
		pb.Error_UNAVAILABLE:         codes.Unavailable,
		pb.Error_DATA_LOSS:           codes.DataLoss,
	}[enum]
}
