#!/usr/bin/env bash
set -euo pipefail

# Sibling repos: BASE_DIR/pramool-core, BASE_DIR/pramool-wallet-service, ...
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# One entry per line: "dirname" (go run .) or "dirname subcommand" (go run . subcommand)
SERVICES=(
  "pramool-core serve"
  "pramool-wallet-service"
  "pramool-auction-service"
)

PIDS=()

cleanup() {
  for pid in "${PIDS[@]:-}"; do
    if kill -0 "$pid" >/dev/null 2>&1; then
      kill "$pid" >/dev/null 2>&1 || true
    fi
  done
}

trap cleanup INT TERM EXIT

for spec in "${SERVICES[@]}"; do
  read -r service_dir subcmd <<<"$spec"
  service_path="$BASE_DIR/$service_dir"
  if [[ ! -d "$service_path" ]]; then
    echo "missing service directory: $service_path" >&2
    exit 1
  fi

  (
    cd "$service_path"
    # Seller auction uploads must land where core still serves /uploads (same paths in DB) until cloud storage.
    if [[ "$service_dir" == "pramool-auction-service" ]]; then
      export PRAMOOL_UPLOAD_ROOT="$BASE_DIR/pramool-core"
    fi
    if [[ -f ".env" ]]; then
      set -a
      # shellcheck disable=SC1091
      source ".env"
      set +a
    fi
    if [[ -n "${subcmd:-}" ]]; then
      echo "starting $service_dir (go run . $subcmd)"
      go run . "$subcmd"
    else
      echo "starting $service_dir (go run .)"
      go run .
    fi
  ) &
  PIDS+=("$!")
done

echo "all services started (core :3001, wallet :3102, auction :3103). Ctrl+C to stop."
wait
