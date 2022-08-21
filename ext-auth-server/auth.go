package authsample

import (
	"bytes"
	"context"
	"log"
	"net"
	"strings"
	"time"

	envoy_api_v3_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type server struct {
	users Users
}

var _ envoy_service_auth_v3.AuthorizationServer = &server{}

// New creates a new authorization server.
func New(users Users) envoy_service_auth_v3.AuthorizationServer {
	return &server{users}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *server) Check(
	ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {

	log.Printf("Request size: %d\n", req.Attributes.Request.Http.Size)
	authorization := req.Attributes.Request.Http.Headers["authorization"]
	log.Println(authorization)

	xForwardedFor := req.Attributes.Request.Http.Headers["x-forwarded-for"]
	xRateLimit := "default" // default limit
	now := time.Now()
	if checkRange(xForwardedFor) {
		xRateLimit = "c1"
	} else if now.Month() == 8 && now.Day() == 4 {
		xRateLimit = "c2"
	} else if now.UTC().Hour() > 18 {
		xRateLimit = "c3"
	}

	extracted := strings.Fields(authorization)
	if len(extracted) == 2 && extracted[0] == "Bearer" {
		valid, user := s.users.Check(extracted[1])
		if valid {
			resp := &envoy_service_auth_v3.CheckResponse{
				HttpResponse: &envoy_service_auth_v3.CheckResponse_OkResponse{
					OkResponse: &envoy_service_auth_v3.OkHttpResponse{
						Headers: []*envoy_api_v3_core.HeaderValueOption{
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									// For a successful request, the authorization server sets the
									// x-current-user value.
									Key:   "x-current-user",
									Value: user,
								},
							},
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									Key:   "x-cluster-header",
									Value: "mock-sms",
								},
							},
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									Key:   "x-ratelimit-policy",
									Value: xRateLimit,
								},
							},
						},
					},
				},
				Status: &status.Status{
					Code: int32(code.Code_OK),
				},
				// DynamicMetadata: &structpb.Struct{
				// 	Fields: map[string]*structpb.Value{
				// 		"rate-limit": {Kind: &structpb.Value_StructValue{
				// 			StructValue: &structpb.Struct{
				// 				Fields: map[string]*structpb.Value{
				// 					"foo": {Kind: &structpb.Value_StringValue{StringValue: "bar"}},
				// 				},
				// 			},
				// 		}},
				// 	},
				// },
			}
			log.Printf("Response: %v", resp)
			return resp, nil
		}
	}

	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_PERMISSION_DENIED),
		},
	}, nil
}

func checkRange(ip string) bool {
	var ip1 = net.ParseIP("216.14.49.184")
	var ip2 = net.ParseIP("216.14.49.191")

	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		log.Printf("%v is not an IPv4 address\n", trial)
		return false
	}
	if bytes.Compare(trial, ip1) >= 0 && bytes.Compare(trial, ip2) <= 0 {
		log.Printf("%v is between %v and %v\n", trial, ip1, ip2)
		return true
	}
	log.Printf("%v is NOT between %v and %v\n", trial, ip1, ip2)
	return false
}
