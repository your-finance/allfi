# AllFi — All-Asset Aggregation Platform

> A personal asset dashboard for Web3 professionals — one page, multiple currencies, full picture.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![GoFrame](https://img.shields.io/badge/GoFrame-v2.10-blue.svg)](https://goframe.org/)
[![Vue](https://img.shields.io/badge/Vue-3.5-brightgreen.svg)](https://vuejs.org/)
[![Vite](https://img.shields.io/badge/Vite-7.3-646CFF.svg)](https://vite.dev/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4-06B6D4.svg)](https://tailwindcss.com/)

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

### Option 1: Docker Deployment (Recommended) 🐳

**Only requires Docker — no need to install Go / Node.js / pnpm locally.**

Prerequisites: Docker 20.10+, Docker Compose v2+

#### One-click Script Deployment

```bash
# Clone repo and deploy
git clone https://github.com/your-finance/allfi.git
cd allfi
bash deploy/docker-deploy.sh
```

The script automatically: checks Docker environment → generates `.env` + security keys → builds and starts all services.

#### Manual Docker Deployment

```bash
git clone https://github.com/your-finance/allfi.git
cd allfi

# Generate .env (required for first run)
cp .env.example .env
# Edit .env — set ALLFI_MASTER_KEY (or auto-generate with the line below)
sed -i "s|CHANGE_ME_USE_openssl_rand_base64_32|$(openssl rand -base64 32)|" .env

# Start services
docker compose up -d --build
```

Visit http://localhost:3174 to get started. First-time access requires setting a PIN code (4–8 digits).

```bash
# Common Docker commands
docker compose logs -f       # View logs
docker compose down          # Stop services
docker compose restart       # Restart services
docker compose up -d --build # Rebuild and restart
```

### Option 2: Local Development

For developers who need to modify the code. Requires: Go 1.24+, Node.js 20+, pnpm.

```bash
git clone https://github.com/your-finance/allfi.git
cd allfi
make setup    # Generate .env + install dependencies
make dev      # Start both backend and frontend dev servers
```

Visit http://localhost:3174 to get started. First-time access requires setting a PIN code (4–8 digits).

> **Note**: `make setup` auto-detects your environment. If Go or pnpm is missing, it will skip the corresponding dependency installation and show a warning.

### Option 3: Mock Mode (No Backend)

Just want to see the UI? No backend needed. Requires: Node.js 20+, pnpm.

```bash
cd allfi
cd webapp && pnpm install && pnpm dev:mock
```

Visit http://localhost:3174 — all data is simulated.

> See the [Deployment Guide](./docs/guides/deployment-guide.md) for detailed instructions.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.24 + GoFrame v2.10 + GoFrame ORM + SQLite3 |
| Frontend | Vue 3.5 + Vite 7.3 + Tailwind CSS 4 + Pinia 3 + Chart.js 4 + Phosphor Icons + VueUse |
| Authentication | PIN code bcrypt + JWT Bearer Token |
| Encryption | AES-256-GCM (API key storage) |
| Deployment | Docker Compose (read-only containers + no-new-privileges + healthcheck) |
| API Docs | OpenAPI 3.0 + Swagger UI (`/api/v1/docs`) |

### Architecture

```
Frontend (Vue 3.5 / Vite 7.3 / Tailwind CSS 4)
    ↓ RESTful API
Backend (Go 1.24 / GoFrame v2.10)
    ├── api/              API definitions (RESTful)
    ├── app/              Business modules (26 modules)
    │   ├── controller/   Controllers
    │   ├── logic/        Business logic
    │   └── service/      Service interfaces
    └── integrations/     Third-party integrations (8 modules)
    ↓
Data Layer (GoFrame ORM + SQLite3, 26 entities)
```

---

## Project Structure

```
allfi/
├── core/                       # Backend (Go + GoFrame v2.10)
│   ├── cmd/server/main.go      # Entry point
│   ├── api/v1/                 # API definitions (RESTful)
│   ├── internal/
│   │   ├── app/                # Business modules (26)
│   │   │   └── {module}/
│   │   │       ├── controller/ # Controllers
│   │   │       ├── logic/      # Business logic
│   │   │       └── service/    # Service interfaces
│   │   ├── model/entity/       # Data models (26 entities)
│   │   └── integrations/       # Third-party integrations (8 modules)
│   └── manifest/config/        # Configuration
├── webapp/                     # Frontend (Vue 3.5 + Tailwind CSS 4)
│   └── src/
│       ├── pages/              # 9 pages
│       ├── components/         # 39 components
│       ├── stores/             # 13 Pinia stores
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
3. Write tests (target coverage >70%)
4. Submit a Pull Request

---

## License

[MIT License](LICENSE)

---

**Built in 2026 for Web3 professionals.**

- GitHub: https://github.com/your-finance/allfi
- Issues: https://github.com/your-finance/allfi/issues
