CGO_ENABLED=0 GOOS=linux go build -o ratelimit-linux -ldflags="-w -s" -v github.com/envoyproxy/ratelimit/src/service_cmd
docker build -t ratelimit-renuka:master . 
rm ratelimit-linux
