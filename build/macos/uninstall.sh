#!/bin/bash
set -euo pipefail

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
SERVICE_NAME="com.himawari.server"
PLIST_NAME="$SERVICE_NAME.plist"
PLIST_DST="/Library/LaunchDaemons/$PLIST_NAME"
BIN_PATH="/usr/local/bin/$PROJECT"
DB_NAME="himawari.db"
DB_PATH="/var/lib/himawari"
LOGS="/var/log/himawari-server*.log"
SCRIPTS_PATH=$(echo "$(readlink -f "$0")")

header "Stop $SERVICE_NAME service"
if launchctl list | grep -q "$SERVICE_NAME"; then
    launchctl unload "$PLIST_DST" || true
    success "${GREEN}$SERVICE_NAME${NC} stopped"
else
  warn "The service ${GREEN}$SERVICE_NAME${NC} was not found."
fi

header "Remove $PLIST_NAME LaunchDaemon"
if [ -f "$PLIST_DST" ]; then
    rm -f "$PLIST_DST"
    success "${GREEN}$PLIST_NAME${NC} removed"
else
  warn "The LaunchDaemon ${GREEN}$PLIST_NAME${NC} was not found."
fi

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

header "Remove logs"
rm -f $LOGS 2>/dev/null || true
success "Logs removed"

echo
info "${GREEN}$PROJECT${NC} has been removed"
info "Run ${CYAN}rm $SCRIPTS_PATH${NC} to remove this script"
