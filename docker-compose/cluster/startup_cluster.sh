#!/bin/sh

/opt/redis-stack/bin/redis-server /redis/redis.conf
sleep 5
echo hello world
echo yes | redis-cli --cluster create 192.168.57.10:6379 192.168.57.11:6379 192.168.57.12:6379
sleep infinity