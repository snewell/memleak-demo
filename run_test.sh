#!/bin/sh

while true; do
    /app/memleak-demo client "${@}"
    sleep 5
done
