<#
    Himawari Server Installer (PowerShell)
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

# === Installation process ===
$Project     = "himawari-server"
$BinName     = "$Project.exe"
$SrcDir      = Split-Path -Parent $MyInvocation.MyCommand.Path
$InstallDir  = "C:\Program Files\Himawari"
$BinPath     = Join-Path $InstallDir $BinName
$ServiceName = "HimawariServer"
$UninstallScript = "uninstall.ps1"
$DbDir = "C:\ProgramData\Himawari"

Header "Installing $Project"

if (-not (Test-Path $InstallDir)) {
    Info "Creating install dir: $InstallDir"
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
    Success "Directory created"
} else { Info "Install directory already exists: $InstallDir" }

Info "Copying binary..."
Copy-Item (Join-Path $SrcDir $BinName) $BinPath -Force
if (Test-Path $BinPath) { Success "Binary installed at $BinPath" } else { Fail "Binary copy failed"; pause; exit 1 }

if (Test-Path (Join-Path $SrcDir $UninstallScript)) {
    Copy-Item (Join-Path $SrcDir $UninstallScript) (Join-Path $InstallDir $UninstallScript) -Force
    Success "Uninstall script installed"
} else { Warn "Uninstall script not found, skipping" }

if (-not (Test-Path $DbDir)) {
    Info "Creating database directory: $DbDir"
    New-Item -ItemType Directory -Path $DbDir | Out-Null
    Success "Database directory created"
} else { Info "Database directory already exists: $DbDir" }

Info "Initializing database..."
& $BinPath --init-db
Success "Database initialized"

$svc = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if ($svc) {
    Warn "Old service found, stopping and removing..."
    Stop-Service $ServiceName -Force
    sc.exe delete $ServiceName | Out-Null
    Start-Sleep 2
    Success "Old service removed"
}

Info "Creating Windows service..."
New-Service -Name $ServiceName -BinaryPathName $BinPath -DisplayName "Himawari Server" -StartupType Automatic
Success "Service created"

Info "Starting service..."
Start-Service $ServiceName
Success "Service started"

Write-Host ""
Info "Installation complete. Access $Project at http://localhost:9740"
pause
