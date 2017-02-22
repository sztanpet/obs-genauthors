#!/bin/sh
cd data
go generate
cd ..
GOOS=windows GOARCH=amd64 go build -o obs-genauthors_win-amd64.exe
go build
echo "build.sh done"
