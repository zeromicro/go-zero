#!/usr/bin/env bash
# Scenario 09: compare old vs new goctl output — google well-known types as RPC request/response
# Usage: bash compare.sh
# Requires: go install github.com/zeromicro/go-zero/tools/goctl@latest (auto-installed)
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GOCTL_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
NEW_GOCTL="$GOCTL_ROOT/bin/goctl"
OLD_GOCTL="$(go env GOPATH)/bin/goctl"
OUT_OLD="$SCRIPT_DIR/output_old"
OUT_NEW="$SCRIPT_DIR/output_new"

verify_build() {
  local dir="$1" label="$2"
  echo "Verifying $label ..."
  cd "$dir"
  go mod tidy
  if go build ./...; then
    echo "  ✅ $label: build passed"
  else
    echo "  ❌ $label: build failed"
    exit 1
  fi
  cd "$SCRIPT_DIR"
}

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

# Install released goctl and build local goctl
echo ">>> Installing goctl@latest ..."
go install github.com/zeromicro/go-zero/tools/goctl@latest
echo ">>> Building local goctl ..."
go build -o "$NEW_GOCTL" "$GOCTL_ROOT"

# Generate with old goctl
echo ">>> Generating with old goctl ..."
rm -rf "$OUT_OLD" && mkdir -p "$OUT_OLD/pb"
(cd "$OUT_OLD" && go mod init example.com/demo/s09_google_types > /dev/null 2>&1)
cd "$SCRIPT_DIR"
set +e
"$OLD_GOCTL" rpc protoc service.proto \
  --go_out="$OUT_OLD/pb" \
  --go-grpc_out="$OUT_OLD/pb" \
  --zrpc_out="$OUT_OLD/rpc" \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  --proto_path="$PROTOC_INCLUDE"
GEN_STATUS=$?
set -e
if [ "$GEN_STATUS" -ne 0 ]; then
  echo "  ⚠️  old goctl does not support this feature (expected)"
else
  verify_build "$OUT_OLD" "old"
fi

# Generate with new goctl
echo ">>> Generating with new goctl ..."
rm -rf "$OUT_NEW" && mkdir -p "$OUT_NEW/pb"
(cd "$OUT_NEW" && go mod init example.com/demo/s09_google_types > /dev/null 2>&1)
cd "$SCRIPT_DIR"
"$NEW_GOCTL" rpc protoc service.proto \
  --go_out="$OUT_NEW/pb" \
  --go-grpc_out="$OUT_NEW/pb" \
  --zrpc_out="$OUT_NEW/rpc" \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --proto_path=. \
  --proto_path="$PROTOC_INCLUDE"
verify_build "$OUT_NEW" "new"

# Diff old vs new (exclude go.mod / go.sum)
echo ""
echo ">>> Diff (old vs new):"
if diff -rq --exclude="go.mod" --exclude="go.sum" "$OUT_OLD" "$OUT_NEW" > /dev/null 2>&1; then
  echo "  [identical] no differences between old and new output"
else
  diff -r --exclude="go.mod" --exclude="go.sum" "$OUT_OLD" "$OUT_NEW" || true
fi
