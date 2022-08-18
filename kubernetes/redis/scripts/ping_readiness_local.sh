#!/bin/sh
set -e

REDIS_STATUS_FILE=/tmp/.redis_cluster_check
response=$(
    timeout -s 3 $1 redis-cli ping
)
if [ "$?" -eq "124" ]; then
    echo "Timed out"
    exit 1
fi
if [ "$response" != "PONG" ]; then
    echo "$response"
    exit 1
fi
if [ ! -f "$REDIS_STATUS_FILE" ]; then
    response=$(
        timeout -s 3 $1 \
            redis-cli CLUSTER INFO | grep cluster_state | tr -d '[:space:]'
    )
    if [ "$?" -eq "124" ]; then
        echo "Timed out"
        exit 1
    fi
    if [ "$response" != "cluster_state:ok" ]; then
        echo "$response"
        exit 1
    else
        touch "$REDIS_STATUS_FILE"
    fi
fi
