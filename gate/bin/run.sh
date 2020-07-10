#! /bin/bash

docker run -d  -p 14001:14222/tcp braid-game/gateway:latest \
    -consul http://172.17.0.1:8900 \
    -jaeger http://172.17.0.1:14268/api/traces
