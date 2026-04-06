# Rules Resolution Service (Go + PostgreSQL)

A backend service that resolves workflow configuration using a **multi-dimensional override system**, based on specificity ranking.

---

## Features

- Clean Architecture
- Swagger for API display & Testing
- Specificity-based override resolution (similar to CSS cascade)
- Multi-dimensional selector (state, client, investor, caseType)
- Effective date handling (`effectiveDate`, `expiresDate`)
- Conflict detection between overrides
- Explain the API for debugging resolution decisions
- Override CRUD operations
- Seed data support (steps, defaults, overrides)

---

## Core Concept

Instead of maintaining multiple workflows, the system uses:

- **Defaults** → base configuration  
- **Overrides** → context-specific changes  

Resolution is based on:

1. Matching selector
2. Highest specificity
3. Latest effective date

---

## Architecture
````bash
API → Handler → Service → Repository (interface) → Infra (Postgres)
                    ↓
                 Domain Layer
          (Entities, Value Objects, Aggregates, Domain Services)
````
## Setup & Run

### 1. Clone repository

```bash
git clone <your-repo>
cd rules-resolution-service
````
### 2. Start the app
```bash
docker-compose up -d
````

# APIs
## Resolve Configuration
```http
POST /api/resolve
````
## Explain Resolution
```http
POST /api/resolve/explain
````
## Override APIs
```http
GET    /api/overrides
GET    /api/overrides/{id}
POST   /api/overrides
PUT    /api/overrides/{id}
PATCH  /api/overrides/{id}/status
````
## Conflict Detection
```http
GET /api/overrides/conflicts
````
## Status
This implementation covers core requirements, including resolution, conflict detection, and override management. Some optimizations and advanced features can be further extended.
