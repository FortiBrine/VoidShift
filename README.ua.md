<div align="right">

🌐 [English](README.md) | **Українська**

</div>

<div align="center">

# VoidShift

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker&logoColor=white)](compose.yml)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20FreeBSD-orange)

Веб-панель керування VPN для VPS-серверів. Єдиний бінарник з вбудованим веб-інтерфейсом.
Наразі підтримується WireGuard, у майбутньому планується підтримка інших протоколів.

</div>

---

### ✨ Можливості

- 🔒 **Керування WireGuard** — створення мереж, підняття/зупинка інтерфейсів, додавання та видалення пірів
- 📱 **Конфіги та QR-коди для пірів** — завантаження `.conf`-файлів або сканування QR-коду прямо з інтерфейсу
- 📦 **Єдиний бінарник** — фронтенд вбудований у Go-бінарник на етапі збірки, один файл для деплою
- 🗄️ **Дві бази даних** — SQLite з коробки (без налаштувань) або MySQL через DSN
- 🔑 **Авторизація на основі сесій** — початкові дані адміна задаються через змінні середовища
- 🐳 **Готово до Docker** — `docker compose up` і все запущено

---

### 🖥️ Підтримувані операційні системи

- ✅ **Linux** — повна підтримка
- 🔜 **FreeBSD** — планується
- ❌ **macOS** — не підтримується (WireGuard kernel API недоступні)
- ❌ **Windows** — не підтримується (WireGuard kernel API недоступні)

---

### 🚀 Швидкий старт

#### Docker (рекомендовано)

```sh
cp .env.example .env
# Відредагуйте .env — мінімум: HOST_ADDRESS, ADMIN_USERNAME, ADMIN_PASSWORD
docker compose up -d
```

Відкрийте [http://localhost:8080](http://localhost:8080) і увійдіть з даними адміна.

> Всі необхідні можливості (`NET_ADMIN`, `/dev/net/tun`, `net.ipv4.ip_forward`) вже налаштовані у `compose.yml`.

#### Ручна збірка

```sh
cp .env.example .env
# Відредагуйте .env

make deps     # встановити залежності фронтенду та бекенду
make build    # зібрати фронтенд, потім вбудувати його в Go-бінарник

sudo ./app    # потрібні привілеї NET_ADMIN для WireGuard
```

<details>
<summary>Інші команди make</summary>

| Команда | Опис |
|---|---|
| `make deps` | `bun install` + `go mod download` |
| `make frontend` | Генерація статичного SPA у `frontend/.output/public/` |
| `make backend` | Копіювання фронтенду у embed-директорію, компіляція Go-бінарника |
| `make build` | Повна збірка: frontend + backend |
| `make clean` | Видалення артефактів збірки |

</details>

---

### ⚙️ Конфігурація

Вся конфігурація через змінні середовища. Скопіюйте `.env.example` у `.env`:

| Змінна | За замовч. | Обов'язкова | Опис |
|---|---|---|---|
| `SQLITE_DATABASE_PATH` | | Одна з двох | Шлях до SQLite-файлу (наприклад `./database.db`) |
| `MYSQL_DSN` | | Одна з двох | MySQL DSN; якщо задано, використовується MySQL |
| `HOST_ADDRESS` | | Так | Публічна IP-адреса сервера, використовується у конфігах пірів |
| `ADMIN_USERNAME` | | Так | Ім'я початкового адміна |
| `ADMIN_PASSWORD` | | Так | Пароль початкового адміна |
| `HTTP_ADDRESS` | `:8080` | Ні | Адреса HTTP-сервера |
| `GRACEFUL_TIMEOUT` | `5s` | Ні | Таймаут graceful shutdown |
| `ENVIRONMENT` | | Ні | Встановіть `dev` для розробницького логування |

---

### 📡 API

Усі маршрути знаходяться під `/api`. Захищені маршрути вимагають активну cookie-сесію.

| Метод | Шлях | Авторизація | Опис |
|---|---|---|---|
| `GET` | `/api/health` | Ні | Перевірка стану сервісу |
| `POST` | `/api/auth/login` | Ні | Вхід |
| `POST` | `/api/auth/logout` | Ні | Вихід |
| `GET` | `/api/vpn/wireguard/networks` | Так | Список WireGuard-мереж |
| `POST` | `/api/vpn/wireguard/networks/generate` | Так | Створити нову мережу |
| `GET` | `/api/vpn/wireguard/networks/:id` | Так | Деталі мережі |
| `POST` | `/api/vpn/wireguard/networks/:id/up` | Так | Підняти інтерфейс |
| `POST` | `/api/vpn/wireguard/networks/:id/down` | Так | Зупинити інтерфейс |
| `DELETE` | `/api/vpn/wireguard/networks/:id` | Так | Видалити мережу |
| `POST` | `/api/vpn/wireguard/networks/:id/peers/generate` | Так | Додати пір до мережі |
| `GET` | `/api/vpn/wireguard/peers/:peerId/config` | Так | Отримати конфіг піра |
| `GET` | `/api/vpn/wireguard/peers/:peerId/config/download` | Так | Завантажити `.conf`-файл піра |
| `GET` | `/api/vpn/wireguard/peers/:peerId/qr` | Так | Отримати конфіг піра у вигляді QR-коду |
| `DELETE` | `/api/vpn/wireguard/peers/:peerId` | Так | Видалити пір |

---

### 📦 Структура проєкту

```
backend/               # Корінь Go-модуля (go.mod знаходиться тут)
  cmd/api/             # Точка входу бінарника
  internal/
    app/               # Збірка компонентів застосунку
    auth/              # Обробник, сервіс та middleware авторизації
    session/           # Керування сесіями (GORM, TTL 5 днів)
    user/              # Модель, репозиторій та сервіс користувача
    wireguard/         # Мережі та піри WireGuard
    shared/            # Конфігурація, БД, HTTP, логер, валідатор
    embed/webui/       # Вбудований фронтенд (генерується автоматично, не редагувати)
frontend/              # Nuxt 4 SPA (SSR вимкнено, збірка через Bun)
  app/pages/           # Сторінки маршрутів
  i18n/                # Файли локалізації
```

---

### 📄 Ліцензія

VoidShift розповсюджується під ліцензією [Apache License 2.0](LICENSE).
Повний список сторонніх залежностей: [THIRD_PARTY_LIBRARIES.md](THIRD_PARTY_LIBRARIES.md).
