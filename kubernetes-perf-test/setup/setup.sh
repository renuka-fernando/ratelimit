set -e
scp ~/softwares/bin/apache-jmeter-5.3/bin/rmi_keystore.jks perf-client:~/apache-jmeter-5.5/bin/
scp ~/softwares/bin/apache-jmeter-5.3/bin/rmi_keystore.jks perf-server-1:~/apache-jmeter-5.5/bin/
scp ~/softwares/bin/apache-jmeter-5.3/bin/rmi_keystore.jks perf-server-2:~/apache-jmeter-5.5/bin/


scp ~/git/ratelimit/examples/accuracy-test/PizzaShack-APILevel-default-20KPerMin.jmx perf-client:~/test/PizzaShack-APILevel-default-20KPerMin.jmx

scp ~/git/ratelimit/examples/performance-test/PerformanceTest.jmx perf-client:~/test
scp ~/git/ratelimit/examples/performance-test/PerformanceTestWithoutRateLimit.jmx perf-client:~/test
# scp ~/git/ratelimit/examples/performance-test/PerformanceTestWithoutRateLimitActions.jmx perf-client:~/test
