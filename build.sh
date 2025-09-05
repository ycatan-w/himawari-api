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
TIMESTAMP=$(git log -1 --format="%ct")
if date --version >/dev/null 2>&1; then
    DATE=$(date -u -d "@$TIMESTAMP" +"%Y-%m-%dT%H:%M:%S%Z")
else
    DATE=$(date -u -r "$TIMESTAMP" +"%Y-%m-%dT%H:%M:%SZ")
fi
HASH=$(git log -1 --format="%h")
REFS=$(git log -1 --format="%D")
# BRANCH_OR_TAG=$(echo "$REFS" | sed -E 's/.*-> //; s/tag: //; s/,.*//')
BRANCH_OR_TAG=$(echo "$REFS" | grep -o 'tag: [^,]*' | sed 's/tag: //')
if [ -z "$BRANCH_OR_TAG" ]; then
    BRANCH_OR_TAG=$(echo "$REFS" | sed -E 's/.*-> //; s/,.*//')
fi

PROJECT="himawari-server"
VERSION="$BRANCH_OR_TAG"
BIN_DIR="bin"
BUILD_DIR="dist"
TAR="tar"
ZIP="zip"
PKGBUILD="pkgbuild"
if [[ "$OSTYPE" == "darwin"* ]]; then
  TAR="gtar"
fi
set +e
TAR_BIN="$(command -v $TAR 2> /dev/null)"
ZIP_BIN="$(command -v $ZIP 2> /dev/null)"
PKGBUILD_BIN="$(command -v $PKGBUILD 2> /dev/null)"
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

# === Build Helper ===
generate_build_name() {
  printf "%s-%s-%s_%s" $PROJECT $VERSION $1 $2
}

compile_go() {
  os=$1
  arch=$2
  name=$(generate_build_name $os $arch)
  if [[ "$os" == "windows" ]]; then
    name="${name}.exe"
  fi
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-X main.version=$VERSION -X main.commit=$HASH -X main.date=$DATE" -o $BIN_DIR/${name} ./cmd/server
  success "Compiled: $name"
}

package_linux() {
  arch=$1
  name=$(generate_build_name "linux" "$arch")
  pkg_name=$(printf "%s.tar.gz" $name)

  cp $BIN_DIR/${name} $BUILD_DIR/linux/${PROJECT}
  cp build/linux/install.sh build/linux/uninstall.sh build/linux/himawari.service $BUILD_DIR/linux/
  chmod +x $BUILD_DIR/linux/install.sh $BUILD_DIR/linux/uninstall.sh
  cd $BUILD_DIR/linux
  $TAR_BIN czf ${pkg_name} ${PROJECT} install.sh uninstall.sh himawari.service
  rm ${PROJECT} install.sh uninstall.sh himawari.service
  cd ../../
  success "Packaged: ${pkg_name}"
}

