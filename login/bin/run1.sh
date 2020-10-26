#! /bin/bash

docker run -d  -p 14101:14222/tcp braid-game/login:latest \
    -consul http://172.17.0.1:8500 \
    -jaeger http://172.17.0.1:14268/api/traces \
    -nsqlookup 172.17.0.1:4161 \
    -nsqd 172.17.0.1:4150 \
    -redis redis://172.17.0.1:6379/0
