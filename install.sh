#!/bin/bash
set -euo pipefail

# TG-FileStreamBot Installer for macOS and Linux
# Usage: curl -sL https://filestream.bot/install | bash

# --------------- COLORS & STYLING ---------------
BOLD='\033[1m'
FSB_GREEN='\033[38;2;0;255;133m'    # #00FF85
FSB_DIM='\033[38;2;0;204;106m'      # #00cc6a
INFO='\033[38;2;136;146;176m'       # text-secondary
SUCCESS='\033[38;2;0;229;204m'      # cyan-bright
WARN='\033[38;2;255;176;32m'        # amber
ERROR='\033[38;2;230;57;70m'        # coral-mid
MUTED='\033[38;2;90;100;128m'       # text-muted
NC='\033[0m' # No Color

REPO="EverythingSuckz/TG-FileStreamBot"
GUM_VERSION="0.13.0"
TMPFILES=()
HAS_TTY="false"

if [[ -r /dev/tty ]]; then
    HAS_TTY="true"
fi

can_use_tui() {
    [[ -n "${GUM:-}" && "$HAS_TTY" == "true" ]]
}

prompt_read() {
    local __var_name="$1"
    local __prompt="$2"
    local __default="${3-}"
    local __input=""

    if [[ "$HAS_TTY" == "true" ]]; then
        read -r -p "$__prompt" __input < /dev/tty || __input=""
    else
        read -r -p "$__prompt" __input || __input=""
    fi

    if [[ -z "$__input" && -n "$__default" ]]; then
        __input="$__default"
    fi

    printf -v "$__var_name" '%s' "$__input"
}

gum_choose() {
    if [[ "$HAS_TTY" == "true" ]]; then
        "$GUM" choose "$@" < /dev/tty
    else
        "$GUM" choose "$@"
    fi
}

gum_confirm() {
    if [[ "$HAS_TTY" == "true" ]]; then
        "$GUM" confirm "$@" < /dev/tty
    else
        "$GUM" confirm "$@"
    fi
}

gum_input() {
    if [[ "$HAS_TTY" == "true" ]]; then
        "$GUM" input "$@" < /dev/tty
    else
        "$GUM" input "$@"
    fi
}

# --------------- TEMP RUNTIME ---------------
cleanup() {
    for f in "${TMPFILES[@]:-}"; do
        rm -rf "$f" 2>/dev/null || true
    done
}
trap cleanup EXIT

mktempdir() {
    local d
    d="$(mktemp -d)"
    TMPFILES+=("$d")
    echo "$d"
}

# --------------- LOGGING UI ---------------
ui_info() { echo -e "${MUTED}·${NC} $*"; }
ui_warn() { echo -e "${WARN}!${NC} $*"; }
ui_success() { echo -e "${SUCCESS}✓${NC} $*"; }
ui_error() { echo -e "${ERROR}✗${NC} $*" >&2; }

# --------------- DEPENDENCIES ---------------
DOWNLOADER=""
detect_downloader() {
    if command -v curl &> /dev/null; then
        DOWNLOADER="curl"
        return 0
    fi
    if command -v wget &> /dev/null; then
        DOWNLOADER="wget"
        return 0
    fi
    ui_error "Missing downloader. curl or wget required."
    exit 1
}

download_file() {
    local url="$1" output="$2"
    if [[ -z "$DOWNLOADER" ]]; then detect_downloader; fi
    if [[ "$DOWNLOADER" == "curl" ]]; then
        curl -fsSL --retry 3 -o "$output" "$url"
    else
        wget -q --tries=3 -O "$output" "$url"
    fi
}

download_github_api() {
    local url="$1" output="$2"
    if [[ -z "$DOWNLOADER" ]]; then detect_downloader; fi
    if [[ "$DOWNLOADER" == "curl" ]]; then
        curl -fsSL -H "Accept: application/vnd.github.v3+json" --retry 3 -o "$output" "$url"
    else
        wget -q --header="Accept: application/vnd.github.v3+json" --tries=3 -O "$output" "$url"
    fi
}

