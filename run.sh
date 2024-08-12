#!/bin/sh

# This is just a personalized script to run the application in my own server.
# For production, it's strongly recommend to use a proper tool such as docker.

while true; do 
    go mod tidy
    go build .
    ./ExamSphere $1
    exit_code=$?
    if [ $exit_code -eq 50 ]; then
        echo "Exit code is 50. Breaking the loop."
        break
    fi

    # making sure that we are on the latest version
    git stash; git pull

    # go back and build the ui as well
    cd ../ExamSphere-ui
    git pull
    npm install && npm run build

    cp build/* ../ExamSphere/ui/ -r

    sleep 2
done