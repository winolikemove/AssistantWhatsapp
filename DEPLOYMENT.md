# Panduan Deployment - WhatsApp AI Assistant

## 📱 Opsi 1: Menjalankan di HP Android dengan Termux

### Kelebihan:
- ✅ Gratis 100%
- ✅ Tidak perlu PC/Laptop
- ✅ Bisa berjalan di background

### Kekurangan:
- ❌ Harus tetap menyala HP
- ❌ Koneksi internet harus stabil
- ❌ Performa terbatas

### Langkah-langkah:

#### 1. Install Termux
Download dari F-Droid (JANGAN dari Play Store - versi lama):
```
https://f-droid.org/packages/com.termux/
```

#### 2. Install Dependencies di Termux
Buka Termux, jalankan:

```bash
# Update package
pkg update && pkg upgrade -y

# Install Git
pkg install git -y

# Install Go
pkg install golang -y

# Install necessary tools
pkg install wget curl -y
```

#### 3. Clone & Setup Project
```bash
# Clone repository
git clone https://github.com/winolikemove/AssistantWhatsapp.git
cd AssistantWhatsapp

# Setup Google Sheets credentials
# Buat folder credentials
mkdir -p credentials

# Copy credentials.json ke folder tersebut
# Cara: Buka file manager, copy credentials.json ke:
# /data/data/com.termux/files/home/AssistantWhatsapp/credentials/
```

#### 4. Konfigurasi
```bash
# Copy .env.example ke .env
cp .env.example .env

# Edit konfigurasi
nano .env
```

Isi dengan:
```env
OWNER_NUMBER=6281234567890
LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-api-key
OPENAI_MODEL=gpt-4o-mini
```

#### 5. Download Dependencies & Run
```bash
# Download dependencies
go mod tidy

# Jalankan aplikasi
go run ./cmd/bot
```

#### 6. Scan QR Code
- QR code akan muncul di Termux
- Buka WhatsApp di HP lain → Menu → Perangkat tertaut → Tautkan perangkat
- Scan QR code yang muncul

### Tips Agar Tetap Berjalan:
```bash
# Install termux-wake-lock agar tidak sleep
pkg install termux-api

# Atau gunakan screen untuk background
pkg install screen
screen -S bot
go run ./cmd/bot

# Tekan Ctrl+A+D untuk detach (tetap jalan di background)
# Untuk kembali: screen -r bot
```

---

## ☁️ Opsi 2: Server Cloud Gratis (RECOMMENDED)

### A. Railway.app (Paling Mudah) ⭐⭐⭐⭐⭐

**Kelebihan:**
- ✅ Gratis $5 credit/bulan
- ✅ Deploy dari GitHub otomatis
- ✅ Mudah dikonfigurasi
- ✅ Tidak perlu kartu kredit

**Kekurangan:**
- ❌ Credit habis jika intensif

**Langkah-langkah:**

1. **Fork repository** ke GitHub Anda

