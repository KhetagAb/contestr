#!/bin/bash
set -e

# Цвета
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # нет цвета

print_status() {
  local msg=$1
  local status=$2
  local color=$3
  printf "%-30s: %b%s%b\n" "$msg" "$color" "$status" "$NC"
}

echo -n "oapi-codegen version: "
version=$(oapi-codegen --version 2>/dev/null || echo "unknown")
echo "$version"

if oapi-codegen -config api/config.yaml api/openapi.yaml; then
  print_status "Generating OpenAPI code" "OK" "$GREEN"
else
  print_status "Generating OpenAPI code" "FAILED" "$RED"
  exit 1
fi

if (cd internal/generated/app && wire); then
  print_status "Generating Wire code" "OK" "$GREEN"
else
  print_status "Generating Wire code" "FAILED" "$RED"
  exit 1
fi

print_status "Code generation completed" "OK" "$GREEN"
