#!/usr/bin/env bash
# Scenario 06: well-known type imports
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
  [ -f "$CANDIDATE/google/protobuf/timestamp.proto" ] && PROTOC_INCLUDE="$CANDIDATE"
fi
if [ -z "$PROTOC_INCLUDE" ]; then
  for d in /opt/homebrew/include /usr/local/include; do
    [ -f "$d/google/protobuf/timestamp.proto" ] && PROTOC_INCLUDE="$d" && break
  done
fi
if [ -z "$PROTOC_INCLUDE" ]; then
  echo "Error: cannot find google/protobuf/timestamp.proto. Install protobuf (brew install protobuf)."
  exit 1
fi
echo "Well-known types: $PROTOC_INCLUDE"

# Build goctl from source
go build -o "$GOCTL" "$GOCTL_ROOT"

# Clean and initialize output directory
rm -rf "$OUT" && mkdir -p "$OUT/pb"
(cd "$OUT" && go mod init example.com/demo/s06_wellknown > /dev/null 2>&1)

# Generate code
cd "$SCRIPT_DIR"
"$GOCTL" rpc protoc events.proto \
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
