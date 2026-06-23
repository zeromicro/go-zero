#!/usr/bin/env bash
# Scenario 09: google well-known types as RPC request/response
# Usage: bash gen.sh
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GOCTL_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
GOCTL="$GOCTL_ROOT/bin/goctl"
OUT="$SCRIPT_DIR/output"

# Find well-known types include path
PROTOC_INCLUDE=""
PROTOC_BIN="$(which protoc 2>/dev/null || true)"
if [ -n "$PROTOC_BIN" ]; then
  CANDIDATE="$(cd "$(dirname "$PROTOC_BIN")/.." && pwd)/include"
  [ -f "$CANDIDATE/google/protobuf/empty.proto" ] && PROTOC_INCLUDE="$CANDIDATE"
fi
if [ -z "$PROTOC_INCLUDE" ]; then
  for d in /opt/homebrew/include /usr/local/include; do
    [ -f "$d/google/protobuf/empty.proto" ] && PROTOC_INCLUDE="$d" && break
  done
fi
if [ -z "$PROTOC_INCLUDE" ]; then
  echo "Error: cannot find google/protobuf/empty.proto. Install protobuf (brew install protobuf)."
  exit 1
fi
echo "Well-known types: $PROTOC_INCLUDE"

# Build goctl from source
go build -o "$GOCTL" "$GOCTL_ROOT"

# Clean and initialize output directory
rm -rf "$OUT" && mkdir -p "$OUT/pb"
(cd "$OUT" && go mod init example.com/demo/s09_google_types > /dev/null 2>&1)

# Generate code
cd "$SCRIPT_DIR"
"$GOCTL" rpc protoc service.proto \
  --go_out="$OUT/pb" \
  --go-grpc_out="$OUT/pb" \
  --zrpc_out="$OUT/rpc" \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  --proto_path="$PROTOC_INCLUDE"

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
