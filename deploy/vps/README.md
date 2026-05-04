# VPS Deployment (Docker Compose)

This stack deploys:

- `postgres` (16, data in Docker volume `pgdata`)
- `pramool-core`
- `pramool-wallet-service`
- `pramool-auction-service`
- `pramool-frontend` (Next.js `pramool.in.th`, **built with** `NEXT_PUBLIC_*` API URLs)
- `postgrest`

## Option A — one command (recommended for a fresh VPS)

On the server (as `root`), after this file exists on `main` / `develop` in your clone:

```bash
bash /opt/pramool/pramool-core/deploy/vps/bootstrap-pramool.sh
```

Or copy only the script from your laptop, then run it (it will clone repos under `/opt/pramool` itself):

```bash
chmod +x bootstrap-pramool.sh
INSTALL_ROOT=/opt/pramool CORE_BRANCH=develop bash bootstrap-pramool.sh
```

The script will:

- Install Docker + Compose and open ports 3000, 3001, 3102–3104 (plus SSH / optional 80–443).
- Clone **four** repos over **HTTPS** (no GitHub SSH key on the VPS): `pramool-core`, `pramool-wallet-service`, `pramool-auction-service`, `pramool.in.th`.
- Write `deploy/vps/.env` with generated `JWT_SECRET`, DB password, internal API keys, and `NEXT_PUBLIC_*` URLs (public/LAN IP from `api.ipify.org` or `hostname -I`).
- Build all images (including the Next.js app with API base URLs), start Postgres, run `pramool-core migrate`, then `docker compose up -d`.

Open the site: **`http://<server-ip>:3000`**

Override branches if needed:

```bash
CORE_BRANCH=develop FRONTEND_BRANCH=main \
  WALLET_BRANCH=main AUCTION_BRANCH=main bash bootstrap-pramool.sh
```

## Option B — manual

### 1) Repo layout on VPS

Clone **from `/opt/pramool`** (not from inside `pramool-core`), so all repos are **siblings**:

```text
/opt/pramool/
  pramool-core/
  pramool-wallet-service/
  pramool-auction-service/
  pramool.in.th/
```

The compose file uses relative `build.context` to each repo. If you already cloned `pramool.in.th` under `pramool-core/`, either run `mv /opt/pramool/pramool-core/pramool.in.th /opt/pramool/`, or set in `.env`: `FRONTEND_BUILD_CONTEXT=../../pramool.in.th` (path is relative to `deploy/vps`).

### 2) Environment

```bash
cd /opt/pramool/pramool-core/deploy/vps
cp .env.example .env
```

Edit `.env`: set `DATABASE_*` / DSNs to hostname **`postgres`**, and set **`NEXT_PUBLIC_*`** to the URLs the **browser** will use (same host as the site, with ports 3001 / 3102 / 3103 for APIs). **CORS** must list your frontend origin (e.g. `http://your-ip:3000`).

Use **`chmod 0644 .env`** after saving so the non-root user inside the `pramool-core` image can read `/app/.env` (root-only `0600` causes `permission denied`).

Build and start:

```bash
docker compose build
docker compose up -d postgres
# wait until healthy
docker compose run --rm pramool-core ./pramool-core migrate
docker compose up -d
```

**Changing `NEXT_PUBLIC_*`:** values are baked at **`docker compose build`** time for the frontend. After edits:

```bash
docker compose build pramool-frontend --no-cache
docker compose up -d pramool-frontend
```

### 3) Verify

```bash
docker compose ps
docker compose logs -f pramool-frontend
docker compose logs -f pramool-core
```

## Notes

- `pramool-core` mounts `./.env` into the container as `/app/.env` because the binary uses Viper to read that file for `serve` and `migrate`.
- `pramool-auction-service` and `pramool-core` share the `uploads` volume (core serves `/uploads`, auction writes under `/srv/pramool/uploads`).
- `.env.example` uses `PGRST_DB_ANON_ROLE=postgres` for a quick lab; for production, create a restricted role (e.g. `web_anon`) and grants instead.
- Omise / ThaiBulkSMS: leave empty until you add real keys; wallet top-up / OTP will need them.
