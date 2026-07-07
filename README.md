<div align="center">

# VoidShift

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker&logoColor=white)](compose.yml)
[![License](https://img.shields.io/badge/license-MIT%20%2F%20Apache--2.0-blue.svg)](#license)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20FreeBSD%20%7C%20OpenBSD%20%7C%20macOS%20%7C%20Windows-orange)

VPN management panel for VPS servers. Single binary with an embedded, server-rendered web UI.
Currently supports WireGuard, with more protocols planned.

</div>

---

### Features

- Create and manage WireGuard networks and peers; bring interfaces up or down from the web UI or API
- Download peer `.conf` files or display them as a QR code
- The web UI is server-rendered (templ + Tailwind) and compiled into the binary at build time, one file to deploy
- Admin credentials come from environment variables and are bootstrapped on startup
- Ships with `compose.yml` for one-command deployment

---

### Supported operating systems

| OS | Status | Notes |
|---|---|---|
| Linux | Stable | Full support via `netlink` + `wgctrl` |
| FreeBSD / OpenBSD | Experimental | Userspace WireGuard via `golang.zx2c4.com/wireguard`; OpenBSD interface names are auto-generated (`tun0`, `tun1`, …) since `tun(4)` can't be renamed |
| macOS (Darwin) | Experimental | Userspace WireGuard via `golang.zx2c4.com/wireguard`; requires `utun` interface support |
| Windows | Experimental | Userspace WireGuard via `golang.zx2c4.com/wireguard`; address configuration via the IP Helper API (`winipcfg`) |

> Experimental platforms compile and run but haven't been tested in production. Interface state is in-memory only and lost on process restart; re-running "up" after a restart may conflict with the still-running kernel interface.

---

### Quick start

#### Docker (recommended)

```sh
cp .env.example .env
# Edit .env: set HOST_ADDRESS, ADMIN_USERNAME, ADMIN_PASSWORD at minimum
docker compose up -d
```

Open [http://localhost:8080](http://localhost:8080) and log in with your admin credentials.

> All required capabilities (`NET_ADMIN`, `/dev/net/tun`, `net.ipv4.ip_forward`) are pre-configured in `compose.yml`.

#### Manual build

Builds run through [mise](https://mise.jdx.dev): it pins the exact Go and Bun versions
(`mise.toml`) and downloads them itself on first run, so nothing needs to be pre-installed —
same command on Linux, macOS, or Windows.

```sh
cp .env.example .env
# Edit .env

curl https://mise.run | sh   # or: winget install jdx.mise / scoop install mise / brew install mise

mise run build   # generate (sqlc + templ) -> bundle frontend deps (assets/ + templui via bun) -> go build -o app ./cmd/api

sudo ./app    # NET_ADMIN privileges required for WireGuard
```

Other tasks: `mise run install`, `mise run generate`, `mise run assets`, `mise run run`,
`mise run clean` — see `mise.toml` for the full task graph.

---

### Configuration

All configuration is via environment variables (loaded from `.env` if present). Copy `.env.example` to `.env`:

| Variable | Default | Description |
|---|---|---|
| `SQLITE_DATABASE_PATH` | `./store.db` | Path to the SQLite database file |
| `HOST_ADDRESS` | `1.2.3.4` | Public IP of this server, embedded in generated peer configs |
| `HTTP_ADDRESS` | `:8080` | HTTP listen address |
| `GRACEFUL_TIMEOUT` | `5s` | Graceful shutdown timeout |
| `ADMIN_USERNAME` | `admin` | Bootstrap admin username (upserted on every boot) |
| `ADMIN_PASSWORD` | `password` | Bootstrap admin password (upserted on every boot) |
| `ENVIRONMENT` | `dev` | Set to `prod` to disable debug logging and enforce the session cookie's `Secure` flag |

---

### API reference

The JSON API lives under `/api`. The web UI itself is server-rendered separately (`/`, `/login`, `/wireguard/...`) and isn't part of this API. Protected routes require an active session cookie.

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
| `GET` | `/api/vpn/wireguard/peers/:peerId/qr` | Yes | Get peer config as a QR code image |
| `DELETE` | `/api/vpn/wireguard/peers/:peerId` | Yes | Remove peer |

---

### Project structure

```
cmd/api/                 # Binary entrypoint (main.go)
internal/
  app/                   # Application wiring, route registration
  auth/                  # Auth handler, service, JSON API routes
  user/                  # Bootstrap admin user, repository, service
  wireguard/             # WireGuard networks/peers: netlink + wgctrl, JSON API routes
  webui/                 # Server-rendered UI routes/handlers, embedded static assets
  config/                # Environment-variable configuration
  middleware/            # Auth, error handling
  store/                 # sqlc-generated DB layer + goose migration runner
  validator/             # go-playground/validator adapter for Fiber
  logger/                # slog setup
view/
  pages/                 # templ page templates
  layouts/               # templ layout templates
  components/ui/         # vendored templui component set
db/
  migrations/            # goose SQL migrations
  queries/               # sqlc input queries
assets/                  # Frontend dependency sources (Tailwind v4 entrypoint), bundled via bun into internal/webui/static/app.css
mise.toml                # Pinned Go/Bun toolchain + build tasks (mise run build)
```

---

### Roadmap

- MySQL as an alternative to SQLite
- Additional VPN protocols beyond WireGuard
- Stable non-Linux support (persistent interface state across restarts)

---

### License

VoidShift is dual-licensed under [MIT](LICENSE-MIT) or [Apache License 2.0](LICENSE-APACHE), at your option.
