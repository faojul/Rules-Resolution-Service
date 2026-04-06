# Rules Resolution Service (Go + PostgreSQL)

A backend service that resolves configuration using a **multi-dimensional override system**, inspired by CSS specificity.

## 🚀 Features

- Specificity-based override resolution (0–4 dimensions)
- Effective date filtering (`effectiveDate`, `expiresDate`)
- Conflict detection between overrides
- Explain API for debugging resolution
- Override CRUD with audit trail
- Bulk resolution support

---

## 🧠 Core Idea

Instead of maintaining multiple workflows, we define:
- **Defaults** (base config)
- **Overrides** (context-specific changes)

Resolution is based on:
1. Matching selector
2. Highest specificity
3. Latest effective date

---

## 📦 Tech Stack

- Go 1.21+
- PostgreSQL
- pgxpool
- chi router

---

## ⚙️ Run Locally

```bash
git clone <repo>
cd project

docker-compose up -d

go run cmd/main.go
