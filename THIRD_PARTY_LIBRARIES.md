# Third-Party Libraries

VoidShift is built on top of the following open-source libraries and frameworks. We are grateful to their authors and contributors.

---

## Backend (Go)

| Library | Version | License | Purpose |
|---|---|---|---|
| [labstack/echo/v5](https://github.com/labstack/echo) | v5.1.0 | MIT | HTTP web framework |
| [gorm.io/gorm](https://github.com/go-gorm/gorm) | v1.31.1 | MIT | ORM for database access |
| [gorm.io/driver/mysql](https://github.com/go-gorm/mysql) | v1.6.0 | MIT | GORM MySQL driver |
| [glebarez/sqlite](https://github.com/glebarez/sqlite) | v1.11.0 | MIT | Pure-Go SQLite driver (no CGO) |
| [golang.zx2c4.com/wireguard/wgctrl](https://github.com/WireGuard/wgctrl-go) | v0.0.0-20241231184526 | MIT | WireGuard kernel interface control |
| [vishvananda/netlink](https://github.com/vishvananda/netlink) | v1.3.1 | Apache 2.0 | Linux network interface management via netlink |
| [go-playground/validator/v10](https://github.com/go-playground/validator) | v10.30.2 | MIT | Request struct validation |
| [caarlos0/env/v11](https://github.com/caarlos0/env) | v11.4.0 | MIT | Environment variable config parsing |
| [joho/godotenv](https://github.com/joho/godotenv) | v1.5.1 | MIT | `.env` file loading |
| [skip2/go-qrcode](https://github.com/skip2/go-qrcode) | v0.0.0-20200617195104 | MIT | QR code generation for peer configs |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | v0.50.0 | BSD-3-Clause | Cryptographic utilities |
| [google/uuid](https://github.com/google/uuid) | v1.6.0 | BSD-3-Clause | UUID generation |
| [mdlayher/netlink](https://github.com/mdlayher/netlink) | v1.11.1 | MIT | Low-level Linux netlink (transitive) |
| [golang.zx2c4.com/wireguard](https://github.com/WireGuard/wireguard-go) | v0.0.0-20250521234502 | MIT | WireGuard userspace (transitive) |

---

## Frontend (JavaScript / TypeScript)

| Library | Version | License | Purpose |
|---|---|---|---|
| [nuxt](https://github.com/nuxt/nuxt) | ^4.4.2 | MIT | Full-stack Vue framework (used in SPA/static mode) |
| [vue](https://github.com/vuejs/core) | ^3.5.30 | MIT | UI component framework |
| [vue-router](https://github.com/vuejs/router) | ^5.0.3 | MIT | Client-side routing |
| [vuetify](https://github.com/vuetifyjs/vuetify) | ^4.0.5 | MIT | Material Design component library |
| [vuetify-nuxt-module](https://github.com/vuetifyjs/vuetify-loader) | ^0.19.5 | MIT | Vuetify integration for Nuxt |
| [axios](https://github.com/axios/axios) | ^1.13.1 | MIT | HTTP client for API requests |
| [@mdi/font](https://github.com/Templarian/MaterialDesign-Webfont) | ^7.4.47 | Apache 2.0 | Material Design Icons webfont |

---

The complete list of transitive Go dependencies and their exact versions is available in [`backend/go.mod`](backend/go.mod) and [`backend/go.sum`](backend/go.sum).
The complete list of transitive frontend dependencies is available in [`frontend/bun.lock`](frontend/bun.lock).
