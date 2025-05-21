#!/bin/sh

# 等 redis-master 的 6379 端口准备好
until nc -z redis-master 6379; do
  echo "Waiting for redis-master..."
  sleep 1
done

echo "Starting Sentinel..."
exec redis-server /sentinel.conf --sentinel
