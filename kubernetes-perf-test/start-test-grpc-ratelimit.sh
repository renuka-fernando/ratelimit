ssh perf-client << EOF

cd ~/apache-jmeter-5.5/bin/
RATE_LIMIT_IP=10.224.1.3

USERS_COUNT=200 # USER COUNT
nohup ./jmeter -n -t ~/test/grpc-direct.jmx \
    -GuserCount=$(($USERS_COUNT / 2)) \
    -Gduration=900 \
    -GratelimitIP="${RATE_LIMIT_IP}" \
    -GratelimitPort=9081 \
    -GprotoDir=/home/renuka/proto-lib/envoy/service/ratelimit/v3 \
    -GprotoLib=/home/renuka/proto-lib \
    -l "/home/renuka/test-results/jtl/basic/ratelimit-direct/${USERS_COUNT}Users.jtl" \
    -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new &
tail -f nohup.out
EOF
