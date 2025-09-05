@echo off
setlocal

SET SERVICE_NAME=HimawariServer
SET INSTALL_DIR="C:\Program Files\Himawari"
SET DATABASE_DIR="C:\ProgramData\Himawari"
SET BIN_NAME=himawari-server.exe
SET BIN_PATH=%INSTALL_DIR%\%BIN_NAME%

:: === Verify and restart as admin ===
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo Requesting administrator privileges...
    powershell -Command "Start-Process '%~f0' -Verb runAs"
    exit /b
)

echo === Uninstalling Himawari Server at %DATE% %TIME% ===

:: Stop service
sc.exe stop %SERVICE_NAME% >nul 2>&1

:: Delete service
sc.exe delete %SERVICE_NAME% >nul 2>&1

:: Remove install dir
if exist %INSTALL_DIR% (
    echo Removing install dir: %INSTALL_DIR%
    rmdir /S /Q %INSTALL_DIR%
)

:: Remove database dir
if exist %DATABASE_DIR% (
    echo Removing database dir: %DATABASE_DIR%
    rmdir /S /Q %DATABASE_DIR%
)

echo ✅ Uninstall complete.

pause
endlocal
