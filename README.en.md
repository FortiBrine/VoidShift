<div align="right">

🌐 **English** | [Українська](README.ua.md)

> Also available as [README.md](README.md) (default for GitHub)

</div>

<div align="center">

# VoidShift

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker&logoColor=white)](compose.yml)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20FreeBSD-orange)

VPN management panel for VPS servers. Single binary with an embedded web UI.
Currently supports WireGuard, with more protocols planned.

</div>

---

### ✨ Features

- 🔒 **WireGuard management** — create networks, bring interfaces up/down, add and remove peers
- 📱 **Peer configs and QR codes** — download `.conf` files or scan a QR code directly from the UI
- 📦 **Single binary** — the frontend is embedded into the Go binary at build time, one file to deploy
- 🗄️ **Dual database** — SQLite out of the box (zero config), or MySQL via DSN
- 🔑 **Session-based auth** — bootstrap admin credentials via environment variables
- 🐳 **Docker-ready** — `docker compose up` and you're running

---

### 🖥️ Supported Operating Systems

- ✅ **Linux** — fully supported
- 🔜 **FreeBSD** — planned
- ❌ **macOS** — not supported (WireGuard kernel APIs unavailable)
- ❌ **Windows** — not supported (WireGuard kernel APIs unavailable)

---

### 🚀 Quick Start

#### Docker (recommended)

```sh
cp .env.example .env
# Edit .env: set HOST_ADDRESS, ADMIN_USERNAME, ADMIN_PASSWORD at minimum
docker compose up -d
```

Open [http://localhost:8080](http://localhost:8080) and log in with your admin credentials.

> All required capabilities (`NET_ADMIN`, `/dev/net/tun`, `net.ipv4.ip_forward`) are pre-configured in `compose.yml`.

#### Manual Build

```sh
cp .env.example .env
# Edit .env

make deps     # install frontend and backend dependencies
make build    # build frontend, then embed it into the Go binary

sudo ./app    # NET_ADMIN privileges required for WireGuard
```

<details>
<summary>Other make targets</summary>

| Command | Description |
|---|---|
| `make deps` | `bun install` + `go mod download` |
| `make frontend` | Generate static SPA to `frontend/.output/public/` |
| `make backend` | Copy frontend output into embed dir, compile Go binary |
| `make build` | Full build: frontend + backend |
| `make clean` | Remove build artifacts |

</details>

---

### ⚙️ Configuration

All configuration is via environment variables. Copy `.env.example` to `.env`:

| Variable | Default | Required | Description |
|---|---|---|---|
| `SQLITE_DATABASE_PATH` | | One of these two | Path to SQLite file (e.g. `./database.db`) |
| `MYSQL_DSN` | | One of these two | MySQL DSN; if set, MySQL is used instead of SQLite |
| `HOST_ADDRESS` | | Yes | Public IP of this server, used in generated peer configs |
| `ADMIN_USERNAME` | | Yes | Bootstrap admin username |
| `ADMIN_PASSWORD` | | Yes | Bootstrap admin password |
| `HTTP_ADDRESS` | `:8080` | No | HTTP listen address |
| `GRACEFUL_TIMEOUT` | `5s` | No | Graceful shutdown timeout |
| `ENVIRONMENT` | | No | Set to `dev` for development logging |

---

### 📡 API Reference

All routes are under `/api`. Protected routes require an active session cookie.

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/health` | No | Health check |
| `POST` | `/api/auth/login` | No | Log in |
| `POST` | `/api/auth/logout` | No | Log out |
| `GET` | `/api/vpn/wireguard/networks` | Yes | List WireGuard networks |
| `POST` | `/api/vpn/wireguard/networks/generate` | Yes | Create a new network |
| `GET` | `/api/vpn/wireguard/networks/:id` | Yes | Get network details |
| `POST` | `/api/vpn/wireguard/networks/:id/up` | Yes | Bring interface up |
| `POST` | `/api/vpn/wireguard/networks/:id/down` | Yes | Bring interface down |
| `DELETE` | `/api/vpn/wireguard/networks/:id` | Yes | Delete network |
| `POST` | `/api/vpn/wireguard/networks/:id/peers/generate` | Yes | Add peer to network |
| `GET` | `/api/vpn/wireguard/peers/:peerId/config` | Yes | Get peer config text |
| `GET` | `/api/vpn/wireguard/peers/:peerId/config/download` | Yes | Download peer `.conf` file |
| `GET` | `/api/vpn/wireguard/peers/:peerId/qr` | Yes | Get peer config as QR code |
| `DELETE` | `/api/vpn/wireguard/peers/:peerId` | Yes | Remove peer |

---

### 📦 Project Structure

```
backend/               # Go module root (go.mod lives here)
  cmd/api/             # Binary entrypoint
  internal/
    app/               # Application wiring
    auth/              # Auth handler, service, middleware
    session/           # Session management (GORM-backed, 5-day TTL)
    user/              # User model, repository, service
    wireguard/         # WireGuard networks and peers
    shared/            # Config, database, HTTP plumbing, logger, validator
    embed/webui/       # Embedded frontend (auto-generated, do not edit)
frontend/              # Nuxt 4 SPA (SSR disabled, built with Bun)
  app/pages/           # Route pages
  i18n/                # Localisation files
```

---

### 📄 License

VoidShift is licensed under the [Apache License 2.0](LICENSE).
See [THIRD_PARTY_LIBRARIES.md](THIRD_PARTY_LIBRARIES.md) for open-source dependencies used in this project.
