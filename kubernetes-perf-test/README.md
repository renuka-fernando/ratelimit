[Setup JMeter scripts](./setup)

## Start K8s Deployment

Apply azure patch if do with Azure Redis

```sh
./create-k8s-deployment.sh
```

Sample test
```sh
curl "https://10.224.1.2:8888/pizzashack/1.0.0/menu" -i \
    -H 'Authorization: Bearer user2.Today' \
    -H Host:cc-envoy \
    -H x-ratelimit-api-policy:default \
    -H x-cluster-header:clusterProd_cc-envoy_PizzaShack1.0.0 \
    -d '{"hello":"world"}' \
    -k

curl "https://10.224.1.2:8888/perfapi/2.1.1/perf" -i \
    -H Host:cc-envoy \
    -H x-cluster-header:clusterProd_cc-envoy_PerfAPI2.1.1 \
    -H x-ratelimit-api-policy:default \
    -d '{"hello":"world"}' \
    -k
```

## Start JMeter Servers

```sh
ssh perf-server-1 < start-server.sh
```

```sh
ssh perf-server-2 < start-server.sh
```

## Run Test

### 1. Watch Pods

```sh
k get po -w
```

### 2. Watch Resource Usage

```sh
for i in {1..15}; do sleep 60; k top po --containers & done
```

### 3. Run

Accuracy

```sh
ssh perf-client < start-accuracy-test.sh
```

Performance

```sh
ssh perf-client < start-perf-test.sh
```

### Direct Ratelimit Service

## Destroy K8s Deployment

```sh
./destroy-k8s-deployment.sh
```

## Evaluate

ssh perf-client
cd test-results/jtl/accuracy/


rm 1000Users.jtl filter* print-count.sh startTime.js

[Check](./evaluate)

```
java -jar jtl-splitter-0.4.6-SNAPSHOT.jar -t 5 -s -p -f 1000Users.jtl
```
