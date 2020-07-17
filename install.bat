@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set GOPROXY=https://goproxy.cn
set GO111MODULE=off

gofmt -w -s .

go build -o server/server gofile/server

:end
echo finished