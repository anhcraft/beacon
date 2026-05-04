#!/usr/bin/env bash

set -e

echo "Starting Beacon installation..."

# 1. Require root privileges for /etc and /usr/local/bin
if [[ $EUID -ne 0 ]]; then
   echo "Error: This script must be run as root. Try 'sudo bash $0'"
   exit 1
fi

# 2. Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"
DOWNLOAD_URL=""

if [[ "$OS" == "Linux" && "$ARCH" == "x86_64" ]]; then
    echo "Detected Linux amd64..."
    DOWNLOAD_URL="https://github.com/anhcraft/beacon/releases/latest/download/beacon-linux-amd64"
elif [[ "$OS" == "Darwin" && "$ARCH" == "arm64" ]]; then
    echo "Detected macOS (Darwin) arm64..."
    DOWNLOAD_URL="https://github.com/anhcraft/beacon/releases/latest/download/beacon-darwin-arm64"
else
    echo "Error: Unsupported OS or Architecture ($OS $ARCH)."
    echo "This script only supports Linux amd64 and macOS arm64."
    exit 1
fi

# 3. Download and install the binary
BIN_PATH="/usr/local/bin/beacon"
echo "Downloading beacon binary..."
curl -L -o "$BIN_PATH" "$DOWNLOAD_URL"
chmod +x "$BIN_PATH"
echo "Binary installed to $BIN_PATH"

# 4. Set up configuration directories and files
CONFIG_DIR="/etc/beacon"
CONFIG_FILE="$CONFIG_DIR/config.yml"
GCP_CRED_FILE="$CONFIG_DIR/gcp_credentials.json"

echo "Setting up configuration at $CONFIG_DIR..."
mkdir -p "$CONFIG_DIR"

if [[ ! -f "$CONFIG_FILE" ]]; then
    echo "Creating default config.yml..."
    touch "$CONFIG_FILE"
else
    echo "Config file already exists at $CONFIG_FILE. Skipping."
fi

# 5. OS-Specific Daemon Setup (systemd for Linux)
if [[ "$OS" == "Linux" ]]; then
    SERVICE_FILE="/etc/systemd/system/beacon.service"
    echo "Configuring systemd service at $SERVICE_FILE..."

    # We use a shell wrapper in ExecStart to conditionally check for the GCP file at runtime
    cat <<EOF > "$SERVICE_FILE"
[Unit]
Description=Beacon Daemon
After=network.target

[Service]
Type=simple
# Check for GCP credentials before starting, export if found, then execute binary
ExecStart=/bin/sh -c 'if [ -f $GCP_CRED_FILE ]; then export GOOGLE_APPLICATION_CREDENTIALS=$GCP_CRED_FILE; fi; exec $BIN_PATH -config $CONFIG_FILE'
Restart=always
RestartSec=5

# Log output to standard system journal
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

    echo "Reloading systemd daemon..."
    systemctl daemon-reload
    systemctl enable beacon

    echo "--------------------------------------------------------"
    echo "Installation complete!"
    echo "1. Edit your config:   nano $CONFIG_FILE"
    echo "2. Add GCP keys (opt): nano $GCP_CRED_FILE"
    echo "3. Start the daemon:   systemctl start beacon"
    echo "4. View live logs:     journalctl -u beacon -f"
    echo "--------------------------------------------------------"

elif [[ "$OS" == "Darwin" ]]; then
    echo "--------------------------------------------------------"
    echo "Installation complete!"
    echo "Binary installed to:   $BIN_PATH"
    echo "Config folder created: $CONFIG_DIR"
    echo ""
    echo "To run beacon with optional GCP support on macOS, use:"
    echo "  [ -f $GCP_CRED_FILE ] && export GOOGLE_APPLICATION_CREDENTIALS=$GCP_CRED_FILE"
    echo "  $BIN_PATH -config $CONFIG_FILE"
    echo "--------------------------------------------------------"
fi