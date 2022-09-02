package authsample

import (
	"bytes"
	"context"
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

var appPolicies = map[string]string{
	"MyTour":      "100PerMin",
	"DaySchedule": "9000PerHour",
	"FuelPass":    "5PerSec",
	"MyPet":       "30PerMin",
	"PerfApp":     "6MPerMin",
}

var subPolicies = map[string]string{
	"MyTour/Hotels":          "400PerMin",
	"MyTour/Weather":         "300PerMin",
	"DaySchedule/PizzaShack": "2000PerHour",
	"DaySchedule/SMS":        "100PerMin",
	"DaySchedule/Weather":    "80PerMin",
	"FuelPass/SMS":           "100PerMin",
	"PetsCare/Pets":          "5KPerMin",
	"PetsCare2/Pets":         "1KPerMin",
	"PerfApp/PerfAPI":        "7MPerMin",
}

// New creates a new authorization server.
func New(users Users) envoy_service_auth_v3.AuthorizationServer {
	return &server{users}
}

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
func (s *server) Check(
	ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {

	authorization := req.Attributes.Request.Http.Headers["authorization"]
	extracted := strings.Fields(authorization)

	if len(extracted) == 2 && extracted[0] == "Bearer" {
		valid, user, app := s.users.Check(extracted[1])
		if valid {
			cluster := req.Attributes.ContextExtensions["prodClusterName"]
			api := req.Attributes.ContextExtensions["name"]
			subscription := app + "/" + api

			xRateLimit := getCustomPolicyName(req, user, app, api, subscription)

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
									Value: cluster,
								},
							},
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									Key:   "x-ratelimit-api-policy",
									Value: xRateLimit,
								},
							},
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									Key:   "x-ratelimit-application",
									Value: app,
								},
							},
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									Key:   "x-ratelimit-application-policy",
									Value: appPolicies[app],
								},
							},
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									Key:   "x-ratelimit-subscription",
									Value: subscription,
								},
							},
							{
								Append: &wrappers.BoolValue{Value: false},
								Header: &envoy_api_v3_core.HeaderValue{
									Key:   "x-ratelimit-subscription-policy",
									Value: subPolicies[subscription],
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
			return resp, nil
		}
	}

	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_PERMISSION_DENIED),
		},
	}, nil
}

func getCustomPolicyName(req *envoy_service_auth_v3.CheckRequest, user, app, api, subscription string) (policyName string) {
	policyName = "default"
	xForwardedFor := req.Attributes.Request.Http.Headers["x-forwarded-for"]

	if api == "PizzaShack" && user == "user1" {
		policyName = "c1"
	} else if api == "SMS" && checkRange(xForwardedFor) {
		policyName = "c1"
	} else if strings.Contains(req.Attributes.Request.Http.Path, "/country") {
		httpMethod := req.Attributes.Request.Http.Headers[":method"]
		switch httpMethod {
		case "GET":
			if isOffPeak() {
				policyName = "c1"
			} else if req.Attributes.Request.Http.Headers["foo"] == "bar" {
				policyName = "c2"
			}
		case "POST":
			if isOffPeak() && checkRange(xForwardedFor) {
				policyName = "c1"
			}
		}
	}
	return
}

func checkRange(ip string) bool {
	var ip1 = net.ParseIP("216.14.49.184")
	var ip2 = net.ParseIP("216.14.49.191")

	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, ip1) >= 0 && bytes.Compare(trial, ip2) <= 0 {
		return true
	}
	return false
}

func isOffPeak() bool {
	return time.Now().Hour() >= 18
}
