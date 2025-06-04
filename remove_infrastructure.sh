#!/usr/bin/env bash
# Remove Docker & Node build scripts and folders no longer used by pure LogLineOS

set -euo pipefail

echo "🗑️  Removendo Docker e infrações legadas..."

# Delete Docker config
rm -f Dockerfile docker-compose.yml
rm -rf docker/ docker/env-config.sh

# Remove Nginx config
rm -f docker/nginx.conf

echo "🛠️  Removendo scripts Docker de package.json..."
# Clean up package.json scripts
jq 'del(.scripts["docker:build", "docker:run", "docker:compose"])' package.json > package.tmp.json
mv package.tmp.json package.json

echo "✅ Infraestrutura legada removida com sucesso."