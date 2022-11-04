# SAMPLE

docker-compose -f docker-compose-example.yml up --build --remove-orphans

## 5 Per Min with Default Policy

Enforcer sets "x-wso2-ratelimit-api-policy" as "default" (Here we send it with curl for demo)
curl http://localhost:8888/pets/v2/order -H "x-wso2-ratelimit-api-policy: default" -i

## 5 Per Min with Conditional Policy

Enforcer sets "x-wso2-ratelimit-api-policy" as "c1" (Here we send it with curl for demo)
curl http://localhost:8888/pets/v2/order -H "x-wso2-ratelimit-api-policy: c1" -i
