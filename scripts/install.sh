#!/bin/sh
# WorkPulse Installer
# Usage: curl -sSL https://raw.githubusercontent.com/tauantcamargo/workpulse/main/scripts/install.sh | sh

set -e

REPO="tauantcamargo/workpulse"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="workpulse"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)  echo "darwin" ;;
        Linux*)   echo "linux" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)        error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)             error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest release version
get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/'
}

# Compute SHA256 checksum (cross-platform)
compute_sha256() {
    if command -v sha256sum > /dev/null; then
        sha256sum "$1" | awk '{print $1}'
    elif command -v shasum > /dev/null; then
        shasum -a 256 "$1" | awk '{print $1}'
    else
        error "Neither sha256sum nor shasum found. Cannot verify checksum."
    fi
}

main() {
    info "Installing WorkPulse..."

    OS=$(detect_os)
    ARCH=$(detect_arch)

    # Windows users should download manually or use go install
    if [ "$OS" = "windows" ]; then
        error "Windows installation via this script is not supported. Please download from https://github.com/${REPO}/releases or use 'go install github.com/${REPO}@latest'"
    fi

    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        error "Failed to get latest version. Please check your internet connection."
    fi

    info "Detected: OS=${OS}, ARCH=${ARCH}"
    info "Latest version: ${VERSION}"

    # Build download URL
    VERSION_NUM="${VERSION#v}"
    ARCHIVE_NAME="${BINARY_NAME}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

    info "Downloading ${ARCHIVE_NAME}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    # Download archive with HTTP status check
    if command -v curl > /dev/null; then
        HTTP_CODE=$(curl -sL -w "%{http_code}" "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE_NAME}")
        if [ "$HTTP_CODE" != "200" ]; then
            error "Download failed with HTTP status ${HTTP_CODE}. URL: ${DOWNLOAD_URL}"
        fi
    elif command -v wget > /dev/null; then
        if ! wget -q "$DOWNLOAD_URL" -O "${TMP_DIR}/${ARCHIVE_NAME}"; then
            error "Download failed. URL: ${DOWNLOAD_URL}"
        fi
    else
        error "Neither curl nor wget found. Please install one of them."
    fi

    # Download checksums file
    info "Verifying checksum..."
    if command -v curl > /dev/null; then
        HTTP_CODE=$(curl -sL -w "%{http_code}" "$CHECKSUM_URL" -o "${TMP_DIR}/checksums.txt")
        if [ "$HTTP_CODE" != "200" ]; then
            error "Failed to download checksums file (HTTP ${HTTP_CODE})"
        fi
    elif command -v wget > /dev/null; then
        if ! wget -q "$CHECKSUM_URL" -O "${TMP_DIR}/checksums.txt"; then
            error "Failed to download checksums file"
        fi
    fi

    # Verify checksum
    EXPECTED_CHECKSUM=$(grep "${ARCHIVE_NAME}" "${TMP_DIR}/checksums.txt" | awk '{print $1}')
    if [ -z "$EXPECTED_CHECKSUM" ]; then
        error "Could not find checksum for ${ARCHIVE_NAME} in checksums.txt"
    fi

    ACTUAL_CHECKSUM=$(compute_sha256 "${TMP_DIR}/${ARCHIVE_NAME}")

    if [ "$EXPECTED_CHECKSUM" != "$ACTUAL_CHECKSUM" ]; then
        error "Checksum verification failed!\nExpected: ${EXPECTED_CHECKSUM}\nGot: ${ACTUAL_CHECKSUM}"
    fi

    info "Checksum verified successfully"

    # Extract archive
    info "Extracting archive..."
    cd "$TMP_DIR"
    tar -xzf "$ARCHIVE_NAME"

    # Determine installation directory
    if [ -w "$INSTALL_DIR" ]; then
        TARGET_DIR="$INSTALL_DIR"
    else
        # Try user-local installation first
        USER_BIN="${HOME}/.local/bin"
        warn "Cannot write to ${INSTALL_DIR} (permission denied)"
        info "Installing to ${USER_BIN} instead"
        mkdir -p "$USER_BIN"
        TARGET_DIR="$USER_BIN"
    fi

    # Install binary
    info "Installing to ${TARGET_DIR}..."
    mv "${BINARY_NAME}" "${TARGET_DIR}/"
    chmod +x "${TARGET_DIR}/${BINARY_NAME}"

    # Verify installation
    if [ "$TARGET_DIR" = "$INSTALL_DIR" ] && command -v workpulse > /dev/null; then
        info "WorkPulse installed successfully!"
        info "Run 'workpulse' to start the Pomodoro timer."
    elif [ "$TARGET_DIR" = "${HOME}/.local/bin" ]; then
        info "WorkPulse installed successfully to ${TARGET_DIR}"
        # Check if ~/.local/bin is in PATH
        case ":$PATH:" in
            *":${HOME}/.local/bin:"*)
                info "Run 'workpulse' to start the Pomodoro timer."
                ;;
            *)
                warn "${TARGET_DIR} is not in your PATH."
                warn "Add it to your shell profile:"
                warn "  export PATH=\"\$HOME/.local/bin:\$PATH\""
                warn "Or run directly: ${TARGET_DIR}/${BINARY_NAME}"
                ;;
        esac
    else
        info "Installation completed."
        info "Run '${TARGET_DIR}/${BINARY_NAME}' to start the Pomodoro timer."
    fi
}

main "$@"
