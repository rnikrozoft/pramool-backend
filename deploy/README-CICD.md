# GitHub Actions → EC2 (private Docker registry)

This repo ships a production-oriented compose file: `deploy/docker-compose.prod.yml`. Images are built in CI, pushed to your registry, then the EC2 host pulls and restarts.

## One-time EC2 setup

1. Install Docker + Compose plugin, open ports (SSH, app ports, registry if needed).
2. Configure Docker on the server to trust your HTTP registry (if not using TLS):

   ```bash
   sudo mkdir -p /etc/docker
   echo '{"insecure-registries":["YOUR_HOST:5000"]}' | sudo tee /etc/docker/daemon.json
   sudo systemctl restart docker
   ```

3. Copy `deploy/docker-compose.prod.yml` to e.g. `/opt/pramool/docker-compose.prod.yml`.

   If you created `/opt/pramool` as **root**, either run once on the server  
   `sudo chown -R ubuntu:ubuntu /opt/pramool`  
   or rely on the workflow (it now runs `sudo mkdir` + `sudo chown` for the deploy user — requires passwordless `sudo` as on the default Ubuntu AMI).

   ```bash
   scp -i your-key.pem deploy/docker-compose.prod.yml ubuntu@YOUR_EC2:/opt/pramool/docker-compose.prod.yml
   ```

4. Create `/opt/pramool/.env` on the server with real values (DB, JWT, Omise, `DATABASE_DSN`, `AUCTION_DATABASE_URL`). The frontend image bakes `NEXT_PUBLIC_*` at **build** time in CI — set those as GitHub **Variables** on `pramool.in.th` (see below), not necessarily in `.env` on the server.

5. Export registry + tag when deploying:

   ```bash
   export REGISTRY_HOST=YOUR_HOST:5000
   export IMAGE_TAG=latest
   docker login "$REGISTRY_HOST"
   ```

## GitHub Actions secrets (repository or environment)

| Secret | Purpose |
|--------|---------|
| `REGISTRY_HOST` | Host:port only, e.g. `43.210.19.172:5000` (no `http://`) |
| `REGISTRY_USER` | Registry basic auth user |
| `REGISTRY_PASSWORD` | Registry password |
| `EC2_HOST` | Public IPv4 or DNS |
| `EC2_USER` | e.g. `ubuntu` |
| `EC2_SSH_KEY` | Private key (full PEM) for SSH deploy |
| `EC2_DEPLOY_PATH` | Optional. Default `/opt/pramool` |
| `PRAMOOL_DOTENV_B64` | Optional. Base64 of the full `.env` for the server; if set, the workflow overwrites `/opt/pramool/.env` on each deploy |

### Optional: build args for frontend (repo `pramool.in.th`)

Set these **repository variables** (or secrets) so the Next.js image gets public API URLs:

- `NEXT_PUBLIC_CORE_API_BASE_URL`
- `NEXT_PUBLIC_USER_API_BASE_URL`
- `NEXT_PUBLIC_WALLET_API_BASE_URL`
- `NEXT_PUBLIC_AUCTION_REALTIME_BASE_URL`

## Workflows

Each microservice repo has `.github/workflows/deploy-ec2.yml`:

- Builds the Docker image, tags `latest` and `$GITHUB_SHA` (short).
- Pushes to `$REGISTRY_HOST/pramool-<service>:<tag>`.
- SSH into EC2, `docker login`, `docker compose pull` for that service, `up -d`.

`docker compose up -d <service>` will still start **dependent** services from the same compose file (e.g. Postgres) when they are not already running.

Image names:

| Repo | Image |
|------|--------|
| `pramool-core` | `pramool-core` |
| `pramool-wallet-service` | `pramool-wallet` |
| `pramool-auction-service` | `pramool-auction` |
| `pramool.in.th` | `pramool-frontend` |

After changing only one service, run that repo’s workflow. To roll out a pinned tag, set `IMAGE_TAG` on the server (or extend the workflow to pass the sha).

## GitHub-hosted runner + HTTP registry

The workflow adds your registry to `insecure-registries` on the runner so `docker push` works over plain HTTP. For production, prefer HTTPS (reverse proxy + TLS) and drop that step.

## Health check + rollback (deploy script)

There is no single GitHub Action that covers *private HTTP registry + Docker Compose on one EC2 + rollback* end-to-end. Common production options are **Kamal** (37signals, Ruby) or **ECS/Kubernetes**, which are heavier than this demo setup.

Here we use a small **bash script** (`deploy/ec2-deploy-compose.sh`) — the same pattern many teams use (immutable image tag + `curl` health + re-deploy previous tag):

1. CI pushes images tagged `latest` **and** short git SHA (`metadata-action` + same 7-char SHA as `IMAGE_TAG` on the server).
2. On EC2, the script sets `IMAGE_TAG=<sha>`, `docker compose pull` + `up -d` for **one** service.
3. It polls a **health URL** (Go services: `/healthz`; frontend: `http://127.0.0.1:3000/`).
4. On failure, it reads `/opt/pramool/.last-good-<service>` and redeploys the **previous successful** tag.
5. On success, it writes the new SHA to that file.

**First deploy:** there is no previous tag file yet — a failed first deploy cannot roll back automatically.

Workflows upload the script with **`appleboy/scp-action`** (widely used) and run it over **`appleboy/ssh-action`**.

`scp-action` archives paths with directory structure; uploading `deploy/ec2-deploy-compose.sh` can extract to `.../deploy/ec2-deploy-compose.sh` on the server. Workflows **copy the script to the repo root** before SCP so it lands exactly as `/opt/pramool/deploy/ec2-deploy-compose.sh`.

Per-repo health URLs are set in each workflow as `DEPLOY_HEALTH_URL`. If you change host ports in `.env`, update the workflow `env` block to match.

## Registry garbage-collect after deploy

Each `deploy-ec2.yml` ends with an SSH step that runs **Docker Distribution** garbage collection on the host you deploy to (`EC2_HOST`). It expects a running container named `registry` (same machine as your app stack).

- Removes **unreferenced blobs** (and, when supported, **untagged manifests** via `--delete-untagged`).
- Uses `continue-on-error: true` so a GC quirk does not fail the workflow.
- Optional **repository variable** `REGISTRY_CONTAINER_NAME` if your registry container is not named `registry`.

**Assumption:** the private registry runs on the **same EC2** as `docker compose` (typical for a single demo server). If the registry is on another host, remove that step or point SSH at the registry host.

**Note:** If all four repos run deploy around the same time, GC may run up to four times; harmless but redundant. You can delete the GC step from three workflows and keep it in one repo only if you prefer.
