#!/usr/bin/env sh

for i in $(seq 1 $1); do curl --fail --silent -X POST http://localhost:8080/api/v1/shorten\?url\=https://google.com |> /dev/null; done