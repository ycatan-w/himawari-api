# 🌻 Himawari API

Himawari API is the **backend server** of the Himawari project — a personal **journal, planner, and daily overview tool**.
It provides a RESTful API to manage **users, events, and logs**, with authentication and secure endpoints.

---

## ✨ Features

- 🔐 **Authentication**: login & register
- 👤 **User management**
- 📅 **Events**: create, read, update, delete (CRUD)
- 📓 **Logs**: create, read, update, delete (CRUD)
- 🛡️ **Middleware** for authentication & security
- 📦 **Standardized JSON responses**

---

## 🌐 Demo (Web)

You can try a live demo here:
👉 [https://ycatan-w.github.io/himawari-web/](https://ycatan-w.github.io/himawari-web/)

**Demo credentials**:

- **Username**: `demo`
- **Password**: `demo`

---

## 📂 Project Structure

- `cmd/server/` → Application entrypoint
- `internal/server/` → Server configuration & startup
- `internal/db/` → Database logic (SQLite)
- `internal/api/` → API routes

  - `auth.go` → Authentication endpoints
  - `event.go` → Event endpoints
  - `log.go` → Log endpoints
  - `middleware/` → Security middleware
  - `internal/` → Specialized modules (login, register, utils, log CRUD, event CRUD)

---

## 🚀 Installation

There are two main installation paths: **development setup** and **end-user installation**.

### 🔧 Development Setup

#### Prerequisites

- Go 1.20+
- SQLite3

#### Steps

```bash
# Clone the repository
git clone https://github.com:ycatan-w/himawari-api.git
cd himawari-api

# Install dependencies
go mod tidy

# Build the server
go build -o himawari ./cmd/server

# Run the server
./himawari
```

#### Makefile Commands

- `make build` → build the project
- `make run` → run the server

---

### 💻 Quick User Installation Guide

#### Windows

1. Download the correct `.zip` file (`386` or `amd64`) and extract it.
2. Run `install.bat` or `install.ps1`.
   - The script will automatically request admin rights to complete the installation.
3. Once installed, open your browser: [http://localhost:9740/web](http://localhost:9740/web)

**Uninstall:** Run `uninstall.bat` or `uninstall.ps1` from the install folder.

💡 For full instructions, troubleshooting tips, and checking the service status, see: [WINDOWS.md](WINDOWS.md)

---

#### macOS

1. Download the appropriate `.pkg` (or `.tar.gz`) for your system (`amd64`, `arm64`, or universal).
2. `.pkg`: Double-click and follow the installer (macOS will prompt for admin credentials).
   `.tar.gz`: Extract and run `sudo ./install.sh` from Terminal.
3. Once installed, open your browser: [http://localhost:9740/web](http://localhost:9740/web)

**Uninstall:** Run `sudo himawari-server-uninstall`.

💡 For full instructions, tips, and checking the LaunchDaemon, see: [MACOS.md](MACOS.md)

---

#### Linux

1. Download the appropriate `.tar.gz` for your system (`386`, `amd64`, or `arm64`) and extract it.
2. Run the installer: `sudo ./install.sh`
3. Once installed, open your browser: [http://localhost:9740/web](http://localhost:9740/web)

**Uninstall:** Run `sudo himawari-server-uninstall`.

💡 For full instructions, tips, and checking the systemd service, see: [LINUX.md](LINUX.md)

---

#### Notes

- Windows scripts relaunch automatically with admin rights; no manual admin action is required.
- macOS/Linux scripts will stop if not run with `sudo` and will suggest the correct command.
- After uninstalling, some helper scripts may need manual removal.
- Always verify that the service/daemon has been stopped after uninstall (Windows Services, LaunchDaemons on macOS, systemd on Linux).

---

## 📖 API Documentation

OpenAPI (Swagger/Redoc) documentation will be added soon and published on GitHub Pages.

- OpenAPI spec: [`swagger.json`](./docs/swagger.json)
- Redoc page: [GitHub Pages link](https://ycatan-w.github.io/himawari-api/)

---

## 💡 Usage Examples

### Login

```bash
curl -X POST http://localhost:9740/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "password"}'
```

Response:

```json
{
  "token": "xxxxx.yyyyy.zzzzz"
}
```

### Create Event

```bash
curl -X POST http://localhost:9740/api/events \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Meeting", "date": "2025-01-01"}'
```

---

## 🔗 Himawari Ecosystem

This API is part of the **Himawari ecosystem**:

- **Web**: main UI (frontend client) → [Himawari Web](https://github.com/ycatan-w/himawari-web)
- **Backend API**: this repository
- **Desktop**: planned cross-platform desktop client
