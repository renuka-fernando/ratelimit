ghz --proto=/Users/renuka/git/envoy/api/envoy/service/ratelimit/v3/rls.proto \
    --call=envoy.service.ratelimit.v3.RateLimitService/ShouldRateLimit \
    -i /Users/renuka/Documents/envoy-rate-limit-poc/grpc-direct/proto-libs \
    --cacert=/Users/renuka/git/ratelimit/examples/certs/mg.pem \
    --cert=/Users/renuka/git/ratelimit/examples/certs/mg.pem \
    --key=/Users/renuka/git/ratelimit/examples/certs/mg.key --data '{
  "domain": "default",
  "descriptors": [
    { "entries": [{ "key": "org", "value": "John" }, {"key":"vhost","value":"cc-envoy"}, {"key":"vhost","value":"cc-envoy"}, {"key":"method", "value":"ALL"},{"key":"policy", "value":"3PerMin"}, {"key":"condition", "value":"default"}] }
  ]
}' localhost:8081

# -i /Users/renuka/git/envoy/api,/Users/renuka/git/udpa,/Users/renuka/git/protoc-gen-validate \
# -c 500 -n 100000 