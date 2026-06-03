# Pramool Backend

A backend service for an auction system built with Go, following Clean Architecture principles. This project includes modular components, database migration scripts, and command-based CLI tools.

---

## 📁 Project Structure

```plaintext
.
├── cmd/                  # Entry point for CLI commands (e.g., serve, migrate, rollback, newMigration)
├── config/               # Application configuration setup (e.g., environment loading)
├── const/                # Constant values used across the project
├── docs/                 # Swagger documentation files and generators
├── exception/            # Custom error definitions and utilities
├── handler/              # HTTP route handlers (controllers)
├── mapping/              # Mapping between entities and DTOs
├── migrations/           # SQL schema แยกโฟลเดอร์ (รันผ่าน pramool-core → database เดียว)
│   ├── core/             # users, banks, ที่อยู่
│   ├── wallet/           # transactions, withdrawals
│   ├── auction/          # auctions, bids, bid_transactions
│   └── migration.go
├── model/
│   ├── dto/              # Data Transfer Objects used in APIs
│   └── entity/           # Database models (entities)
├── repository/           # Data access layer (e.g., PostgreSQL repositories)
├── service/              # Business logic layer
├── .env                  # Environment variables (excluded from version control)
├── .gitignore            # Files to be ignored by Git
├── docker-compose.yml    # Docker service definitions
├── Dockerfile            # Docker build instructions
├── go.mod                # Go module metadata
├── go.sum                # Go module checksums
├── main.go               # Main entry point (delegates to cmd/root.go)
└── readme.md             # Project documentation
```

---

## ⚙️ Commands & Usage
This project uses go run . [command] to execute backend tasks via CLI. Available commands:

| Command               | Description                                                                |
| --------------------- | -------------------------------------------------------------------------- |
| `serve`               | Start the backend web server                                               |
| `migrate --db <core\|wallet\|auction\|all>` | Apply migrations from `migrations/<db>/` to `DATABASE_NAME` (database เดียว) |
| `rollback --db <core\|wallet\|auction>` | Revert the latest migration on that database (⚠️ irreversible) |
| `newMigration <db> <name>` | Create a pair under `migrations/<db>/` (db = core, wallet, auction) |

---

## 🚀 Example Usages
🔧 Create a New Migration

```cmd
go run . newMigration core users
```

Creates `migrations/core/<timestamp>_users.up.sql` and `.down.sql` — use `CREATE TABLE` in `.up`, `DROP TABLE` in `.down`.

> ⚠️ Use .down.sql cautiously — it may delete data or schema.

---

## 🧪 Environment Configuration (.env)
The project requires an .env file to define environment-specific variables. You need to create this file manually in the project root directory.

📝 Sample .env Format:

### Database Configuration
```env
DATABASE_HOST=
DATABASE_PORT=
DATABASE_USERNAME=
DATABASE_PASSWORD=
DATABASE_NAME=
```

### JWT Configuration
```env
JWT_SECRET=
JWT_EXPIRE_TIME=
```

### Omise Configuration
```env
OMISE_SECRET_KEY=
OMISE_WEBHOOK_SECRET=
```

### Wallet fees (`pramool-wallet-service` + frontend fallback)
```env
# ขั้นต่ำยอดชำระ Omise ตอนเติมเครดิต (บาท)
WALLET_MIN_TOPUP_GROSS_THB=100
# ขั้นต่ำเครดิตที่ถอนได้ต่อครั้ง (บาท)
WALLET_MIN_WITHDRAW_CREDIT_THB=100
# ค่าธรรมเนียม PromptPay ในหน่วย ppm (17655 ≈ 1.65% + VAT 7% บนค่าธรรมเนียม)
WALLET_OMISE_PROMPTPAY_FEE_PPM=17655
# ค่าธรรมเนียมโอนเข้าธนาคาร Omise ต่อครั้ง (บาท) — หักจากยอดที่ผู้ใช้ได้รับ
WALLET_OMISE_TRANSFER_FEE_THB=21
```

### Auction platform commission (`pramool-auction-service` + `GET /wallet/fees`)
```env
# ค่าคอมมิชชันแพลตฟอร์มเมื่อปิดตามเวลา (%)
AUCTION_PLATFORM_FEE_NORMAL_PCT=25
# ค่าคอมมิชชันเมื่อปิดก่อนเวลา (%)
AUCTION_PLATFORM_FEE_EARLY_PCT=30
# ส่วนที่ผู้ขายได้หลังผู้ซื้อยืนยันรับของ — ปิดตามเวลา (%)
AUCTION_SELLER_KEEP_NORMAL_PCT=75
# ส่วนที่ผู้ขายได้เมื่อปิดก่อนเวลา (%)
AUCTION_SELLER_KEEP_EARLY_PCT=70
```

