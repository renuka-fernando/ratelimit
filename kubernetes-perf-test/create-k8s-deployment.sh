alias k=kubectl
k apply -k ~/git/ratelimit/kubernetes/redis
sleep 60
kubectl exec -it redis-cluster-0 -- redis-cli --cluster create $(kubectl get pods  -l app=redis-cluster -o json | jq -r '.items | map(.status.podIP) | join(":6379 ")'):6379 --cluster-replicas 0 -a password123 --cluster-yes
k apply -k ~/git/ratelimit/kubernetes
k get svc -w
