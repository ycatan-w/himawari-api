@echo off
setlocal

:: === Variables ===
SET SERVICE_NAME=HimawariServer
SET INSTALL_DIR="C:\Program Files\Himawari"
SET BIN_NAME=himawari-server.exe
SET UNINSTALL_NAME="uninstall.bat"
SET SRC_DIR=%~dp0
SET BIN_PATH=%INSTALL_DIR%\%BIN_NAME%
set PORT=9740

:: === Verify and restart as admin ===
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Requesting administrator privileges...
    powershell -Command "Start-Process '%~f0' -Verb runAs"
    exit /b
)

echo === Installing Himawari Server at %DATE% %TIME% ===

:: Create bin dir
if not exist %INSTALL_DIR% (
    echo Creating install dir: %INSTALL_DIR%
    mkdir %INSTALL_DIR%
)

:: Copy binary
echo Copying binary...
copy /Y "%SRC_DIR%%BIN_NAME%" %INSTALL_DIR%
copy /Y "%SRC_DIR%%UNINSTALL_NAME%" %INSTALL_DIR%

:: Check binary
if not exist %BIN_PATH% (
    echo ERROR: Binary not found after copy: %BIN_PATH%
    exit /b 1
)

:: Remove Windows service if existing
sc.exe query %SERVICE_NAME% | findstr /I "FAILED 1060" >nul
if %errorlevel% equ 1 (
    echo Removing old service...
    sc.exe stop %SERVICE_NAME% >nul 2>&1
    sc.exe delete %SERVICE_NAME% >nul 2>&1
)

:: Create the service
echo Creating Windows service...
sc.exe create %SERVICE_NAME% binPath= %BIN_PATH% start= auto DisplayName= "Himawari Server"

:: Init DB
echo Initializing database...
%BIN_PATH% --init-db

:: Start the service
echo Starting service...
sc.exe start %SERVICE_NAME%

echo ✅ Installation complete.

pause
endlocal
