#!/bin/bash

service sendmail start
sleep 5s
# service sendmail status

echo "Starting go test"
go test

sleep 30s
