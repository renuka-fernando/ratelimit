## Setup

### Redis Cluster

<!-- OLD WAY - with bitnami
 ```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
# helm install cc-redis bitnami/redis -f redis/bitnami/sentinel-values.yaml
helm install cc-redis bitnami/redis-cluster -f redis/bitnami/cluster-values.yaml
kubectl apply -k .
``` -->

```sh
k apply -k redis
```

```sh
kubectl exec -it redis-cluster-0 -- redis-cli --cluster create $(kubectl get pods  -l app=redis-cluster -o json | jq -r '.items | map(.status.podIP) | join(":6379 ")'):6379 --cluster-replicas 1 -a password123
```

```sh
k run redis-cli --rm -it --image redis:7.0.4 -- bash
redis-cli -h 172.17.0.7 -c -a password123
```

### Other components

```sh
k apply -k .
```
