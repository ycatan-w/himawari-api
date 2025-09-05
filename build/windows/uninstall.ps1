<#
    Himawari Server Uninstaller (PowerShell)
#>

# === Restart as admin ===
if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "This script must be run as Administrator. Relaunching with elevated privileges..."
    Start-Process powershell -Verb RunAs -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`""
    exit
}

# === Colors ===
$COLOR_SUCCESS = "Green"
$COLOR_FAIL    = "Red"
$COLOR_WARN    = "Yellow"
$COLOR_INFO    = "Blue"
$COLOR_HEADER  = "Cyan"

$ICON_SUCCESS = "[OK]"
$ICON_FAIL    = "[FAIL]"
$ICON_WARN    = "[WARN]"
$ICON_INFO    = "[INFO]"

# === Output Helpers ===
function Header($msg) {
    Write-Host ""
    Write-Host ("==========[ {0} | {1} ]==========" -f (Get-Date -Format "HH:mm:ss"), $msg) -ForegroundColor $COLOR_HEADER
}

function Success($msg) {
    Write-Host ("[{0}] {1} {2}" -f (Get-Date -Format "HH:mm:ss"), $ICON_SUCCESS, $msg) -ForegroundColor $COLOR_SUCCESS
}

function Fail($msg) {
    Write-Host ("[{0}] {1} {2}" -f (Get-Date -Format "HH:mm:ss"), $ICON_FAIL, $msg) -ForegroundColor $COLOR_FAIL
}

function Warn($msg) {
    Write-Host ("[{0}] {1} {2}" -f (Get-Date -Format "HH:mm:ss"), $ICON_WARN, $msg) -ForegroundColor $COLOR_WARN
}

function Info($msg) {
    Write-Host ("[{0}] {1} {2}" -f (Get-Date -Format "HH:mm:ss"), $ICON_INFO, $msg) -ForegroundColor $COLOR_INFO
}

# === Variables ===
$Project       = "himawari-server"
$BinName       = "$Project.exe"
$InstallDir    = "C:\Program Files\Himawari"
$BinPath       = Join-Path $InstallDir $BinName
$UninstallScript = Join-Path $InstallDir "uninstall.ps1"
$ServiceName   = "HimawariServer"
$DbDir         = "C:\ProgramData\Himawari"

# === Removing process ===
Header "Stopping $ServiceName service"
$svc = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if ($svc) {
    Stop-Service $ServiceName -Force
    Success "$ServiceName stopped"
} else {
    Warn "$ServiceName service not found"
}

Header "Removing $ServiceName service"
$svcExists = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if ($svcExists) {
    sc.exe delete $ServiceName | Out-Null
    Success "$ServiceName service removed"
} else {
    Warn "$ServiceName service does not exist"
}

Header "Removing $Project binary"
if (Test-Path $BinPath) {
    Remove-Item $BinPath -Force
    Success "$Project binary removed"
} else { Warn "$Project binary not found" }

if (Test-Path $UninstallScript) {
    Remove-Item $UninstallScript -Force
    Success "Uninstall script removed"
}

if (Test-Path $InstallDir) {
    Remove-Item $InstallDir -Recurse -Force
    Success "Install directory removed"
} else { Info "Install directory does not exist" }

Header "Removing database directory"
if (Test-Path $DbDir) {
    Remove-Item $DbDir -Recurse -Force
    Success "Database directory removed"
} else { Info "Database directory not found" }

Write-Host ""
Info "$Project has been fully uninstalled"
Info "You may now delete this script if desired"
pause
