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
├── migrations/           # Raw SQL migration and rollback scripts
│   ├── 20250524121000_tableName.up.sql
│   ├── 20250524121000_tableName.down.sql
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
| `migrate`             | Apply all `.up.sql` migration scripts located in `migrations/`             |
| `rollback`            | Revert the latest migration using `.down.sql` scripts (⚠️ irreversible)    |
| `newMigration <tableName>` | Generate a migration file pair named `<tableName>.up.sql` and `<tableName>.down.sql` |

---

## 🚀 Example Usages
🔧 Create a New Migration

```cmd
go run . newMigration users
```

Creates two files in migrations/:

20250524121000_users.up.sql: define CREATE, INSERT, or ALTER TABLE

20250524121000_users.down.sql: define rollback logic (DROP TABLE, etc.)

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

> ⚠️ Important: Never commit your .env file to version control (it is already ignored via .gitignore).

---

## 🧩 Usage in the Project
The .env file is used in the following places:

✅ docker-compose.yml

Docker Compose reads the .env file to inject environment variables into containers.

These variables configure the PostgreSQL database and the Go backend service.

Example (from docker-compose.yml):
```yml
environment:
  - DATABASE_HOST=${DATABASE_HOST}
  - DATABASE_PORT=${DATABASE_PORT}
  - DATABASE_USERNAME=${DATABASE_USERNAME}
  - DATABASE_PASSWORD=${DATABASE_PASSWORD}
  - DATABASE_NAME=${DATABASE_NAME}
```

✅ .github/workflows/deployment.yml

In the GitHub Actions CI/CD pipeline, environment variables are reconstructed into an .env file during the workflow using GitHub Secrets.

These values are then copied to the VPS for use in Docker Compose on the server.

Example step:
```yml
- name: Write .env file
  run: |
    echo "DATABASE_HOST=${DATABASE_HOST}" >> .env
    ...
```
By properly setting up your .env file, both local development and production deployment via CI/CD will work seamlessly.