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
  echo -e "${BRIGHT_CYAN}==========[ ${BRIGHT_YELLOW}$(date +%H:%M:%S$)${BRIGHT_CYAN} • $1 ]==========${NC}"
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

# === Removing process ===
PROJECT="himawari-server"
BIN_PATH="/usr/local/bin/$PROJECT"
SERVICE_NAME="himawari"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"
DB_NAME="himawari.db"
DB_PATH="/var/lib/himawari"
SCRIPTS_PATH=$(echo "$(readlink -f "$0")")

header "Stop $SERVICE_NAME service"
systemctl stop $SERVICE_NAME || true
systemctl disable $SERVICE_NAME || true
success "${GREEN}$SERVICE_NAME${NC} stopped"

header "Remove $SERVICE_NAME service"
if [ -f "$SERVICE_FILE" ]; then
  rm -f $SERVICE_FILE
  systemctl daemon-reload
else
  warn "The service file ${GREEN}$SERVICE_NAME.service${NC} was not found."
fi
success "${GREEN}$SERVICE_NAME${NC} removed"

echo "Removing binaries and data..."
header "Remove $PROJECT binary"
if [ -f "$BIN_PATH" ]; then
    rm -f "$BIN_PATH"
    success "${GREEN}$PROJECT${NC} removed"
else
  warn "The binary ${GREEN}$PROJECT${NC} was not found."
fi

header "Remove $DB_NAME database"
if [ -d "$DB_PATH" ]; then
    rm -rf "$DB_PATH"
    success "${GREEN}$DB_NAME${NC} removed"
else
  warn "The database ${GREEN}$DB_NAME${NC} was not found."
fi

echo
info "${GREEN}$PROJECT${NC} has been removed"
info "Run ${CYAN}rm $SCRIPTS_PATH${NC} to remove this script"
