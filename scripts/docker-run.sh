#!/bin/bash

set -e
set -u

sudo docker ps -f name=http-client-test -qa | xargs sudo docker rm -f

sudo docker run -d -p '8000:8000' \
    -e 'ENV_NAME=test' \
    -e 'SERVICE_BASE_URL=http://ec2-54-197-30-116.compute-1.amazonaws.com:8001/service' \
    -e 'LOG_LEVEL=info' \
    -e 'HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST=100' \
    -e 'HTTP_CLIENT_MAX_IDLE_CONNS=100' \
    -e 'HTTP_CLIENT_DIALER_TIMEOUT_MS=500' \
    -e 'HTTP_CLIENT_DIALER_KEEPALIVE_MS=30000' \
    -e 'HTTP_CLIENT_IDLE_CONN_TIMEOUT_MS=90000' \
    -e 'HTTP_CLIENT_TLS_HANDSHAKE_TIMEOUT_MS=1000' \
    -e 'HTTP_CLIENT_EXPECT_CONTINUE_TIMEOUT_MS=1000' \
    -e 'HTTP_CLIENT_TIMEOUT_MS=1500' \
    http-client-test
