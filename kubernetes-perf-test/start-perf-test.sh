USERS_COUNT=500

cd ~/apache-jmeter-5.5/bin/
CC_IP=10.224.1.2


######## Backend Direct

# nohup ./jmeter -n -t ~/test/backend.jmx \
#     -Gprotocol=http -Ghost="${CC_IP}" -Gport=9999 -Gpath="/abc"  \
#     -Gusers=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/backend/${USERS_COUNT}Users.jtl" &
# tail -f nohup.out

######## Envoy with Backend

# nohup ./jmeter -n -t ~/test/envoy-with-backend.jtl \
#     -Gprotocol=https -Ghost="${CC_IP}" -Gport=8888 -Gpath="/perfapi/2.1.1/perf"  \
#     -Gusers=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/envoy_with_backend/${USERS_COUNT}Users.jtl" &
# tail -f nohup.out

######## Ratelimit + OSS Redis + Full Setup

nohup ./jmeter -n -t ~/test/ratelimit-one-route-api-level-ratelimit.jtl \
    -Gprotocol=https -Ghost="${CC_IP}" -Gport=8888 -Gpath="/perfapi/2.1.1/perf"  \
    -Gusers=$(($USERS_COUNT / 2)) -Gduration=900 \
    -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
    -l "/home/renuka/test-results/jtl/perf/ratelimit-oss/${USERS_COUNT}Users-API-level.jtl" &
tail -f nohup.out




######## With Ratelimits OSS

# nohup ./jmeter -n -t ~/test/PerformanceTest.jmx -GccIP="${CC_IP}" -GuserCount=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/oss-with-ratelimit/${USERS_COUNT}Users.jtl" &
# tail -f nohup.out

######## With Ratelimits Azure

# nohup ./jmeter -n -t ~/test/PerformanceTest.jmx -GccIP="${CC_IP}" -GuserCount=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/azure-with-ratelimit/${USERS_COUNT}Users.jtl" &
# tail -f nohup.out





######## Without Ratelimits

# nohup ./jmeter -n -t ~/test/PerformanceTestWithoutRateLimit.jmx -GccIP="${CC_IP}" -GuserCount=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/oss-without-ratelimit/${USERS_COUNT}Users.jtl" &
# tail -f nohup.out





######## One Per Day

# nohup ./jmeter -n -t ~/test/PizzaShack-APILevel-default-20KPerMin.jmx -GccIP="${CC_IP}" -GuserCount=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/one-per-day/${USERS_COUNT}Users.jtl" &
# tail -f nohup.out