# --------------- GUM TUI BOOTSTRAP ---------------
GUM=""
bootstrap_gum() {
    if command -v gum &> /dev/null; then
        GUM="gum"
        return 0
    fi

    local os arch asset base gum_tmpdir
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    case "$(uname -m)" in
        x86_64|amd64) arch="x86_64" ;;
        arm64|aarch64) arch="arm64" ;;
        i386|i686) arch="i386" ;;
        *) return 1 ;;
    esac

    asset="gum_${GUM_VERSION}_${os}_${arch}.tar.gz"
    base="https://github.com/charmbracelet/gum/releases/download/v${GUM_VERSION}"
    gum_tmpdir="$(mktempdir)"

    if download_file "${base}/${asset}" "$gum_tmpdir/$asset"; then
        tar -xzf "$gum_tmpdir/$asset" -C "$gum_tmpdir" >/dev/null 2>&1
        local gum_path
        gum_path="$(find "$gum_tmpdir" -type f -name gum | head -n1)"
        if [[ -n "$gum_path" ]]; then
            chmod +x "$gum_path"
            GUM="$gum_path"
            return 0
        fi
    fi
    return 1
}

# --------------- DETECTION ---------------
detect_system() {
    OS="unknown"
    ARCH="unknown"
    
    case "$(uname -s)" in
        Darwin) OS="darwin" ;;
        Linux) OS="linux" ;;
        MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *)
            ui_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    if [[ "$OS" == "unknown" ]]; then
        ui_error "Unsupported operating system."
        exit 1
    fi
}

