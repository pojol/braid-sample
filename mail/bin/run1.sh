#! /bin/bash

docker run -d  -p 14301:14222/tcp braid-game/mail:latest \
    -consul http://172.17.0.1:8500 \
    -jaeger http://172.17.0.1:14268/api/traces
