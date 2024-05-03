#!/bin/bash

#go build -race -ldflags "-s -w"

export GIT_COMMIT=$(git log -1 --date=iso | grep "Date" | tr ' ' '_')
echo "build GIT:$GIT_COMMIT"
CGO_ENABLED=0 go build -ldflags "-X main.GIT=GIT:$GIT_COMMIT"
