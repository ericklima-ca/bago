#!/bin/bash
export RABBITMQ_URL=amqp://root:root@localhost:5672
export JWT_SECRET=jwt_secret_test
export REDIS_URL=redis://localhost:6379
go test $1 ./tests/...