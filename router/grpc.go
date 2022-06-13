package router

import (
	"belajar-grpc/api/user"
	"belajar-grpc/configs"
	"belajar-grpc/package/errors"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	userpb "belajar-grpc/models/user"
)

func NewGRPCServer(config *configs.Config, logger *logrus.Logger) error {
	lis,err:=net.Listen("tcp","0.0.0.0:"+config.PORT_GRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v",err)
		return err
	}
	
	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	logrusEntry := logrus.NewEntry(logger)
	opts:=[]grpc_logrus.Option{grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}){
		return "grpc.time_ns",duration.Nanoseconds()
	}),}

	alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string,servingObject interface{}) bool{ return true}


	// register grpc service server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(grpc_ctxtags.UnaryServerInterceptor(),grpc_logrus.UnaryServerInterceptor(logrusEntry,opts...),
		RequestIDInterceptor(),AuthInterceptor,
		grpc_logrus.PayloadUnaryServerInterceptor(logrusEntry,alwaysLoggingDeciderServer))),
	)

	userpb.RegisterUsersServer(grpcServer,user.New(config,logger))
	
	// add reflection service
	reflection.Register(grpcServer)

	// running gRPC server
	log.Println("[SERVER] GRPC server is ready")
	grpcServer.Serve(lis)

	return nil
}

const (
    // XRequestIDKey is a key for getting request id.
    XRequestIDKey    = "x-user-id"
    unknownRequestID = "<unknown>"
)

// RequestIDInterceptor is a interceptor of access control list.
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
    return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
        requestID := requestIDFromContext(ctx)
        ctx = context.WithValue(ctx, "x-user-id", requestID)
        return handler(ctx, req)
    }
}

func requestIDFromContext(ctx context.Context) string {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return unknownRequestID
    }

    key := strings.ToLower(XRequestIDKey)
    header, ok := md[key]
    if !ok || len(header) == 0 {
        return unknownRequestID
    }

    requestID := header[0]
    if requestID == "" {
        return unknownRequestID
    }

    return requestID
}

// func AuthInterceptor()grpc.UnaryServerInterceptor {
// 	return authInterceptor
// }

// AuthInterceptor middleware validatetoken validation
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// register service to oauth
	//RegServiceToJWTToken(string(cocdpb.File_cocd_proto.Package()))

	// read header from incoming request
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		errn := errors.New("failed to read header from incoming request")
		return nil, errors.FormatError(errors.InternalServer, codes.Internal, false, "100", fmt.Sprintf("%v", errn),"")
	}
	if headers["x-api-key"][0]!="asdf" {
		errn := errors.New("API KEY not Valid")
		return nil, errors.FormatError(errors.InternalServer, codes.Internal, false, "100", fmt.Sprintf("%v", errn),"")
	}

	return handler(ctx, req)
}

func Map(structs interface{}) (results map[string]interface{}){
	data,_:=json.Marshal(structs)
	json.Unmarshal(data,&results)
	return
}