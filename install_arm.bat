@echo off

setlocal

if exist install_arm.bat goto ok
echo install_arm.bat must be run from its folder
goto end

:ok

set GOPROXY=https://goproxy.cn
set GO111MODULE=off
set GOROOT=c:\go
set GOARCH=arm
set GOOS=linux

gofmt -w -s .

go build -o server/server gofile/server

:end
echo finished