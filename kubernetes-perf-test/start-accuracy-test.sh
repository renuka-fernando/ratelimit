USERS_COUNT=100

cd ~/apache-jmeter-5.5/bin/
CC_IP=10.224.1.2

nohup ./jmeter -n -t ~/test/PizzaShack-APILevel-default-20KPerMin.jmx -GccIP="${CC_IP}" -GuserCount=$(($USERS_COUNT / 2)) -Gduration=300 \
    -R envoy-rate-limit-perf-test-server-1-new,envoy-rate-limit-perf-test-server-2-new \
    -l "/home/renuka/test-results/jtl/accuracy/${USERS_COUNT}Users.jtl" &
tail -f nohup.out
