#! /bin/bash

docker run -d  -p 14001:1202/tcp gateway:latest \
    -consul http://172.17.0.1:8900 \
    -jaeger http://172.17.0.1:14268/api/traces