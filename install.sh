#!/bin/bash
set -e

OWNER="zahidcakici"
REPO="tokenizer-cli"
BINARY="tokenizer"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     OS='linux';;
    Darwin*)    OS='darwin';;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Detect Arch
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)    ARCH='amd64';;
    arm64)     ARCH='arm64';;
    aarch64)   ARCH='arm64';;
    *)         echo "Unsupported Architecture: ${ARCH}"; exit 1;;
esac

echo "Detected platform: ${OS}/${ARCH}"

# Get latest version
echo "Fetching latest version..."
LATEST_URL="https://api.github.com/repos/${OWNER}/${REPO}/releases/latest"
VERSION_JSON=$(curl -s $LATEST_URL)

# Check if we got a valid response (simple check)
if [[ "$VERSION_JSON" == *"message"* && "$VERSION_JSON" == *"Not Found"* ]]; then
    echo "Error: Could not find latest release. Please check if any releases adhere to the naming convention."
    exit 1
fi

TAG_NAME=$(echo "$VERSION_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG_NAME" ]; then
    echo "Error: Failed to fetch latest version tag."
    exit 1
fi

echo "Latest version: ${TAG_NAME}"

# Strip 'v' prefix for the filename construction if present, as GoReleaser usually drops it for artifacts
CLEAN_VERSION="${TAG_NAME#v}"

# Construct filename matching .goreleaser.yaml: {{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}
FILENAME="${REPO}_${CLEAN_VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${TAG_NAME}/${FILENAME}"

echo "Downloading ${FILENAME}..."
TMP_DIR=$(mktemp -d)
curl -sL "$DOWNLOAD_URL" -o "${TMP_DIR}/${FILENAME}"

if [ ! -f "${TMP_DIR}/${FILENAME}" ]; then 
    echo "Error: Download failed."
    exit 1
fi

echo "Extracting..."
tar -xzf "${TMP_DIR}/${FILENAME}" -C "$TMP_DIR"

echo "Installing ${BINARY} to /usr/local/bin requires sudo..."
if sudo mv "${TMP_DIR}/${BINARY}" /usr/local/bin/; then
    echo "Successfully installed ${BINARY} to /usr/local/bin/"
else
    echo "Error: Failed to move binary to /usr/local/bin/"
    rm -rf "$TMP_DIR"
    exit 1
fi

echo "Cleaning up..."
rm -rf "$TMP_DIR"

echo "Done! You can now run '${BINARY}'."
