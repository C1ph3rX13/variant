@echo off

echo Configure GOPROXY Environment Variable
set "GOPROXY=https://goproxy.io,direct"

echo Set private repositories or groups that do not go through the proxy
set "GOPRIVATE=git.mycompany.com,github.com/my/private"

echo Initialize Go Modules
go mod init variant
go mod tidy

echo Installing go-wines
go install github.com/tc-hib/go-winres@latest

echo Installing goimports
go install golang.org/x/tools/cmd/goimports@latest

echo Installing garble
go install mvdan.cc/garble@latest

echo Initialization Completed

pause
