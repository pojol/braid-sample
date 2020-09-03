#! /bin/bash

docker run -d  -p 14101:14222/tcp braid-game/login:latest \
    -consul http://172.17.0.1:8900 \
    -jaeger http://172.17.0.1:14268/api/traces \
    -nsqlookup 172.17.0.1:4161 \
    -nsqd 172.17.0.1:4150
