#!/bin/bash

set -e

# Configuration
APP_NAME="db"
MAIN_PACKAGE="./cmd"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR="./build"
RELEASE_DIR="./release"
CONFIG_FILE="config.yaml"

# Color definitions
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[34m"
RESET="\033[0m"

# Print colored message
print_msg() {
  echo -e "${BLUE}$1${RESET}"
}

# Check if Go is installed
check_dependencies() {
  print_msg "Checking dependencies..."
  if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    exit 1
  fi
}

# Clean build artifacts
clean() {
  print_msg "Cleaning build artifacts..."
  rm -rf $BUILD_DIR $RELEASE_DIR
}

# Run tests
test() {
  print_msg "Running tests..."
  go test -v ./...
}

# Build for current platform
build() {
  print_msg "Building $APP_NAME v$VERSION for current platform..."
  mkdir -p $BUILD_DIR
  
  go build -ldflags "-X main.version=$VERSION" -o $BUILD_DIR/$APP_NAME $MAIN_PACKAGE
  cp $CONFIG_FILE $BUILD_DIR/ 2>/dev/null || touch $BUILD_DIR/$CONFIG_FILE
  
  print_msg "Build complete: $BUILD_DIR/$APP_NAME"
}

# Build for all platforms
release() {
  print_msg "Creating release builds for v$VERSION..."
  mkdir -p $RELEASE_DIR
  
  # Linux (amd64)
  GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o $RELEASE_DIR/${APP_NAME}-linux-amd64 $MAIN_PACKAGE
  
  # MacOS (amd64)
  GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o $RELEASE_DIR/${APP_NAME}-darwin-amd64 $MAIN_PACKAGE
  
  # MacOS (arm64)
  GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$VERSION" -o $RELEASE_DIR/${APP_NAME}-darwin-arm64 $MAIN_PACKAGE
  
  # Windows (amd64)
  GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o $RELEASE_DIR/${APP_NAME}-windows-amd64.exe $MAIN_PACKAGE
  
  print_msg "Release builds created in $RELEASE_DIR"
}

# Install dependencies
install_deps() {
  print_msg "Installing dependencies..."
  go mod tidy
}

# Show help
show_help() {
  echo "Build script for $APP_NAME"
  echo
  echo "Usage: $0 [command]"
  echo
  echo "Commands:"
  echo "  build       Build for the current platform"
  echo "  test        Run tests"
  echo "  release     Build releases for multiple platforms"
  echo "  clean       Remove build artifacts"
  echo "  deps        Install dependencies"
  echo "  all         Clean, install deps, test, and build"
  echo "  help        Show this help message"
  echo
}

# Main execution
check_dependencies

case "$1" in
  "build")
    build
    ;;
  "test")
    test
    ;;
  "release")
    release
    ;;
  "clean")
    clean
    ;;
  "deps")
    install_deps
    ;;
  "all")
    clean
    install_deps
    test
    build
    ;;
  "help"|"")
    show_help
    ;;
  *)
    echo "Unknown command: $1"
    show_help
    exit 1
    ;;
esac