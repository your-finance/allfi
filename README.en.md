# AllFi — All-Asset Aggregation Platform

> A personal asset dashboard for Web3 professionals — one page, multiple currencies, full picture.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)

[中文文档](./README.md)

---

## About

AllFi is an **open-source, self-hosted** all-asset aggregation platform that unifies your crypto and traditional assets:

- **CEX Exchanges**: Binance, OKX, Coinbase
- **On-chain Wallets**: Ethereum, BSC, Polygon (+ Arbitrum/Optimism/Base queries)
- **DeFi Protocols**: Lido, RocketPool, Aave, Compound, Uniswap V2/V3, Curve
- **NFT Collections**: Alchemy integration for browsing and valuation
- **Traditional Assets**: Bank deposits, cash, stocks, mutual funds

All data is stored locally. API keys are encrypted with AES-256-GCM and never leave your machine.

---

## Key Features

| Category | Features |
|----------|----------|
| Asset Aggregation | CEX + on-chain + DeFi + NFT + traditional assets in one view |
| Multi-currency Pricing | USDC / BTC / ETH / CNY — switch freely |
| Transaction History | Unified CEX + on-chain records with incremental sync and cursor pagination |
| Analytics | Daily PnL, fee analytics, attribution analysis, benchmark comparison (vs BTC/ETH/S&P500) |
| Strategy Engine | Target allocation + rebalance suggestions |
| Reports | Auto-generated daily / weekly / monthly / annual reports |
| Achievements | 11 investment achievement badges |
| Notifications | Price alerts + WebPush browser notifications |
| Privacy Mode | One-click blur all amounts for screen sharing |
| Themes | 4 professional financial themes (3 dark + 1 light) |
| Languages | Simplified Chinese / Traditional Chinese / English |
| PWA | Add to home screen, works offline |

---

## Quick Start

### Docker (Recommended)

```bash
git clone https://github.com/your-finance/allfi.git
cd allfi

# Configure environment variables
cp .env.example .env
# Edit .env — set ALLFI_MASTER_KEY (generate with: openssl rand -base64 32)

# Start services
docker-compose up -d
```

Visit http://localhost:3174 to get started. First-time access requires setting a PIN code (4–8 digits).

### Manual Development

```bash
# Backend
cd core
go run cmd/server/main.go     # http://localhost:8080

# Frontend (new terminal)
cd webapp
pnpm install && pnpm dev      # http://localhost:3174
```

> See the [Deployment Guide](./docs/guides/deployment-guide.md) for detailed instructions.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.24 + standard library `net/http` + GORM v2 + SQLite3/MySQL |
| Frontend | Vue 3.5 + Vite 7.3 + Pinia 3 + Chart.js 4 + Phosphor Icons |
| Authentication | PIN code bcrypt + JWT Bearer Token |
| Encryption | AES-256-GCM (API key storage) |
| Deployment | Docker Compose (read-only containers + no-new-privileges + healthcheck) |
| API Docs | OpenAPI 3.0 + Swagger UI (`/api/v1/docs`) |

### Architecture

```
Frontend (Vue 3 / Vite)
    ↓ RESTful API (~75 routes)
Backend (Go / net/http)
    ├── api/handlers/     HTTP handlers (24)
    ├── service/          Business logic (16 services)
    ├── repository/       Data access (18 repos)
    └── integrations/     Third-party integrations (8 modules)
    ↓
Data Layer (SQLite3 / MySQL, 18 tables)
```

---

## Project Structure

```
allfi/
├── core/                       # Backend (Go)
│   ├── cmd/server/main.go      # Entry point
│   ├── internal/
│   │   ├── api/handlers/       # HTTP handlers
│   │   ├── service/            # Business logic
│   │   ├── repository/         # Data access
│   │   ├── models/             # Data models
│   │   └── integrations/       # Third-party integrations
│   └── manifest/config/        # Configuration
├── webapp/                     # Frontend (Vue 3)
│   └── src/
│       ├── pages/              # 9 pages
│       ├── components/         # 39 components
│       ├── stores/             # 12 Pinia stores
│       └── composables/        # 8 composables
└── docs/                       # Documentation
    ├── product/                # Product docs
    ├── tech/                   # Technical docs
    ├── specs/                  # Requirement specs
    ├── design/                 # Design docs
    └── guides/                 # Guides
```

---

## Documentation

Full documentation index: [docs/README.md](./docs/README.md)

| Category | Documents |
|----------|-----------|
| Product | [Business Overview](./docs/product/biz-overview.md) · [Feature Overview](./docs/product/feature-overview.md) |
| Technical | [Tech Baseline](./docs/tech/tech-baseline.md) · [API Reference](./docs/tech/api-reference.md) · [Swagger UI](http://localhost:8080/api/v1/docs) |
| Guides | [Deployment Guide](./docs/guides/deployment-guide.md) · [Dev Guide](./docs/guides/dev-guide.md) · [User Guide](./docs/guides/user-guide.md) |
| Design | [UI/UX Standards](./docs/design/ui-ux-standards.md) · [i18n System](./docs/design/i18n.md) |
| Specs | [Frontend Spec](./docs/specs/frontend-spec.md) · [Backend Spec](./docs/specs/backend-spec.md) |

---

## Security

- API keys encrypted with **AES-256-GCM** — no plaintext in the database
- PIN code hashed with **bcrypt** — irreversible
- Fully **self-hosted** — data never leaves your server
- Recommend **read-only** API key permissions — no withdrawal/trading access
- Docker containers run as **non-privileged + read-only**

---

## Contributing

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/xxx`)
3. Follow coding conventions (see [CLAUDE.md](./CLAUDE.md))
4. Write tests (target coverage >70%)
5. Submit a Pull Request

---

## License

[MIT License](LICENSE)

---

**Built in 2026 for Web3 professionals.**

- GitHub: https://github.com/your-finance/allfi
- Issues: https://github.com/your-finance/allfi/issues
