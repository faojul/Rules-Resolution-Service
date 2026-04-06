# Rules Resolution Service (Go + PostgreSQL)

A backend service that resolves workflow configuration using a **multi-dimensional override system**, based on specificity ranking.

---

## 🚀 Features

- Specificity-based override resolution (similar to CSS cascade)
- Multi-dimensional selector (state, client, investor, caseType)
- Effective date handling (`effectiveDate`, `expiresDate`)
- Conflict detection between overrides
- Explain API for debugging resolution decisions
- Override CRUD operations
- Seed data support (steps, defaults, overrides)

---

## 🧠 Core Concept

Instead of maintaining multiple workflows, the system uses:

- **Defaults** → base configuration  
- **Overrides** → context-specific changes  

Resolution is based on:

1. Matching selector
2. Highest specificity
3. Latest effective date

---

## 🏗 Architecture
`Handler → Service → Repository → PostgreSQL`


- **Handler** → HTTP layer
- **Service** → business logic
- **Repository** → DB interaction
- **Resolver** → core rule evaluation logic

---

## ⚙️ Setup & Run

### 1. Clone repository

```bash
git clone <your-repo>
cd rules-resolution-service
````
### 2. Start the app
```bash
docker-compose up -d
````

📡 APIs

