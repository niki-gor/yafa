#!/bin/sh

build() {
    docker build -t park .
}

start() {
    build
    docker run -d --memory 2G --log-opt max-size=5M --log-opt max-file=3 --name park_perf -p 5000:5000 -p 5432:5432 park
    ../tester func -u http://localhost:5000/api -r report.html
}

restart() {
    docker stop park_perf
    docker rm park_perf
    start
}

update_easyjson() {
    rm internal/model/*_easyjson.go
    easyjson -snake_case -all internal/model/*
}
