package main

import (
	"authorization/config"
	"authorization/domain"
	"authorization/infrastructure/worker"
	"authorization/view"
	"context"
	"log"
	"net"
	"strings"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_status "google.golang.org/grpc/status"
)

type AuthorizationServer struct {
}

// inject a header that can be used for future rate limiting
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	userID := req.Attributes.Request.Http.Headers["user-id"]
	method := req.Attributes.Request.Http.Method
	rawPath := req.Attributes.Request.Http.Path
	path := strings.Split(rawPath, "?")[0]

	log.Printf("authorization for user_id: %s to path %s and method %s", userID, path, method)

	endpoints := make(map[string]domain.Endpoint)
	isAuthorized, err := view.Authorization(ctx, userID, method, path, endpoints)
	if err != nil {
		log.Printf("Error while authorizing: %v", err)
		return nil, _status.Errorf(codes.Internal, "Error while authorizing: %v", err)
	}

	if isAuthorized {
		return &auth.CheckResponse{
			Status:       &status.Status{Code: int32(rpc.OK)},
			HttpResponse: &auth.CheckResponse_OkResponse{},
		}, nil
	}

	return &auth.CheckResponse{
		Status: &status.Status{Code: int32(rpc.PERMISSION_DENIED)},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Forbidden,
				},
				Body: "You are not authorized to access this resource",
			},
		},
	}, nil
}

func main() {
	// create a TCP listener on port 4000
	lis, err := net.Listen("tcp", config.AppConfig.AppHost+":"+config.AppConfig.AppExtAuthzPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	mailerClient := worker.CreateMailerClient()
	worker.CreateMailer(mailerClient)
	defer mailerClient.Close()

	grpcServer := grpc.NewServer()
	authServer := &AuthorizationServer{}
	auth.RegisterAuthorizationServer(grpcServer, authServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
