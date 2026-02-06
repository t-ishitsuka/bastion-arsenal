#!/usr/bin/env bash
# Arsenal installer for Linux/macOS
set -euo pipefail

# Configuration
REPO="t-ishitsuka/bastion-arsenal"
BINARY_NAME="bastion-arsenal"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${CYAN}==>${NC} $*"
}

log_success() {
    echo -e "${GREEN}✓${NC} $*"
}

log_error() {
    echo -e "${RED}✗${NC} $*" >&2
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $*"
}

detect_platform() {
    local os arch

    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    case "$os" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            log_error "サポートされていないOS: $os"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            log_error "サポートされていないアーキテクチャ: $arch"
            exit 1
            ;;
    esac

    log_info "検出されたプラットフォーム: ${OS}-${ARCH}"
}

get_latest_release() {
    log_info "最新リリース情報を取得中..."

    local api_url="https://api.github.com/repos/${REPO}/releases/latest"
    local release_info

    if command -v curl >/dev/null 2>&1; then
        release_info=$(curl -sL "$api_url")
    elif command -v wget >/dev/null 2>&1; then
        release_info=$(wget -qO- "$api_url")
    else
        log_error "curl または wget が必要です"
        exit 1
    fi

    VERSION=$(echo "$release_info" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        log_error "リリース情報の取得に失敗しました"
        exit 1
    fi

    log_success "最新バージョン: $VERSION"
}

download_archive() {
    local archive_name="${BINARY_NAME}-${VERSION}-${OS}-${ARCH}.tar.gz"
    local download_url="https://github.com/${REPO}/releases/download/${VERSION}/${archive_name}"
    local checksum_url="${download_url}.sha256"

    log_info "ダウンロード中: $archive_name"

    local tmp_dir
    tmp_dir=$(mktemp -d)
    cd "$tmp_dir"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL -o "$archive_name" "$download_url"
        curl -fsSL -o "${archive_name}.sha256" "$checksum_url" || true
    elif command -v wget >/dev/null 2>&1; then
        wget -q -O "$archive_name" "$download_url"
        wget -q -O "${archive_name}.sha256" "$checksum_url" || true
    fi

    # Verify checksum if available
    if [ -f "${archive_name}.sha256" ]; then
        log_info "チェックサムを検証中..."
        if command -v shasum >/dev/null 2>&1; then
            shasum -a 256 -c "${archive_name}.sha256" >/dev/null 2>&1
            log_success "チェックサム検証成功"
        elif command -v sha256sum >/dev/null 2>&1; then
            sha256sum -c "${archive_name}.sha256" >/dev/null 2>&1
            log_success "チェックサム検証成功"
        fi
    fi

    ARCHIVE_PATH="$tmp_dir/$archive_name"
    TMP_DIR="$tmp_dir"
}

extract_and_install() {
    log_info "展開中..."
    cd "$TMP_DIR"

    tar xzf "$ARCHIVE_PATH"

    local binary_file="${BINARY_NAME}-${OS}-${ARCH}"

    if [ ! -f "$binary_file" ]; then
        log_error "バイナリファイルが見つかりません: $binary_file"
        exit 1
    fi

    # Backup existing binary if it exists
    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        log_info "既存のバイナリをバックアップ中..."
        mv "${INSTALL_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}.backup"
    fi

    log_info "インストール中: ${INSTALL_DIR}/${BINARY_NAME}"

    # Try to install with sudo if needed
    if [ -w "$INSTALL_DIR" ]; then
        mv "$binary_file" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        log_warning "インストールには管理者権限が必要です"
        sudo mv "$binary_file" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    log_success "インストールが完了しました"
}

cleanup() {
    if [ -n "${TMP_DIR:-}" ] && [ -d "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}

verify_installation() {
    log_info "インストールを確認中..."

    if ! command -v "$BINARY_NAME" >/dev/null 2>&1; then
        log_warning "${BINARY_NAME} が PATH に見つかりません"
        log_warning "次のコマンドを実行して PATH に追加してください:"
        echo ""
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
        echo ""
        return
    fi

    local installed_version
    installed_version=$("$BINARY_NAME" version 2>/dev/null | head -n1 || echo "unknown")

    log_success "インストール済みバージョン: $installed_version"
    log_success "使用方法: ${BINARY_NAME} --help"
}

print_next_steps() {
    echo ""
    echo "次のステップ:"
    echo ""
    echo "  1. シェルの設定に Arsenal を追加:"
    echo "     eval \"\$(${BINARY_NAME} init-shell bash)\"  # Bash の場合"
    echo "     eval \"\$(${BINARY_NAME} init-shell zsh)\"   # Zsh の場合"
    echo "     ${BINARY_NAME} init-shell fish | source      # Fish の場合"
    echo ""
    echo "  2. ツールをインストール:"
    echo "     ${BINARY_NAME} install node 20.10.0"
    echo ""
    echo "  3. .toolversions から同期:"
    echo "     ${BINARY_NAME} sync"
    echo ""
}

main() {
    log_info "Arsenal インストーラー"
    echo ""

    detect_platform
    get_latest_release
    download_archive
    extract_and_install
    cleanup
    verify_installation
    print_next_steps
}

# Trap to cleanup on error
trap cleanup EXIT

main
