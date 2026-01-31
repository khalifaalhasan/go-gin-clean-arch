# Secure Anonymous Comment System (Backend)

A production-ready backend service for anonymous comments, built with Golang and Clean Architecture. Focuses on high security against Spam, DDoS, XSS, and Bots.

## 🛠 Tech Stack
- **Language:** Golang 1.25+
- **Framework:** Gin Gonic
- **Rate Limiting:** Redis
- **Security:** Cloudflare Turnstile (Anti-Bot), Bluemonday (Anti-XSS)
- **Validation:** go-playground/validator

## 🛡️ Security Features (Defense in Depth)
1. **Rate Limiting:** Max 3 comments/minute per IP (via Redis). Returns `429 Too Many Requests`.
2. **Bot Verification:** Server-side validation of Cloudflare Turnstile tokens.
3. **XSS Sanitization:** Strips dangerous HTML/Scripts from content before processing.
4. **Strict Validation:** Struct-level validation (alphanumeric check, length constraints).

## 🚀 How to Run

### 1. Prerequisites
- Go installed
- Docker (for Redis)

### 2. Start Redis
```bash
docker run --name my-redis -p 6379:6379 -d redis
```

### 3. Setup Environment (Optional)
```bash
export REDIS_ADDR="localhost:6379"
export TURNSTILE_SECRET_KEY="your_cloudflare_secret_key"
```

### 4. Run the Server
```bash
go mod tidy
go run cmd/api/main.go
```

Server will run on http://localhost:8080.

## 📂 Project Structure


```
├── 📁 api
│   ├── 📁 cmd
│   │   └── 📁 api
│   │       └── 🐹 main.go
│   ├── 📁 internal
│   │   ├── 📁 config
│   │   ├── 📁 handler
│   │   │   └── 🐹 comment_handler.go
│   │   ├── 📁 middleware
│   │   │   └── 🐹 ratelimit.go
│   │   ├── 📁 model
│   │   │   └── 🐹 comment.go
│   │   └── 📁 service
│   │       └── 🐹 security_service.go
│   ├── 📄 go.mod
│   └── 📄 go.sum
├── 📁 web
│   ├── 📁 public
│   │   ├── 🖼️ file.svg
│   │   ├── 🖼️ globe.svg
│   │   ├── 🖼️ next.svg
│   │   ├── 🖼️ vercel.svg
│   │   └── 🖼️ window.svg
│   ├── 📁 src
│   │   └── 📁 app
│   │       ├── 📄 favicon.ico
│   │       ├── 🎨 globals.css
│   │       ├── 📄 layout.tsx
│   │       └── 📄 page.tsx
│   ├── ⚙️ .gitignore
│   ├── 📄 eslint.config.mjs
│   ├── 📄 next-env.d.ts
│   ├── 📄 next.config.ts
│   ├── ⚙️ package-lock.json
│   ├── ⚙️ package.json
│   ├── 📄 postcss.config.mjs
│   └── ⚙️ tsconfig.json
└── 📝 README.md
```
