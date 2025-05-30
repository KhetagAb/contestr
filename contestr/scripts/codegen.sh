#!/bin/bash
set -e

echo "Using oapi-codegen version:"
oapi-codegen --version || echo "unknown"

mkdir -p internal/generated

echo "Generating code from OpenAPI specification..."
oapi-codegen -config api/config.yaml api/openapi.yaml
echo "Code generation completed successfully!"
