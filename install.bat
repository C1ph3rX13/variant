@echo off

echo init ...
go mod tidy

echo Installing go-wines
go install github.com/tc-hib/go-winres@latest

echo Installing goimports
go install golang.org/x/tools/cmd/goimports@latest

echo Installing garble
go install mvdan.cc/garble@latest

echo Done
pause