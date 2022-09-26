### Create/Destroy K8s Deployment

./create-k8s-deployment.sh
./destroy-k8s-deployment.sh

### Deploy/Invoke API

LB_IP=192.168.205.2

apictl mg add env ratelimit --adapter "https://${LB_IP}:9843"
apictl mg login ratelimit -u admin -p admin -k
apictl mg deploy api -f . -e ratelimit -k

TOKEN=$(curl -X POST "https://${LB_IP}:9095/testkey" -d "scope=read:pets" -H "Authorization: Basic YWRtaW46YWRtaW4=" -k -v -H "Host: cc-envoy")

```sh
curl "https://${LB_IP}:9095/perfapi/2.1.1/perf" -i \
    -H Host:cc-envoy \
    -H "Authorization:Bearer $TOKEN" \
    -H x-ratelimit-api-policy:default \
    -d '{"hello":"world"}' \
    -k
```