2. **Buat akun Railway** di [railway.app](https://railway.app)

3. **New Project → Deploy from GitHub repo**

4. **Set Environment Variables:**
```
OWNER_NUMBER=6281234567890
LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-api-key
OPENAI_MODEL=gpt-4o-mini
```

5. **Upload credentials.json:**
   - Railway tidak mendukung file upload langsung
   - Solusi: Encode ke base64 dan simpan di env var

```bash
# Di komputer lokal, encode credentials
base64 -i credentials.json

# Copy hasilnya ke env var:
GOOGLE_CREDENTIALS_BASE64=<hasil-base64>
```

6. **Modifikasi main.go** untuk membaca dari env:
```go
// Tambahkan kode ini di awal main()
if credsBase64 := os.Getenv("GOOGLE_CREDENTIALS_BASE64"); credsBase64 != "" {
    credsJSON, err := base64.StdEncoding.DecodeString(credsBase64)
    if err != nil {
        log.Fatalf("Failed to decode credentials: %v", err)
    }
    err = os.WriteFile("credentials/credentials.json", credsJSON, 0644)
    if err != nil {
        log.Fatalf("Failed to write credentials: %v", err)
    }
}
```

### B. Render.com ⭐⭐⭐⭐

**Kelebihan:**
- ✅ Gratis 750 jam/bulan
- ✅ Mudah setup
- ✅ Auto-deploy dari GitHub

**Kekurangan:**
- ❌ Service sleep setelah 15 menit idle
- ❌ Perlu kartu kredit untuk verifikasi

**Langkah:**

1. Buat akun di [render.com](https://render.com)
2. New → Web Service → Connect GitHub
3. Pilih repository
4. Set Build Command: `go build -o bot ./cmd/bot`
5. Set Start Command: `./bot`
6. Add Environment Variables

### C. Fly.io ⭐⭐⭐⭐

**Kelebihan:**
- ✅ Gratis 3 VM small
- ✅ Performa bagus
- ✅ Tidak sleep

**Kekurangan:**
- ❌ Perlu kartu kredit
- ❌ Setup lebih teknis

**Langkah:**

1. Install flyctl:
```bash
curl -L https://fly.io/install.sh | sh
```

2. Login:
```bash
fly auth login
```

3. Buat `fly.toml` di root project:
```toml
app = "assistant-whatsapp"
primary_region = "sin"

[build]
  builder = "paketobuildpacks/builder:base"

[env]
  PORT = "8080"

[mounts]
  source = "data"
  destination = "/app/data"
```

4. Deploy:
```bash
fly launch
fly secrets set OWNER_NUMBER=6281234567890
fly secrets set OPENAI_API_KEY=sk-your-key
```

### D. Oracle Cloud (Always Free) ⭐⭐⭐⭐⭐

**Kelebihan:**
- ✅ Gratis selamanya (Always Free)
- ✅ VM ARM 4 OCPU + 24GB RAM
- ✅ Tidak sleep
- ✅ Sangat powerful

**Kekurangan:**
- ❌ Setup kompleks
- ❌ Perlu kartu kredit verifikasi
- ❌ Region terbatas

**Langkah:**

1. Daftar Oracle Cloud Free Tier
2. Create Compute Instance (ARM)
3. Pilih OS: Ubuntu 22.04
4. SSH ke instance
5. Install Go & jalankan aplikasi

### E. Koyeb ⭐⭐⭐

**Kelebihan:**
- ✅ Gratis $5.50 credit/bulan
- ✅ Mudah setup
- ✅ Global CDN

**Langkah:**

1. Daftar di [koyeb.com](https://www.koyeb.com)
2. Create App → GitHub
3. Set environment variables
4. Deploy

---

## 🔧 Konfigurasi Tambahan untuk Cloud

### Persistent Volume untuk Session

WhatsApp session disimpan di file SQLite. Di cloud, perlu persistent storage:

**Untuk Railway:**
```toml
# railway.toml
[build]
builder = "heroku/buildpacks:20"

[[services]]
  internal_port = 8080
  protocol = "tcp"
  
  [[services.volume]]
  name = "session-data"
  mount_path = "/app/session"
```

**Untuk Docker:**
Buat `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bot ./cmd/bot

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bot .
COPY --from=builder /app/credentials ./credentials
VOLUME ["/app/session", "/app/credentials"]
CMD ["./bot"]
```

### Environment Variables Lengkap:

```env
# Owner WhatsApp number (tanpa +)
OWNER_NUMBER=6281234567890

# Google Sheets ID
GOOGLE_SHEETS_ID=1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms

# LLM Configuration
LLM_PROVIDER=openai
OPENAI_API_KEY=sk-your-api-key
OPENAI_MODEL=gpt-4o-mini
# Atau gunakan OpenRouter:
# LLM_PROVIDER=openrouter
# OPENROUTER_API_KEY=sk-or-your-key

# Atau gunakan Ollama (local):
# LLM_PROVIDER=ollama
# OLLAMA_BASE_URL=http://localhost:11434
# OLLAMA_MODEL=llama3

# Google Credentials (untuk cloud)
GOOGLE_CREDENTIALS_BASE64=<base64-encoded-json>
```

---

## 📋 Rekomendasi

| Kebutuhan | Rekomendasi |
|-----------|-------------|
| Coba-coba / Development | Termux di HP |
| Production Ringan | Railway.app |
| Production Stabil | Oracle Cloud Free |
| Paling Mudah | Render.com |
| Performa Terbaik | Fly.io |

---

## ❓ FAQ

### Q: Apakah bisa jalan 24/7 di Termux?
A: Bisa, tapi HP harus tetap menyala dan koneksi internet stabil. Tidak disarankan untuk production.

### Q: Apakah QR code harus di-scan ulang?
A: Tidak, selama file session tidak terhapus. Di cloud, gunakan persistent volume.

### Q: Berapa biaya LLM (OpenAI)?
A: Dengan gpt-4o-mini, sekitar $0.15-0.50/bulan untuk penggunaan normal.

### Q: Apakah bisa pakai LLM gratis?
A: Bisa! Gunakan Ollama di server sendiri, atau Groq (gratis dengan rate limit).

### Q: Kenapa service saya sleep di Render?
A: Free tier Render sleep setelah 15 menit tidak ada request. Gunakan UptimeRobot untuk ping berkala, atau upgrade ke paid plan.

---

## 🆘 Butuh Bantuan?

Jika mengalami masalah:
1. Buka Issue di GitHub: https://github.com/winolikemove/AssistantWhatsapp/issues
2. Sertakan log error dan konfigurasi (tanpa API key)
