# Migrations

All SQL migrations live here. Single track table: `bun_migrations`.

```bash
go run . migrate --db all
```

| Folder | Contents |
|--------|----------|
| `core/` | users, geo, categories, PDPA, notifications, appeals, … |
| `wallet/` | transactions, withdrawals |
| `auction/` | auctions, bids, reports, fees, … |
| `admin/` | platform_settings, site_announcements, admin moderation, platform_withdrawals, … |

Fresh install:

```bash
docker compose down -v && docker compose up -d postgres
go run . migrate --db all
```

**pramool-backoffice-core** no longer runs migrations — use this repo only.
