#!/usr/bin/env bash
# One-shot VPS bootstrap: Docker, firewall, clone repos (HTTPS), write .env with
# generated secrets, Postgres in Compose, migrations, then full stack.
#
# Usage (as root on Ubuntu 22.04+):
#   curl -fsSL .../bootstrap-pramool.sh | bash
# or:
#   bash bootstrap-pramool.sh
#
set -euo pipefail

# ===================== tweak if needed =====================
INSTALL_ROOT="${INSTALL_ROOT:-/opt/pramool}"
GIT_CORE_URL="${GIT_CORE_URL:-https://github.com/rnikrozoft/pramool-core.git}"
GIT_WALLET_URL="${GIT_WALLET_URL:-https://github.com/rnikrozoft/pramool-wallet-service.git}"
GIT_AUCTION_URL="${GIT_AUCTION_URL:-https://github.com/rnikrozoft/pramool-auction-service.git}"
GIT_FRONTEND_URL="${GIT_FRONTEND_URL:-https://github.com/rnikrozoft/pramool.in.th.git}"
CORE_BRANCH="${CORE_BRANCH:-develop}"
WALLET_BRANCH="${WALLET_BRANCH:-main}"
AUCTION_BRANCH="${AUCTION_BRANCH:-main}"
FRONTEND_BRANCH="${FRONTEND_BRANCH:-main}"
# ===========================================================

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing command: $1" >&2
    exit 1
  }
}

rand_hex() { openssl rand -hex "$1"; }

echo "[1/10] Install base packages"
export DEBIAN_FRONTEND=noninteractive
require_cmd openssl
apt-get update -y
apt-get install -y ca-certificates curl git ufw

echo "[2/10] Install Docker + Compose plugin"
if ! command -v docker >/dev/null 2>&1; then
  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
  chmod a+r /etc/apt/keyrings/docker.asc
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" >/etc/apt/sources.list.d/docker.list
  apt-get update -y
  apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
  systemctl enable docker
  systemctl restart docker
fi

echo "[3/10] Firewall"
ufw allow OpenSSH
ufw allow 80/tcp || true
ufw allow 443/tcp || true
ufw allow 3001/tcp
ufw allow 3102/tcp
ufw allow 3103/tcp
ufw allow 3104/tcp
ufw allow 3000/tcp
ufw --force enable || true

echo "[4/10] Clone / refresh repositories"
mkdir -p "$INSTALL_ROOT"
cd "$INSTALL_ROOT"
rm -rf pramool-core pramool-wallet-service pramool-auction-service pramool.in.th || true
git clone -b "$CORE_BRANCH" "$GIT_CORE_URL" pramool-core
git clone -b "$WALLET_BRANCH" "$GIT_WALLET_URL" pramool-wallet-service
git clone -b "$AUCTION_BRANCH" "$GIT_AUCTION_URL" pramool-auction-service
git clone -b "$FRONTEND_BRANCH" "$GIT_FRONTEND_URL" pramool.in.th

DEPLOY_DIR="$INSTALL_ROOT/pramool-core/deploy/vps"
test -f "$DEPLOY_DIR/docker-compose.yml" || {
  echo "missing $DEPLOY_DIR/docker-compose.yml (wrong branch?)" >&2
  exit 1
}

echo "[5/10] Generate secrets (no manual .env editing)"
DB_PASS="$(rand_hex 16)"
JWT_SECRET="$(rand_hex 32)"
INTERNAL_KEY="$(rand_hex 24)"

PUB_IP=""
command -v curl >/dev/null 2>&1 && PUB_IP="$(curl -fsS --max-time 5 https://api.ipify.org 2>/dev/null || true)"
if [[ -z "${PUB_IP:-}" ]]; then
  PUB_IP="$(hostname -I 2>/dev/null | awk '{print $1}')"
fi
[[ -z "${PUB_IP:-}" ]] && PUB_IP="127.0.0.1"

CORS_ALLOW_ORIGINS="http://localhost:3000,http://127.0.0.1:3000"
if [[ "${PUB_IP}" != "127.0.0.1" ]]; then
  CORS_ALLOW_ORIGINS+=",http://${PUB_IP}:3000,http://${PUB_IP}:3001"
