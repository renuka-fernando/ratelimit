## Setup
```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install cc-redis bitnami/redis -f redis/values.yaml
kubectl apply -k .
```
