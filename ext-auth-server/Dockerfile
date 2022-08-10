# FROM golang:alpine AS builder

# RUN apk --no-cache add make
# COPY . /app
# RUN make -C /app/grpc-service

FROM alpine:3.15.4
COPY server-linux /bin/server
CMD ["/bin/server", "-users", "/etc/users.json"]
