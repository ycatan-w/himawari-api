#!/bin/bash
set -e

# === Colors ===
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
YELLOW_UNDERLINE="\033[4;33m"
BRIGHT_YELLOW="\033[2;93m"
CYAN="\033[0;36m"
BRIGHT_CYAN="\033[0;96m"
BLUE="\033[0;34m"
NC="\033[0m"

# === Output Helper ===
header() {
  echo
  echo -e "${BRIGHT_CYAN}==========[ ${BRIGHT_YELLOW}$(date +%H:%M:%S)${BRIGHT_CYAN} • $1 ]==========${NC}"
}
success() {
    echo -e "[${BRIGHT_YELLOW}$(date +%H:%M:%S)${NC}] ${GREEN}✔ ${NC} $1"
}

fail() {
    echo -e "[${BRIGHT_YELLOW}$(date +%H:%M:%S)${NC}] ${RED}✖ ${NC} $1"
}

warn() {
    echo -e "[${BRIGHT_YELLOW}$(date +%H:%M:%S)${NC}] ${YELLOW}‼ ${NC} $1"
}

info() {
    echo -e "[${BRIGHT_YELLOW}$(date +%H:%M:%S)${NC}] ${BLUE}ⓘ ${NC} $1"
}

if [ "$EUID" -ne 0 ]; then
  warn "This script requires to run with sudo: \n> ${CYAN}sudo $0${NC}"
  exit 1
fi

# === Installation process ===
PROJECT="himawari-server"
BIN_PATH="/usr/local/bin/$PROJECT"
SERVICE_NAME="himawari"

header "Install $PROJECT into bin directory"
cp "$PROJECT" "$BIN_PATH"
chmod +x "$BIN_PATH"
success "${GREEN}$PROJECT${NC} installed"

header "Install $PROJECT-uninstall into bin directory"
cp "uninstall.sh" "$BIN_PATH-uninstall"
chmod +x "$BIN_PATH-uninstall"
success "${GREEN}$PROJECT-uninstall${NC} installed"

header "Create $PROJECT database directory"
mkdir -p /var/lib/himawari
success "Directories created"

header "Initialize database"
"$BIN_PATH" --init-db
success "database created"

header "Install $SERVICE_NAME service"
systemctl stop $SERVICE_NAME 2>/dev/null || true
cp "$SERVICE_NAME.service" /etc/systemd/system/
systemctl daemon-reload
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME
success "${GREEN}$SERVICE_NAME${NC} installed"

echo
info "${GREEN}$PROJECT${NC} is accessible from ${YELLOW_UNDERLINE}http://localhost:9740${NC}"