package_macos_pkgbuild() {
  arch=$1
  name=$(generate_build_name "darwin" "$arch")
  pkg_name=$(printf "%s-%s-macos_%s.pkg" $PROJECT $VERSION $arch)

  mkdir -p $BUILD_DIR/macos-pkg/pkgroot/usr/local/bin
  mkdir -p $BUILD_DIR/macos-pkg/pkgroot/Library/LaunchDaemons
  mkdir -p $BUILD_DIR/macos-pkg/Scripts
  cp $BIN_DIR/${name} $BUILD_DIR/macos-pkg/pkgroot/usr/local/bin/${PROJECT}
  cp build/macos-pkg/uninstall.sh $BUILD_DIR/macos-pkg/pkgroot/usr/local/bin/${PROJECT}-uninstall
  chmod +x $BUILD_DIR/macos-pkg/pkgroot/usr/local/bin/${PROJECT}-uninstall
  cp build/macos-pkg/com.himawari.server.plist $BUILD_DIR/macos-pkg/pkgroot/Library/LaunchDaemons/
  cp build/macos-pkg/preinstall build/macos-pkg/postinstall $BUILD_DIR/macos-pkg/Scripts/
  chmod +x $BUILD_DIR/macos-pkg/Scripts/*
  $PKGBUILD_BIN --root $BUILD_DIR/macos-pkg/pkgroot \
          --scripts $BUILD_DIR/macos-pkg/Scripts \
          --identifier com.himawari.server \
          --version $VERSION \
          $BUILD_DIR/macos-pkg/$pkg_name &> /dev/null
  rm -r $BUILD_DIR/macos-pkg/Scripts $BUILD_DIR/macos-pkg/pkgroot
  success "Packaged: ${pkg_name}"
}

package_macos_tar_gz() {
  arch=$1
  name=$(generate_build_name "darwin" "$arch")
  pkg_name=$(printf "%s-%s-macos_%s.tar.gz" $PROJECT $VERSION $arch)

  cp $BIN_DIR/${name} $BUILD_DIR/macos/${PROJECT}
  cp build/macos/install.sh build/macos/uninstall.sh build/macos/com.himawari.server.plist $BUILD_DIR/macos/
  chmod +x $BUILD_DIR/macos/install.sh $BUILD_DIR/macos/uninstall.sh
  cd $BUILD_DIR/macos
  $TAR_BIN czf ${pkg_name} ${PROJECT} install.sh uninstall.sh com.himawari.server.plist
  rm ${PROJECT} install.sh uninstall.sh com.himawari.server.plist
  cd ../../
  success "Packaged: ${pkg_name}"
}

package_windows() {
  arch=$1
  name=$(generate_build_name "windows" "$arch")
  pkg_name=$(printf "%s.zip" $name)
  cp $BIN_DIR/${name}.exe $BUILD_DIR/windows/${PROJECT}.exe
  cp build/windows/install.bat build/windows/uninstall.bat build/windows/install.ps1 build/windows/uninstall.ps1 $BUILD_DIR/windows/
  cd $BUILD_DIR/windows
  $ZIP_BIN -r ${pkg_name} ${PROJECT}.exe install.bat uninstall.bat install.ps1 uninstall.ps1 &> /dev/null
  rm ${PROJECT}.exe install.bat uninstall.bat install.ps1 uninstall.ps1
  cd ../../
  success "Packaged: ${pkg_name}"
}

if [[ -z "$TAR_BIN" ]] || [[ -z "$ZIP_BIN" ]] || [[ -z "$PKGBUILD_BIN" ]]; then
  echo
  fail "Missing required dependencies: ${GREEN}$TAR${NC}, ${GREEN}$ZIP${NC}, ${GREEN}$PKGBUILD${NC}"
  exit 1
fi

# === Build Directories ===
header "Create build directories"
rm -rf $BIN_DIR $BUILD_DIR
mkdir -p $BIN_DIR $BUILD_DIR
mkdir -p $BUILD_DIR/linux
mkdir -p $BUILD_DIR/macos
if [[ "$OSTYPE" == "darwin"* ]]; then
  mkdir -p $BUILD_DIR/macos-pkg/pkgroot/usr/local/bin
  mkdir -p $BUILD_DIR/macos-pkg/pkgroot/Library/LaunchDaemons
  mkdir -p $BUILD_DIR/macos-pkg/Scripts
fi
mkdir -p $BUILD_DIR/windows
success "Directory created"

# === Go Compilation ===
header "Compile Go binaries"

  # === Linux ===
subheader "Linux compilation"
compile_go "linux" "amd64"
compile_go "linux" "arm64"
compile_go "linux" "386"

  # === MacOs ===
subheader "Macos compilation"
compile_go "darwin" "amd64"
compile_go "darwin" "arm64"
if [[ "$OSTYPE" == "darwin"* ]]; then
  lipo -create -output $BIN_DIR/$(generate_build_name "darwin" "all") \
      $BIN_DIR/$(generate_build_name "darwin" "amd64") \
      $BIN_DIR/$(generate_build_name "darwin" "arm64")
  success "Compiled: $(generate_build_name "darwin" "all")"
fi

  # === Windows ===
subheader "Windows compilation"
compile_go "windows" "amd64"
compile_go "windows" "386"

# === Platform packaging ===
header "Build platform packages"

  # === Linux ===
subheader "Linux packaging"
package_linux "amd64"
package_linux "arm64"
package_linux "386"

  # === MacOs ===
if [[ "$OSTYPE" == "darwin"* ]]; then
  subheader "Macos packaging (.pkg)"
  package_macos_pkgbuild "amd64"
  package_macos_pkgbuild "arm64"
  if [[ "$OSTYPE" == "darwin"* ]]; then
    package_macos_pkgbuild "all"
  fi
fi
subheader "Macos packaging (.tar.gz)"
package_macos_tar_gz "amd64"
package_macos_tar_gz "arm64"
if [[ "$OSTYPE" == "darwin"* ]]; then
  package_macos_tar_gz "all"
fi

  # === Windows ===
subheader "Windows packaging"
package_windows "amd64"
package_windows "386"

header "Remove build directories"
rm -rf $BIN_DIR
success "Directory removed"

echo
success "Build success"
