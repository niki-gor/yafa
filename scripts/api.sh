#!/bin/sh

run_api() {
    until nc -vz localhost 5432; do 
        echo waiting for db...
        sleep 1
    done

    echo starting API...
    /api
}

run_api &