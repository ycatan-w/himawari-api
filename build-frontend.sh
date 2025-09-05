#!/bin/bash
set -e

# === Colors ===
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BRIGHT_YELLOW="\033[2;93m"
CYAN="\033[0;36m"
BRIGHT_CYAN="\033[0;96m"
BLUE="\033[0;34m"
NC="\033[0m"

# === Build Variables ===
WEB_REPOSITORY="git@github.com:ycatan-w/himawari-web.git"
WEB_SRC_DIR=".web"
WEB_DST_DIR="internal/server/web"
NPM="npm"
GIT="git"
set +e
NPM_BIN="$(command -v $NPM 2> /dev/null)"
GIT_BIN="$(command -v $GIT 2> /dev/null)"
set -e

# === Output Helper ===
echo -e "${YELLOW}
!      _    _ _                                    _
!     | |  | (_)                                  (_)
!     | |__| |_ _ __ ___   __ ___      ____ _ _ __ _
!     |  __  | | '_ \` _ \ / _\` \ \ /\ / / _\` | '__| |
!     | |  | | | | | | | | (_| |\ V  V / (_| | |  | |
!     |_|  |_|_|_| |_| |_|\__,_| \_/\_/ \__,_|_|  |_|
!
!${NC}"

header() {
  echo
  echo -e "${BRIGHT_CYAN}==========[ ${BRIGHT_YELLOW}$(date +%H:%M:%S$)${BRIGHT_CYAN} • $1 ]==========${NC}"
}
subheader() {
  echo
  echo -e "${CYAN}--- $1${NC}"
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

if [[ -z "$NPM_BIN" ]] || [[ -z "$GIT_BIN" ]]; then
  echo
  fail "Missing required dependencies: ${GREEN}$NPM${NC}, ${GREEN}$GIT${NC}"
  exit 1
fi

# === Build Directories ===
header "Clone repository"
rm -rf $WEB_SRC_DIR $WEB_DST_DIR/*
$GIT_BIN clone $WEB_REPOSITORY $WEB_SRC_DIR --quiet
success "Repository cloned"

header "Build web app"
cd $WEB_SRC_DIR
$NPM_BIN --silent install
$NPM_BIN --silent run build &> /dev/null
cd ../
success "Web application build"

header "Copy web app"
cp -r $WEB_SRC_DIR/dist/* $WEB_DST_DIR/.
success "Web application has been copy"

rm -rf $WEB_SRC_DIR

echo
success "Web app has been successfully added into the project"

