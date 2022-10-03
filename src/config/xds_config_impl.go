package config

import (
	"context"
	"io"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	pb_struct "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/ratelimit/src/settings"
	"github.com/envoyproxy/ratelimit/src/stats"
	"github.com/golang/protobuf/ptypes/any"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	rls_conf_v3 "github.com/envoyproxy/ratelimit/src/api/ratelimit/config/ratelimit/v3"
	rls_svc_v3 "github.com/envoyproxy/ratelimit/src/api/ratelimit/service/ratelimit/v3"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

const (
	configTypeURL string = "type.googleapis.com/ratelimit.config.ratelimit.v3.RateLimitConfig"
)

type rateLimitXdsGrpcConfigLoader struct {
	s                 settings.Settings
	xdsStream         rls_svc_v3.RateLimitConfigDiscoveryService_StreamRlsConfigsClient
	lastAckedResponse *discovery.DiscoveryResponse
	// TODO: (renuka) lastAckedResponse and lastReceivedResponse are equal
	lastReceivedResponse *discovery.DiscoveryResponse
	// If a connection error occurs, true event would be returned
	connectionFaultChannel chan bool
	// streamConfChannel chan
}

func (l *rateLimitXdsGrpcConfigLoader) Load(
	configs []RateLimitConfigToLoad, statsManager stats.Manager, mergeDomainConfigs bool) RateLimitConfig {
	return &RateLimitConfigXds{}
}

type RateLimitConfigXds struct{}

// Dump implements RateLimitConfig
func (rlc *RateLimitConfigXds) Dump() string {
	return "FOOO: impl xds"
}

// GetLimit implements RateLimitConfig
func (rlc *RateLimitConfigXds) GetLimit(ctx context.Context, domain string, descriptor *pb_struct.RateLimitDescriptor) *RateLimit {
	return &RateLimit{}
}

func (l *rateLimitXdsGrpcConfigLoader) getGrpcConnection() (*grpc.ClientConn, error) {
	backOff := grpc_retry.BackoffLinearWithJitter(l.s.ConfigGrpcXdsServerConnectRetryInterval, 0.5)
	logger.Infof("Dialing xDS Configuration Server: %q", l.s.ConfigGrpcXdsServerUrl)
	return grpc.Dial(
		l.s.ConfigGrpcXdsServerUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithTransportCredentials(generateTLSCredentialsForXdsClient()),
		grpc.WithBlock(),
		grpc.WithStreamInterceptor(
			grpc_retry.StreamClientInterceptor(grpc_retry.WithBackoff(backOff))))

}

func (l *rateLimitXdsGrpcConfigLoader) initConnection() (*grpc.ClientConn, error) {
	conn, err := l.getGrpcConnection()
	if err != nil {
		logger.Errorf("Error initializing gRPC connection to xDS Configuration Server: %s", err.Error())
		return nil, err
	}
	l.xdsStream, err = rls_svc_v3.NewRateLimitConfigDiscoveryServiceClient(conn).StreamRlsConfigs(context.Background())
	if err != nil {
		logger.Error("Error initializing gRPC stream to xDS Configuration Server: %s", err.Error())
		return nil, err
	}
	logger.Info("Connection to xDS Configuration Server is successful")
	return conn, nil
}

func (l *rateLimitXdsGrpcConfigLoader) initializeAndWatch() *grpc.ClientConn {
	conn, err := l.initConnection()
	if err != nil {
		l.connectionFaultChannel <- true
		return conn
	}
	go l.watchConfigs()

	// TODO: (renuka) check this, no nil for all cases
	var lastAppliedVersion string
	if l.lastAckedResponse != nil {
		// If the connection is interrupted in the middle, we need to apply if the version remains same
		lastAppliedVersion = l.lastAckedResponse.VersionInfo
	} else {
		lastAppliedVersion = ""
	}
	discoveryRequest := &discovery.DiscoveryRequest{
		Node:        &core.Node{Id: l.s.ConfigGrpcXdsNodeId},
		VersionInfo: lastAppliedVersion,
		TypeUrl:     configTypeURL,
	}
	l.xdsStream.Send(discoveryRequest)
	return conn
}

func (l *rateLimitXdsGrpcConfigLoader) applyConfigs(resources []*any.Any) {
	for _, res := range resources {
		config := &rls_conf_v3.RateLimitConfig{}
		err := anypb.UnmarshalTo(res, config, proto.UnmarshalOptions{}) // err := ptypes.UnmarshalAny(res, config)
		if err != nil {
			logger.Errorf("Error while unmarshalling config from xDS Configuration Server: %s", err.Error())
			l.nack(err.Error())
			return
		}

		logger.Info("RENUKA TEST: %v", config)
	}
	l.ack()
}

func (l *rateLimitXdsGrpcConfigLoader) watchConfigs() {
	for {
		discoveryResponse, err := l.xdsStream.Recv()
		if err == io.EOF {
			// reinitialize again, if stream ends
			logger.Error("EOF is received from xDS Configuration Server")
			l.connectionFaultChannel <- true
			return
		}
		if err != nil {
			logger.Error("Failed to receive the discovery response from xDS Configuration Server: %s", err.Error())
			errStatus, _ := grpcStatus.FromError(err)
			if errStatus.Code() == codes.Unavailable {
				logger.Errorf("Connection unavailable. errorCode: %s errorMessage: %s",
					errStatus.Code().String(), errStatus.Message())
				l.connectionFaultChannel <- true
				return
			}
			logger.Errorf("Error while xDS communication; errorCode: %s errorMessage: %s",
				errStatus.Code().String(), errStatus.Message())
			l.nack(errStatus.Message())
		} else {
			l.lastReceivedResponse = discoveryResponse
			logger.Debugf("Discovery response is received from xDS Configuration Server with response version: %s", discoveryResponse.VersionInfo)
			logger.Tracef("Discovery response received from xDS Configuration Server: %v", discoveryResponse)
			l.applyConfigs(discoveryResponse.Resources)
		}
	}
}

func (l *rateLimitXdsGrpcConfigLoader) ack() {
	l.lastAckedResponse = l.lastReceivedResponse
	discoveryRequest := &discovery.DiscoveryRequest{
		Node:          &core.Node{Id: l.s.ConfigGrpcXdsNodeId},
		VersionInfo:   l.lastAckedResponse.VersionInfo,
		TypeUrl:       configTypeURL,
		ResponseNonce: l.lastReceivedResponse.Nonce,
	}
	l.xdsStream.Send(discoveryRequest)
}

func (l *rateLimitXdsGrpcConfigLoader) nack(errorMessage string) {
	if l.lastAckedResponse == nil { // TODO: (renuka) why? if last acked response is nil, shouldn't we send discovery request?
		return
	}
	discoveryRequest := &discovery.DiscoveryRequest{
		Node:        &core.Node{Id: l.s.ConfigGrpcXdsNodeId},
		VersionInfo: l.lastAckedResponse.VersionInfo,
		TypeUrl:     configTypeURL,
		ErrorDetail: &status.Status{
			Message: errorMessage,
		},
	}
	if l.lastReceivedResponse != nil {
		discoveryRequest.ResponseNonce = l.lastReceivedResponse.Nonce
	}
	l.xdsStream.Send(discoveryRequest)
}

func (l *rateLimitXdsGrpcConfigLoader) initXdsClient() {
	logger.Info("Starting xDS client connection for rate limit configurations")
	conn := l.initializeAndWatch()
	for retryTrueReceived := range l.connectionFaultChannel {
		if !retryTrueReceived {
			continue
		}
		if conn != nil {
			conn.Close()
		}
		conn = l.initializeAndWatch()
	}
}

func NewRateLimitXdsGrpcConfigLoaderImpl(s settings.Settings) RateLimitConfigLoader {
	loader := &rateLimitXdsGrpcConfigLoader{s: s}
	loader.connectionFaultChannel = make(chan bool)
	loader.initXdsClient()
	return loader
}
