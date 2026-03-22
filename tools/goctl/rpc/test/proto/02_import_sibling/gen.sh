#!/usr/bin/env bash
# Scenario 02: sibling import
# Usage: bash gen.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GOCTL_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
GOCTL="$GOCTL_ROOT/bin/goctl"
OUT="$SCRIPT_DIR/output"

# Build goctl from source
go build -o "$GOCTL" "$GOCTL_ROOT"

# Clean and initialize output directory
rm -rf "$OUT" && mkdir -p "$OUT/pb"
(cd "$OUT" && go mod init example.com/demo/s02_import_sibling > /dev/null 2>&1)

# Generate code
cd "$SCRIPT_DIR"
"$GOCTL" rpc protoc user.proto \
  --go_out="$OUT/pb" \
  --go-grpc_out="$OUT/pb" \
  --zrpc_out="$OUT/rpc" \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --proto_path=.

# Verify build
echo "Running go mod tidy..."
cd "$OUT" && go mod tidy
echo "Checking build..."
if go build ./...; then
  echo "✅ Build passed"
else
  echo "❌ Build failed"
  exit 1
fi

echo "Done. Output: $OUT"