# --------------- API LOGIC ---------------
get_latest_version() {
    local include_prerelease="$1"
    local api_resp
    api_resp="$(mktempdir)/api.json"
    
    ui_info "Fetching latest releases from GitHub..." >&2
    download_github_api "https://api.github.com/repos/${REPO}/releases" "$api_resp"
    
    # We parse manually using robust grep/sed combos since jq might not be installed
    local latest_tag=""
    if [[ "$include_prerelease" == "yes" ]]; then
        latest_tag=$(grep -m 1 '"tag_name":' "$api_resp" | sed -E 's/.*"([^"]+)".*/\1/')
    else
        # Find first non-prerelease
        # This basic parsing finds the first release block where "prerelease": false
        # We use awk to grab the tag_name of the first object where prerelease is false
        latest_tag=$(awk '
            /"tag_name":/ { tag=$2; gsub(/["|,]/, "", tag) }
            /"prerelease": false/ { print tag; exit }
        ' "$api_resp")
    fi

    if [[ -z "$latest_tag" || "$latest_tag" == "null" ]]; then
        ui_error "Could not determine latest version from GitHub API."
        exit 1
    fi
    echo "$latest_tag"
}

# --------------- MAIN INSTALLER ---------------
main() {
    echo -e "${FSB_GREEN}${BOLD}"
    echo "  TG-FileStreamBot Installer"
    echo -e "${NC}${INFO}  Stream Telegram files with direct HTTP links.${NC}"
    echo ""

    detect_system
    bootstrap_gum || true

    # Existing Detection
    local EXISTING_FSB=""
    if command -v fsb &> /dev/null; then
        EXISTING_FSB="$(command -v fsb)"
    elif [[ -x "$PWD/fsb" ]]; then
        EXISTING_FSB="$PWD/fsb"
    elif [[ -x "$PWD/FSB/fsb" ]]; then
        EXISTING_FSB="$PWD/FSB/fsb"
    fi

    local SKIP_DIR_PROMPT="no"
    if [[ -n "$EXISTING_FSB" ]]; then
        local raw_ver
        raw_ver=$("$EXISTING_FSB" -v 2>/dev/null | grep -i 'version' | awk '{print $NF}' || true)
        local current_ver="unknown"
        if [[ -n "$raw_ver" && "$raw_ver" != "unknown" ]]; then
            if [[ "$raw_ver" != v* ]]; then
                current_ver="v${raw_ver}"
            else
                current_ver="$raw_ver"
            fi
        fi

        echo ""
        ui_info "Found existing installation at ${BOLD}$EXISTING_FSB${NC} (version: $current_ver)"
        
        local latest_stable
        latest_stable=$(get_latest_version "no")
        
        echo ""
        local ACTION
        if [[ "$current_ver" == "$latest_stable" ]]; then
            if can_use_tui; then
                ACTION=$(gum_choose "Keep installed version" "Reinstall $latest_stable" "Update to Pre-release" "Uninstall")
            else
                echo "1) Keep installed version"
                echo "2) Reinstall $latest_stable"
                echo "3) Update to Pre-release"
                echo "4) Uninstall"
                prompt_read act "Select action [1]: " "1"
                case "${act:-1}" in
                    2) ACTION="Reinstall $latest_stable" ;;
                    3) ACTION="Update to Pre-release" ;;
                    4) ACTION="Uninstall" ;;
                    *) ACTION="Keep installed version" ;;
                esac
            fi
        else
            ui_info "New version available: ${BOLD}$latest_stable${NC}"
            if can_use_tui; then
                ACTION=$(gum_choose "Update to latest stable ($latest_stable)" "Update to Pre-release" "Uninstall" "Cancel")
            else
                echo "1) Update to latest stable ($latest_stable)"
                echo "2) Update to Pre-release"
                echo "3) Uninstall"
                echo "4) Cancel"
                prompt_read act "Select action [1]: " "1"
                case "${act:-1}" in
                    2) ACTION="Update to Pre-release" ;;
                    3) ACTION="Uninstall" ;;
                    4) ACTION="Cancel" ;;
                    *) ACTION="Update to latest stable ($latest_stable)" ;;
                esac
            fi
        fi

        case "$ACTION" in
            "Keep installed version"|"Cancel")
                ui_info "Exiting."
                exit 0
                ;;
            "Uninstall")
                local target_dir
                target_dir=$(dirname "$EXISTING_FSB")
                local env_file="$target_dir/fsb.env"
                
                echo ""
                if [[ "$(basename "$target_dir")" == "FSB" ]]; then
                    ui_warn "The following directory and all its contents will be deleted:"
                    echo -e "  ${ERROR}$target_dir${NC}"
                else
                    ui_warn "The following files will be deleted:"
                    echo -e "  ${ERROR}$EXISTING_FSB${NC}"
                    if [[ -f "$env_file" ]]; then
                        echo -e "  ${ERROR}$env_file${NC}"
                    fi
                fi
                echo ""
                
                if can_use_tui; then
                    if ! gum_confirm "Proceed with uninstallation?" --default=false; then
                        ui_info "Uninstallation cancelled."
                        exit 0
                    fi
                else
                    prompt_read uninst_resp "Proceed with uninstallation? [y/N] " "N"
                    case "$uninst_resp" in
                        [yY][eE][sS]|[yY]) ;;
                        *)
                            ui_info "Uninstallation cancelled."
                            exit 0
                            ;;
                    esac
                fi

                ui_info "Uninstalling..."
                if [[ "$(basename "$target_dir")" == "FSB" ]]; then
                    rm -rf "$target_dir" 2>/dev/null || sudo rm -rf "$target_dir"
                else
                    rm -f "$EXISTING_FSB" 2>/dev/null || sudo rm -f "$EXISTING_FSB"
                    if [[ -f "$env_file" ]]; then
                        rm -f "$env_file" 2>/dev/null || sudo rm -f "$env_file"
                    fi
                fi
                ui_success "Uninstalled successfully."
                exit 0
                ;;
            "Update to Pre-release")
                WANT_PRERELEASE="yes"
                TARGET_VER=$(get_latest_version "yes")
                INSTALL_DIR=$(dirname "$EXISTING_FSB")
                SKIP_DIR_PROMPT="yes"
                WANT_ENV="no"
                ;;
            *)
                WANT_PRERELEASE="no"
                TARGET_VER="$latest_stable"
                INSTALL_DIR=$(dirname "$EXISTING_FSB")
                SKIP_DIR_PROMPT="yes"
                WANT_ENV="no"
                ;;
        esac
    fi

    # Defaults
    if [[ "$SKIP_DIR_PROMPT" != "yes" ]]; then
        WANT_PRERELEASE="No"
        INSTALL_DIR="$PWD/FSB"
        WANT_ENV="no"
        ENV_DIR=""

    # Interactive flow if TUI is available
    if can_use_tui; then
        echo -e "${BOLD}Do you want to install the latest pre-release version?${NC}"
        if gum_confirm "Include pre-releases? (Recommended: No)" --default=false; then
            WANT_PRERELEASE="yes"
        else
            WANT_PRERELEASE="no"
        fi

        echo ""
        echo -e "${BOLD}Where should we install the binary?${NC}"
        INSTALL_DIR=$(gum_choose \
            "Current Directory/FSB (Recommended)" \
            "Current Directory ($PWD)" \
            "$HOME/.local/bin" \
            "/usr/local/bin (Requires sudo/root)" \
            "Custom Path" | sed 's/ (.*//')

        if [[ "$INSTALL_DIR" == "Custom Path" ]]; then
            INSTALL_DIR=$(gum_input --placeholder "/path/to/install/dir")
        elif [[ "$INSTALL_DIR" == "Current Directory/FSB" ]]; then
            INSTALL_DIR="$PWD/FSB"
        elif [[ "$INSTALL_DIR" == "Current Directory" ]]; then
            INSTALL_DIR="$PWD"
        fi

        echo ""
        echo -e "${BOLD}Do you want to download a sample 'fsb.env' configuration file?${NC}"
        if gum_confirm "Download sample env?" --default=true; then
            WANT_ENV="yes"
        fi

        if [[ "$WANT_ENV" == "yes" ]]; then
            if [[ "$INSTALL_DIR" == "$PWD/FSB" || "$INSTALL_DIR" == "$PWD" ]]; then
                ENV_DIR="$INSTALL_DIR"
            else
                echo ""
                echo -e "${BOLD}Where should we save the 'fsb.env' file?${NC}"
                ENV_DIR=$(gum_choose \
                    "Current Directory ($PWD) (Recommended)" \
                    "Installation Folder ($INSTALL_DIR)" \
                    "Custom Path" | sed 's/ (.*//')

                if [[ "$ENV_DIR" == "Custom Path" ]]; then
                    ENV_DIR=$(gum_input --placeholder "/path/to/save/env")
                elif [[ "$ENV_DIR" == "Current Directory" ]]; then
                    ENV_DIR="$PWD"
                elif [[ "$ENV_DIR" == "Installation Folder" ]]; then
                    ENV_DIR="$INSTALL_DIR"
                fi
            fi
        fi
    else
        # Fallback text mode prompts
        prompt_read response "Install pre-release version? [y/N] " "N"
        case "$response" in
            [yY][eE][sS]|[yY]) WANT_PRERELEASE="yes" ;;
            *) WANT_PRERELEASE="no" ;;
        esac

        echo ""
        echo "Select installation directory:"
        echo "1) Current Directory/FSB (Recommended)"
        echo "2) Current Directory ($PWD)"
        echo "3) $HOME/.local/bin"
        echo "4) /usr/local/bin (Requires sudo)"
        prompt_read dir_resp "Selection [1]: " "1"
        case "${dir_resp:-1}" in
            2) INSTALL_DIR="$PWD" ;;
            3) INSTALL_DIR="$HOME/.local/bin" ;;
            4) INSTALL_DIR="/usr/local/bin" ;;
            *) INSTALL_DIR="$PWD/FSB" ;;
        esac

        echo ""
        prompt_read env_response "Download sample 'fsb.env' configuration file? [Y/n] " "Y"
        case "$env_response" in
            [nN][oO]|[nN]) WANT_ENV="no" ;;
            *) WANT_ENV="yes" ;;
        esac

        if [[ "$WANT_ENV" == "yes" ]]; then
            if [[ "$INSTALL_DIR" == "$PWD/FSB" || "$INSTALL_DIR" == "$PWD" ]]; then
                ENV_DIR="$INSTALL_DIR"
            else
                echo ""
                echo "Select where to save 'fsb.env':"
                echo "1) Current Directory ($PWD) (Recommended)"
                echo "2) Installation Folder ($INSTALL_DIR)"
                prompt_read env_dir_resp "Selection [1]: " "1"
                case "${env_dir_resp:-1}" in
                    2) ENV_DIR="$INSTALL_DIR" ;;
                    *) ENV_DIR="$PWD" ;;
                esac
            fi
        fi
    fi

        # Prep path
        case "$INSTALL_DIR" in
            ~*) INSTALL_DIR="${INSTALL_DIR/#\~/$HOME}" ;;
        esac
        if [[ "$WANT_ENV" == "yes" ]]; then
            case "$ENV_DIR" in
                ~*) ENV_DIR="${ENV_DIR/#\~/$HOME}" ;;
            esac
        fi

        TARGET_VER=$(get_latest_version "$WANT_PRERELEASE")
    fi
    
    local file_ext=".tar.gz"
    if [[ "$OS" == "windows" ]]; then
        file_ext=".zip"
    fi
    
    ASSET_NAME="TG-FileStreamBot-${TARGET_VER}-${OS}-${ARCH}${file_ext}"
    ASSET_URL="https://github.com/${REPO}/releases/download/${TARGET_VER}/${ASSET_NAME}"
    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${TARGET_VER}/TG-FileStreamBot-${TARGET_VER}-checksums.txt"
    SIG_URL="${CHECKSUM_URL}.sig"

    echo ""
    ui_info "Target Version    : ${BOLD}${TARGET_VER}${NC}"
    ui_info "OS / Arch         : ${BOLD}${OS} / ${ARCH}${NC}"
    ui_info "Binary Location   : ${BOLD}${INSTALL_DIR}${NC}"
    if [[ "$WANT_ENV" == "yes" ]]; then
        ui_info "Env File Location : ${BOLD}${ENV_DIR}/fsb.env${NC}"
    fi
    echo ""

    # Confirmation
    if can_use_tui; then
        if ! gum_confirm "Proceed with installation?" --default=true; then
            ui_info "Installation cancelled."
            exit 0
        fi
    else
        prompt_read proceed_resp "Proceed with installation? [Y/n] " "Y"
        case "$proceed_resp" in
            [nN][oO]|[nN])
                ui_info "Installation cancelled."
                exit 0
                ;;
        esac
    fi

    if [[ -n "$GUM" ]]; then
        "$GUM" spin --spinner dot --title "Downloading components..." -- sleep 0.1
    fi

    local dl_dir
    dl_dir="$(mktempdir)"

    ui_info "Downloading ${ASSET_NAME}..."
    if ! download_file "$ASSET_URL" "$dl_dir/$ASSET_NAME"; then
        ui_error "Failed to download the binary archive. Is the release available for this architecture?"
        exit 1
    fi

    ui_info "Downloading checksums..."
    download_file "$CHECKSUM_URL" "$dl_dir/checksums.txt"

    # Verify checksum if sha256sum or shasum exists
    ui_info "Verifying SHA256 checksum..."
    pushd "$dl_dir" >/dev/null
    local checksum_matched=false
    if command -v sha256sum &> /dev/null; then
        if grep "$ASSET_NAME" checksums.txt | sha256sum -c - >/dev/null 2>&1; then
            checksum_matched=true
        fi
    elif command -v shasum &> /dev/null; then
        if grep "$ASSET_NAME" checksums.txt | shasum -a 256 -c - >/dev/null 2>&1; then
            checksum_matched=true
        fi
    else
        ui_warn "No sha256sum or shasum found; skipping checksum verification."
        checksum_matched=true
    fi
    popd >/dev/null

    if [[ "$checksum_matched" != "true" ]]; then
        ui_error "Checksum verification failed! The file may be corrupted or compromised."
        exit 1
    fi
    ui_success "Checksum matched."

    # Verify Signature if tools exist
    ui_info "Downloading signature..."
    download_file "$SIG_URL" "$dl_dir/checksums.txt.sig"

    ui_info "Verifying signature..."
    pushd "$dl_dir" >/dev/null
    local sig_verified=false

    if command -v gpg &> /dev/null; then
        ui_info "gpg detected. Fetching developer's public key..."
        download_file "https://github.com/EverythingSuckz.gpg" "dev_key.asc"
        
        # Create a temporary keyring to avoid polluting the user's keychain
        local keyring="--no-default-keyring --keyring ./temp_keyring.gpg"
        gpg $keyring --import dev_key.asc >/dev/null 2>&1
        
        if gpg $keyring --verify checksums.txt.sig checksums.txt >/dev/null 2>&1; then
           sig_verified=true
           ui_success "Signature verified successfully."
        else
           ui_warn "Signature verification failed. The release might be compromised."
           sig_verified=false
        fi
    else
        ui_warn "No gpg executable found; skipping signature validation."
        sig_verified=true
    fi
    popd >/dev/null

    if [[ "$sig_verified" != "true" ]]; then
        ui_error "Signature verification failed!"
        exit 1
    fi

    ui_info "Extracting binary..."
    if [[ "$ASSET_NAME" == *.zip ]]; then
        if ! command -v unzip &> /dev/null; then
            ui_error "unzip is required to extract Windows binaries."
            exit 1
        fi
        unzip -q "$dl_dir/$ASSET_NAME" -d "$dl_dir"
    else
        tar -xzf "$dl_dir/$ASSET_NAME" -C "$dl_dir"
    fi
    
    # The tarball contains the binary, possibly named `fsb` or inside a folder. 
    # Usually it extracts right out as `fsb`.
    local bin_file
    bin_file="$(find "$dl_dir" -type f \( -name "fsb" -o -name "TG-FileStreamBot" \) | head -n1)"

    if [[ -z "$bin_file" ]]; then
        ui_error "Could not find the executable inside the extracted archive."
        exit 1
    fi

    chmod +x "$bin_file"

    if [[ ! -d "$INSTALL_DIR" ]]; then
        ui_info "Creating directory $INSTALL_DIR..."
        mkdir -p "$INSTALL_DIR" 2>/dev/null || sudo mkdir -p "$INSTALL_DIR"
    fi

    ui_info "Installing to $INSTALL_DIR..."
    if mv "$bin_file" "$INSTALL_DIR/fsb" 2>/dev/null; then
        ui_success "Installed successfully!"
    else
        ui_info "Elevated permissions required to write to $INSTALL_DIR"
        # Use sudo on unix, but check if we are on windows first
        if command -v sudo &> /dev/null; then
            if sudo mv "$bin_file" "$INSTALL_DIR/fsb"; then
                ui_success "Installed successfully!"
            else
                ui_error "Installation failed due to permissions."
                exit 1
            fi
        else
            ui_error "Insufficient permissions and sudo not found. Please run script as Admin/Root."
            exit 1
        fi
    fi

    echo ""
    echo -e "${FSB_GREEN}${BOLD}✓ Setup Complete!${NC}"
    
    # Path logic checking
    local in_path=false
    if builtin type -p fsb &>/dev/null; then
        in_path=true
    elif [[ ":$PATH:" == *":$INSTALL_DIR:"* ]]; then
        in_path=true
    fi

    local is_portable=false
    if [[ "$INSTALL_DIR" != "/usr/local/bin" && "$INSTALL_DIR" != "/usr/bin" && "$INSTALL_DIR" != "$HOME/.local/bin" ]]; then
        is_portable=true
    fi

    echo ""
    if [[ "$in_path" == false && "$is_portable" == false ]]; then
        ui_warn "The directory $INSTALL_DIR is not in your PATH."
        echo -e "You may need to add it to your shell configuration file (e.g. ~/.bashrc or ~/.zshrc):"
        echo -e "  ${BOLD}export PATH=\"\$PATH:$INSTALL_DIR\"${NC}"
        echo ""
    fi
    
    local bin_cmd="fsb"
    if [[ "$in_path" == false ]]; then
        bin_cmd="./fsb"
    fi

    # Download Env if requested
    if [[ "$WANT_ENV" == "yes" ]]; then
        ui_info "Downloading sample configuration to $ENV_DIR/fsb.env..."
        if [[ ! -d "$ENV_DIR" ]]; then
            mkdir -p "$ENV_DIR" 2>/dev/null || sudo mkdir -p "$ENV_DIR"
        fi
        
        local env_branch="main"
        if [[ "$WANT_PRERELEASE" == "yes" ]]; then
            env_branch="dev"
        fi
        
        local env_url="https://raw.githubusercontent.com/${REPO}/refs/heads/${env_branch}/fsb.sample.env"
        local dl_env
        dl_env="$(mktempdir)/fsb.sample.env"
        
        if download_file "$env_url" "$dl_env"; then
            if mv "$dl_env" "$ENV_DIR/fsb.env" 2>/dev/null; then
                ui_success "Saved fsb.env"
            else
                sudo mv "$dl_env" "$ENV_DIR/fsb.env" 2>/dev/null || ui_warn "Failed to move fsb.env to $ENV_DIR"
            fi
        else
            ui_warn "Failed to download sample fsb.env configuration."
        fi
    fi

    echo ""
    echo -e "${BOLD}Next Steps:${NC}"
    if [[ "$INSTALL_DIR" != "$PWD" && "$is_portable" == true ]]; then
        # If it's a sub-directory, e.g. PWD/FSB or custom folder
        echo -e "  cd $INSTALL_DIR"
    fi

    if [[ "$WANT_ENV" == "yes" ]]; then
        echo -e "  ${INFO}# Edit fsb.env to configure your bot variables${NC}"
    else
        echo -e "  ${INFO}# Supply your environmental variables (e.g. in fsb.env) before running${NC}"
    fi
    echo -e "  ${bin_cmd} run"
    echo ""
}

main "$@"
