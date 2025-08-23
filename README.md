# URL Shortener

A simple URL shortening service written in Go, using **SQLite** for storage and **chi router** for routing.  
Supports creating, retrieving (redirect), and deleting short links.  

---

## 🚀 Features
- Create short links (`POST /url`)
- Redirect by alias (`GET /{alias}`)
- Delete short links (`DELETE /{alias}`)
- Structured logging (text or JSON format)
- Middleware:
  - `RequestID` — unique request ID
  - `Recoverer` — prevents crashes on `panic`
  - Custom logger for better tracing

---

## 🛠️ Tech Stack
- Go 1.21+
- [chi](https://github.com/go-chi/chi) — lightweight router
- SQLite — persistent storage
- `log/slog` — structured logging
- YAML-based configuration (different environments: local, dev, prod)

---

## 📦 Installation & Run

1. Clone the repository:

```bash
git clone https://github.com/Ya095/url-shortener.git
cd url-shortener
```