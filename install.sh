#!/bin/bash

# --- Конфигурация ---
REPO="algrvvv/goproxy"
RAW_CONFIG_PATH="config.example.yaml"
RAW_CONFIG_FILE="config.yaml"
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
    Linux)
        BIN_FILE="linux_goproxy"
        ;;
    Darwin)
      echo "To install on your OS, use the first installation option or clone the repository yourself and run the build"
      BIN_FILE="err"
        ;;
    CYGWIN*|MINGW32*|MSYS*|MINGW*)
        BIN_FILE="windows_goproxy.exe"
        ;;
    *)
        echo "Error: unsupported os - $OS"
        exit 1
        ;;
esac

# --- Выбор установки ---
echo "Select installation option :"
echo "1) Clone repository"
echo "2) Install the binary from the latest release"
read -p "Enter number (1 or 2): " choice

if [ "$choice" == "1" ]; then
    # --- Клонируем ---
  git clone https://github.com/algrvvv/goproxy
  cd goproxy

  # --- Делаем билд ---
  # LOWER_OS=$(echo "$OS" | awk '{print tolower($0)}')
  # GOOS=$LOWER_OS GOARCH=$ARCH
  go build -o bin/goproxy cmd/goproxy/main.go

  START_LINE="For start proxy server use: cd goproxy && ./bin/goproxy"
elif [ "$choice" == "2" ]; then
  echo "Getting the latest release version $REPO... for your os: $OS ($ARCH)"
  if [ "$BIN_FILE" == "err" ]; then
    echo "Error: This installation option unsupported for your os, select 1 installation option"
    exit 1
  fi

  LATEST_RELEASE_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep "browser_download_url.*$BIN_FILE" | cut -d '"' -f 4)

  if [ -z "$LATEST_RELEASE_URL"]; then
    echo "Error: Could not get latest release URL"
    exit 1
  fi

  # --- Установка ---
  echo "Downloading the latest release..."
  curl -Lo goproxy $LATEST_RELEASE_URL

  if [ "$OS" != "CYGWIN" ] && [ "$OS" != "MINGW32" ] && [ "$OS" != "MSYS" ] && [ "$OS" != "MINGW" ]; then
    chmod +x goproxy
  fi

  START_LINE="For start proxy server use: ./goproxy"
else
  echo "Error: wrong installation option"
  exit 1
fi

echo "Downloading config example..."
RAW_URL="https://raw.githubusercontent.com/$REPO/main/$RAW_CONFIG_PATH"
echo $RAW_URL
curl -o $RAW_CONFIG_FILE $RAW_URL

# --- Завершение ---
echo "Installation is complete, customize config.yaml"
echo "$START_LINE"