> ⚠️ Important: Never commit your .env file to version control (it is already ignored via .gitignore).

---

## 🧩 Usage in the Project
The .env file is used in the following places:

✅ docker-compose.yml

Docker Compose reads the .env file to inject environment variables into containers.

These variables configure the PostgreSQL database, Redis (live auction bidders), and the Go backend services.

`pramool-auction-service` uses `REDIS_URL` (default in compose: `redis://redis:6379/0`).

## Database (PostgreSQL ตัวเดียว)

SQL ทั้งหมดอยู่ที่ `pramool-core/migrations/` แยกโฟลเดอร์ตาม domain แต่ **รันลง database เดียว** (`DATABASE_NAME` เช่น `pramool`)

| `--db` | โฟลเดอร์ SQL | ตารางหลัก |
|--------|--------------|-----------|
| `core` | `migrations/core/` | users, geo, categories, PDPA, notifications, appeals, … |
| `wallet` | `migrations/wallet/` | `transactions`, `withdrawals` |
| `auction` | `migrations/auction/` | auctions, bids, reports, fees, … |
| `admin` | `migrations/admin/` | platform_settings, announcements, admin moderation, … |
| `all` | ทั้ง 4 โฟลเดอร์ (เรียง timestamp รวม) | **ใช้ตอน setup / deploy** |

ทุก service ใช้ `DATABASE_NAME` / `DATABASE_DSN` ชี้ database เดียว — migrate ครั้งเดียวที่ repo นี้

ทุก service (`pramool-core`, `wallet`, `auction`) ใช้ `DATABASE_NAME` / `DATABASE_DSN` ชี้ database เดียวกัน

Fresh install:

```bash
docker compose down -v   # ลบ volume เก่า (ถ้ามี 3 DB แยก)
docker compose up -d postgres
cd pramool-core && go run . migrate --db all
```

```env
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=...
DATABASE_NAME=pramool
# ทางเลือก: DSN เต็ม (migrate อ่านตัวนี้ก่อน)
# DATABASE_DSN=postgres://postgres:...@localhost:5432/pramool?sslmode=disable
```

`pramool-wallet-service`: top-up credits net of Omise PromptPay fee; withdrawals deduct transfer fee from payout.

`pramool-auction-service`: เมื่อปลด escrow ให้ผู้ขาย — ส่วนแบ่งผู้ขาย = `trunc(winner × seller_keep% / 100)` เศษทั้งหมดบันทึกใน `platform_sale_fees` (ไม่หายจากระบบ).

```bash
docker compose up -d redis              # Redis only (localhost:6379)
docker compose up -d redis redis-insight   # + dashboard http://localhost:5540
```

**Redis Insight** (`http://localhost:5540`): Insight รัน**ใน Docker** — การเชื่อมทำจาก container ของ Insight ไม่ใช่จาก Mac

| ใช้เมื่อ | Host / URL |
|----------|------------|
| เพิ่ม DB ใน Redis Insight (compose) | Host **`redis`**, Port **6379** หรือ `redis://redis:6379` |
| `redis-cli` / auction-service บน Mac | `127.0.0.1:6379` หรือ `redis://localhost:6379/0` |

อย่าใส่ `127.0.0.1` ใน Insight — จะชี้ไปที่ container ตัวเองแล้ว error  
compose ตั้ง `RI_REDIS_HOST=redis` ให้แล้ว — หลัง `docker compose up -d redis-insight` ควรเห็น connection **pramool-local** อัตโนมัติ

Example (from docker-compose.yml):
```yml
environment:
  - DATABASE_HOST=${DATABASE_HOST}
  - DATABASE_PORT=${DATABASE_PORT}
  - DATABASE_USERNAME=${DATABASE_USERNAME}
  - DATABASE_PASSWORD=${DATABASE_PASSWORD}
  - DATABASE_NAME=${DATABASE_NAME}
```

✅ `.github/workflows/deploy-ec2.yml`

CI builds and pushes images to your private registry, then SSHs into EC2 and runs `deploy/ec2-deploy-compose.sh` (health check + rollback). Server `.env` can be supplied via the `PRAMOOL_DOTENV_B64` secret or maintained on the host. See `deploy/README-CICD.md`.

By properly setting up your `.env` file (local) and secrets (CI / server), both local development and production deployment work as intended.