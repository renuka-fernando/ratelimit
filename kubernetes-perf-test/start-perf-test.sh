USERS_COUNT=200

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

# nohup ./jmeter -n -t ~/test/ratelimit-one-route-api-level-ratelimit.jtl \
#     -Gprotocol=https -Ghost="${CC_IP}" -Gport=8888 -Gpath="/perfapi/2.1.1/perf"  \
#     -Gusers=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/ratelimit-oss/${USERS_COUNT}Users-API-level.jtl" &
# tail -f nohup.out




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


#### Ratelimit One API

# nohup ./jmeter -n -t ~/test/ratelimit-one-route-api-level-ratelimit.jtl \
#     -Gprotocol=https -Ghost="${CC_IP}" -Gport=8888 -Gpath="/perfapi/2.1.1/perf"  \
#     -Gusers=$(($USERS_COUNT / 2)) -Gduration=900 \
#     -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
#     -l "/home/renuka/test-results/jtl/perf/ratelimit-oss/${USERS_COUNT}Users-API-level.jtl" &
# tail -f nohup.out


##### Chore Connect

alias apictl=/home/renuka/api/apictl
# apictl mg add env ratelimit --adapter https://10.224.1.2:9843
apictl mg login ratelimit -u admin -p admin -k
apictl mg deploy api -f /home/renuka/api/ratelimit-api -e ratelimit -k -o

TOKEN=$(curl -X POST "https://${CC_IP}:9195/testkey" -d "scope=read:pets" -H "Authorization: Basic YWRtaW46YWRtaW4=" -k -v -H "Host: cc-envoy")

# scp /Users/renuka/git/ratelimit/kubernetes/choreo-connect/choreo-connect-one-route-api-level-ratelimit.jtl perf-client:~/test
nohup ./jmeter -n -t ~/test/choreo-connect-one-route-api-level-ratelimit.jtl \
    -Gprotocol=https -Ghost="${CC_IP}" -Gport=9195 -Gpath="/perfapi/2.1.1/perf"  \
    -Gusers=$(($USERS_COUNT / 2)) -Gduration=900 \
    -Gauth="Bearer ${TOKEN}" \
    -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
    -l "/home/renuka/test-results/jtl/perf/ratelimit-oss/${USERS_COUNT}Users-API-level.jtl" &
tail -f nohup.out
