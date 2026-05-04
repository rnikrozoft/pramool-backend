# VPS Deployment (Docker Compose)

This stack deploys:

- `pramool-core`
- `pramool-wallet-service`
- `pramool-auction-service`
- `postgrest`

## 1) Required repo layout on VPS

Keep these repos as siblings:

```text
/opt/pramool/
  pramool-core/
  pramool-wallet-service/
  pramool-auction-service/
```

The compose file here relies on relative build contexts to sibling repos.

## 2) Prepare environment

```bash
cd /opt/pramool/pramool-core/deploy/vps
cp .env.example .env
```

Edit `.env` with real values (JWT/Omise/DB/etc).

## 3) Build and start

```bash
docker compose up -d --build
```

## 4) Verify

```bash
docker compose ps
docker compose logs -f pramool-core
docker compose logs -f pramool-wallet-service
docker compose logs -f pramool-auction-service
docker compose logs -f postgrest
```

## 5) Notes

- `pramool-auction-service` and `pramool-core` share the `uploads` volume:
  - core serves static files from `/app/uploads`
  - auction writes files to `/srv/pramool/uploads` via `PRAMOOL_UPLOAD_ROOT=/srv/pramool`
- `postgrest` needs an existing DB role in `PGRST_DB_ANON_ROLE` (for example `web_anon`) with proper grants.
- If your DB runs in another container/host, point `DATABASE_*`, `DATABASE_DSN`, `AUCTION_DATABASE_URL`, and `PGRST_DB_URI` to that reachable address.
