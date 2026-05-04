# VPS Deployment (Docker Compose)

This stack deploys:

- `postgres` (16, data in Docker volume `pgdata`)
- `pramool-core`
- `pramool-wallet-service`
- `pramool-auction-service`
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

- Install Docker + Compose and open ports 3001, 3102–3104 (plus SSH / optional 80–443).
- Clone the three repos over **HTTPS** (no GitHub SSH key on the VPS).
- Write `deploy/vps/.env` with generated `JWT_SECRET`, DB password, and internal API keys.
- Start Postgres, run `pramool-core migrate`, then `docker compose up -d`.

Override branches if needed:

```bash
CORE_BRANCH=develop WALLET_BRANCH=main AUCTION_BRANCH=main bash bootstrap-pramool.sh
```

## Option B — manual

### 1) Repo layout on VPS

Keep these repos as siblings:

```text
/opt/pramool/
  pramool-core/
  pramool-wallet-service/
  pramool-auction-service/
```

The compose file relies on relative build contexts to sibling repos.

### 2) Environment

```bash
cd /opt/pramool/pramool-core/deploy/vps
cp .env.example .env
```

Edit `.env`: set `DATABASE_*` / DSNs to hostname **`postgres`** (the Compose service), not `localhost`.  
Use **`chmod 0644 .env`** after saving so the non-root user inside the `pramool-core` image can read `/app/.env` (root-only `0600` causes `permission denied`).

Run migrations before or after first boot:

```bash
docker compose build
docker compose up -d postgres
# wait until healthy
docker compose run --rm pramool-core ./pramool-core migrate
docker compose up -d
```

### 3) Verify

```bash
docker compose ps
docker compose logs -f pramool-core
```

## Notes

- `pramool-core` mounts `./.env` into the container as `/app/.env` because the binary uses Viper to read that file for `serve` and `migrate`.
- `pramool-auction-service` and `pramool-core` share the `uploads` volume (core serves `/uploads`, auction writes under `/srv/pramool/uploads`).
- `.env.example` uses `PGRST_DB_ANON_ROLE=postgres` for a quick lab; for production, create a restricted role (e.g. `web_anon`) and grants instead.
- Omise / ThaiBulkSMS: leave empty until you add real keys; wallet top-up / OTP will need them.
