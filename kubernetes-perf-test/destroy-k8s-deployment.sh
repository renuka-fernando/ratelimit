alias k=kubectl
k delete -k ~/git/ratelimit/kubernetes/basic-tests/ratelimit-only
k delete -k ~/git/ratelimit/kubernetes/redis
k delete pvc --all
