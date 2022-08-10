CGO_ENABLED=0 GOOS=linux go build -o server-linux -ldflags="-w -s" main/main.go
docker build -t ext-auth-renuka:master . 
rm server-linux
