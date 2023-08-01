package main

import (
	"authorization/config"
	"authorization/domain"
	"authorization/infrastructure/worker"
	"authorization/util"
	"authorization/view"
	"context"
	"net"
	"strings"

	"github.com/rs/zerolog/log"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/gogo/googleapis/google/rpc"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_status "google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
)

type AuthorizationServer struct {
	Endpoints map[string]domain.Endpoint
}

// inject a header that can be used for future rate limiting
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {
	userID := req.Attributes.Request.Http.Headers["user-id"]
	method := req.Attributes.Request.Http.Method
	rawPath := req.Attributes.Request.Http.Path
	path := strings.Split(rawPath, "?")[0]

	log.Printf("authorization for user_id: %s to path %s and method %s", userID, path, method)

	isAuthorized, err := view.Authorization(ctx, userID, method, path, a.Endpoints)
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
		log.Fatal().Caller().Err(err).Msgf("Failed to listen: %s", lis.Addr())
	}
	log.Info().Caller().Msgf("listening on %s", lis.Addr())

	mailerClient := worker.CreateMailerClient()
	worker.CreateMailer(mailerClient)
	defer mailerClient.Close()

	endpointDatas := util.ReadYAML("endpoints.yml")
	var endpointYAML domain.EndpointYAML
	err = yaml.Unmarshal(endpointDatas, &endpointYAML)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to unmarshal endpoint data")
	}

	endpoints := make(map[string]domain.Endpoint)
	for _, endpoint := range endpointYAML.Endpoints {
		endpointData := domain.NewEndpoint(endpoint.Name, endpoint.Path, endpoint.Method)
		endpoints[endpoint.Name] = endpointData
	}

	grpcServer := grpc.NewServer()
	authServer := &AuthorizationServer{Endpoints: endpoints}
	auth.RegisterAuthorizationServer(grpcServer, authServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to start ext-authorization server.")
	}
}
