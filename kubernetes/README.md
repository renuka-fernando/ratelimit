## Setup

### Build Images

Execute from root dir.

1.  Rate Limit Service (DO NOT USE THIS FOR TESTING)
    ```sh
    build-ratelimit-docker.sh (DO NOT USE THIS FOR TESTING)
    ```

2.  Ext Auth Service
    ```sh
    cd ext-auth-server
    ./build-ext-auth-docker.sh
    ```


### Setup Redis Cluster

Create resources in Kubernetes cluster.

```sh
kubectl apply -k redis
```

Initialize Redis cluster (This is a onetime).

```sh
kubectl exec -it redis-cluster-0 -- redis-cli --cluster create $(kubectl get pods  -l app=redis-cluster -o json | jq -r '.items | map(.status.podIP) | join(":6379 ")'):6379 --cluster-replicas 1 -a password123
```

Test the cluster

```sh
kubectl exec -it redis-cluster-1 -- bash

redis-cli -c
keys *
```

### Setup Rate Limit Service and Other Components

```sh
kubectl apply -k .
```

### Test Whole Setup

```
curl localhost:8080/json -d '{
  "domain": "rl",
  "descriptors": [
    { "entries": [{ "key": "org", "value": "John" }, {"key":"resource","value":"/foo"}, {"key":"method", "value":"ALL"},{"key":"policy", "value":"3PerMin"}, {"key":"condition", "value":"default"}] }
  ]
}'
```

### Delete

#### Redis cluster

```sh
kubectl delete pvc --all
```
