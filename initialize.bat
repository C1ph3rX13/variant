@echo off

:: Set Error File
set "ERROR_LOG=errors.log"

echo [+] Check MinGW-w64 Installation
where gcc >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo [-] MinGW-w64 is not installed. Please install it from:
    echo [*] https://sourceforge.net/projects/mingw-w64/files/mingw-w64/mingw-w64-release/
    exit /b %ERRORLEVEL%
)
echo [+] MinGW-w64 is installed.

echo [+] Configure Go Environment
set CGO_ENABLED=1
set GO111MODULE=on

echo [+] Configure GOPROXY Environment Variable
set "GOPROXY=https://goproxy.io,direct"

echo [+] Set private repositories or groups that do not go through the proxy
set "GOPRIVATE=git.mycompany.com,github.com/my/private"

echo [+] Initialize Go Modules
go mod init variant 2>"%ERROR_LOG%"
if %ERRORLEVEL% neq 0 (
    echo [-] Error initializing Go module. Exiting.
    type "%ERROR_LOG%"
    exit /b %ERRORLEVEL%
)

echo [+] Tidy Go Modules
go mod tidy 2>"%ERROR_LOG%"
if %ERRORLEVEL% neq 0 (
    echo [-] Error tidying Go modules. Exiting.
    exit /b %ERRORLEVEL%
)

echo [+] Installing go-wines
go install github.com/tc-hib/go-winres@latest 2>"%ERROR_LOG%"
if %ERRORLEVEL% neq 0 (
    echo [-] Error installing go-wines. Exiting.
    type "%ERROR_LOG%"
    exit /b %ERRORLEVEL%
)

echo [+] Installing goimports
go install golang.org/x/tools/cmd/goimports@latest 2>"%ERROR_LOG%"
if %ERRORLEVEL% neq 0 (
    echo [-] Error installing goimports. Exiting.
    type "%ERROR_LOG%"
    exit /b %ERRORLEVEL%
)

echo [+] Installing garble
go install mvdan.cc/garble@master 2>"%ERROR_LOG%"
if %ERRORLEVEL% neq 0 (
    echo [-] Error installing garble. Exiting.
    type "%ERROR_LOG%"
    exit /b %ERRORLEVEL%
)

:: 检查 errors.log 是否为空并删除
if exist "%ERROR_LOG%" (
    for %%A in ("%ERROR_LOG%") do (
        if %%~zA neq 0 (
            echo [+] Deleted non-empty log file.
        ) else (
            echo [+] Log file was empty.
            del /f "%ERROR_LOG%"
        )
    )
)

echo [+] Initialization Completed

pause
