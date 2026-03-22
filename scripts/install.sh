#!/bin/bash

# homectl installation script
# https://homectl.xyz

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================================${NC}"
echo -e "${BLUE}               homectl installation                   ${NC}"
echo -e "${BLUE}======================================================${NC}"

# Check dependencies
check_dep() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo -e "${RED}Error: $1 is not installed.${NC}"
        return 1
    fi
    return 0
}

echo "Checking dependencies..."
check_dep "docker" || exit 1
check_dep "docker-compose" || check_dep "docker compose" || exit 1

# Create directory
INSTALL_DIR="homectl"
if [ -d "$INSTALL_DIR" ]; then
    echo -e "${RED}Error: Directory '$INSTALL_DIR' already exists.${NC}"
    exit 1
fi

mkdir -p "$INSTALL_DIR/data/db" "$INSTALL_DIR/data/icons"
cd "$INSTALL_DIR"

# Download docker-compose.yml
echo "Downloading docker-compose.yml..."
curl -sSL https://raw.githubusercontent.com/palta-dev/homectl/main/docker-compose.yml -o docker-compose.yml

# Create initial config if missing
if [ ! -f "data/config.yaml" ]; then
    echo "Creating default configuration..."
    cat > data/config.yaml <<EOF
version: 1
settings:
  title: "homectl"
  theme: "dark"
  requestTimeout: "10s"
groups:
  - name: "Infrastructure"
    services:
      - name: "homectl GitHub"
        url: "https://github.com/palta-dev/homectl"
        description: "Documentation and source code"
        icon: "github"
EOF
fi

# Pull and start
echo -e "${GREEN}Starting homectl...${NC}"
if command -v "docker-compose" >/dev/null 2>&1; then
    docker-compose pull
    docker-compose up -d
else
    docker compose pull
    docker compose up -d
fi

echo -e "${GREEN}======================================================${NC}"
echo -e "${GREEN}  Installation complete!                              ${NC}"
echo -e "${GREEN}  Dashboard: http://localhost:7777                    ${NC}"
echo -e "${GREEN}======================================================${NC}"
echo -e "  Config file: $INSTALL_DIR/data/config.yaml"
echo -e "  To stop: cd $INSTALL_DIR && docker-compose down"
echo -e "${GREEN}======================================================${NC}"