fi

BASE_HTTP="http://${PUB_IP}"
NEXT_PUBLIC_CORE_API_BASE_URL="${BASE_HTTP}:3001"
NEXT_PUBLIC_USER_API_BASE_URL="${BASE_HTTP}:3001"
NEXT_PUBLIC_WALLET_API_BASE_URL="${BASE_HTTP}:3102"
NEXT_PUBLIC_AUCTION_REALTIME_BASE_URL="${BASE_HTTP}:3103"

# Password must be URL-safe: pramool-core builds DSN without query escaping.
ENC_PASS="$DB_PASS"

cat >"$DEPLOY_DIR/.env" <<EOF
CORS_ALLOW_ORIGINS=${CORS_ALLOW_ORIGINS}
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRE_TIME=1
JWT_REFRESH_EXPIRE_TIME=168

ADDRESS_REQUEST=
ADDRESS_VERIFY=
API_KEY=
API_SECRET=

DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=${DB_PASS}
DATABASE_NAME=pramool

DATABASE_DSN=postgres://postgres:${ENC_PASS}@postgres:5432/pramool?sslmode=disable
AUCTION_DATABASE_URL=postgres://postgres:${ENC_PASS}@postgres:5432/pramool?sslmode=disable

WALLET_INTERNAL_KEY=${INTERNAL_KEY}
AUCTION_INTERNAL_KEY=${INTERNAL_KEY}
ESCROW_AUTO_CONFIRM_DAYS=14

OMISE_SECRET_KEY=
OMISE_WEBHOOK_SECRET=
OMISE_SYSTEM_RECIPIENT_ID=

PGRST_DB_URI=postgres://postgres:${ENC_PASS}@postgres:5432/pramool?sslmode=disable
PGRST_DB_SCHEMA=public
PGRST_DB_ANON_ROLE=postgres
PGRST_OPENAPI_MODE=follow-privileges
PGRST_DB_POOL=10

NEXT_PUBLIC_CORE_API_BASE_URL=${NEXT_PUBLIC_CORE_API_BASE_URL}
NEXT_PUBLIC_USER_API_BASE_URL=${NEXT_PUBLIC_USER_API_BASE_URL}
NEXT_PUBLIC_WALLET_API_BASE_URL=${NEXT_PUBLIC_WALLET_API_BASE_URL}
NEXT_PUBLIC_AUCTION_REALTIME_BASE_URL=${NEXT_PUBLIC_AUCTION_REALTIME_BASE_URL}

FRONTEND_PORT=3000
CORE_PORT=3001
WALLET_PORT=3102
AUCTION_PORT=3103
POSTGREST_PORT=3104
EOF

# Mount is visible inside the container as uid appuser; root-only mode (0600) causes "permission denied".
chmod 0644 "$DEPLOY_DIR/.env"

echo "[6/10] Build images"
cd "$DEPLOY_DIR"
docker compose build

echo "[7/10] Start Postgres and wait"
docker compose up -d postgres
for i in $(seq 1 60); do
  if docker compose exec -T postgres pg_isready -U postgres -d pramool >/dev/null 2>&1; then
    break
  fi
  sleep 1
  if [[ "$i" -eq 60 ]]; then
    echo "Postgres did not become ready in time" >&2
    exit 1
  fi
done

echo "[8/10] Run pramool-core migrations"
docker compose run --rm pramool-core ./pramool-core migrate

echo "[9/10] Start full stack"
docker compose up -d

echo "[10/10] Done"
docker compose ps

echo ""
echo "Wrote secrets to: $DEPLOY_DIR/.env (mode 644 so non-root app user in containers can read it)"
echo "Database password (postgres user): $DB_PASS"
echo "Tail logs:  cd $DEPLOY_DIR && docker compose logs -f --tail=120"
echo "Frontend:   ${BASE_HTTP}:3000"
echo "Core API:   ${BASE_HTTP}:3001"
echo ""
echo "Note: Omise / ThaiBulkSMS are empty until you add keys to .env. PostgREST uses role 'postgres' for lab only."
echo "If the site cannot reach APIs, fix NEXT_PUBLIC_* in .env and rebuild: docker compose build pramool-frontend --no-cache && docker compose up -d"
