#!/bin/sh

echo "[Sentinel] Waiting for redis-master DNS to resolve..."

# Try to resolve redis-master:6379 in a loop (max 30s)
for i in $(seq 1 30); do
  nslookup redis-master && break
  echo "[$i] redis-master not ready, retrying..."
  sleep 1
done

echo "[Sentinel] Starting Redis Sentinel..."
exec redis-server /sentinel.conf --sentinel
