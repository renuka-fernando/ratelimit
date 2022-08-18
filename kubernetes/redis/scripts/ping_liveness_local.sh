#!/bin/sh
set -e
response=$(
    timeout -s 3 $1 redis-cli ping
)
if [ "$?" -eq "124" ]; then
    echo "Timed out"
    exit 1
fi
responseFirstWord=$(echo $response | head -n1 | awk '{print $1;}')
if [ "$response" != "PONG" ] && [ "$responseFirstWord" != "LOADING" ] && [ "$responseFirstWord" != "MASTERDOWN" ]; then
    echo "$response"
    exit 1
fi
