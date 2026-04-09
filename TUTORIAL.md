# 📚 Tutorial Lengkap AssistantWhatsapp

Panduan step-by-step untuk memasang dan menggunakan AssistantWhatsapp - Bot WhatsApp dengan fitur Keuangan Pribadi dan Sistem Penjualan.

---

## 📋 Daftar Isi

1. [Apa itu AssistantWhatsapp?](#1-apa-itu-assistantwhatsapp)
2. [Persiapan Sebelum Memulai](#2-persiapan-sebelum-memulai)
3. [Cara Mendapatkan Google Service Account](#3-cara-mendapatkan-google-service-account)
4. [Cara Membuat Google Spreadsheet](#4-cara-membuat-google-spreadsheet)
5. [Cara Mendapatkan API Key LLM](#5-cara-mendapatkan-api-key-llm)
6. [Cara Install Aplikasi](#6-cara-install-aplikasi)
7. [Cara Konfigurasi Aplikasi](#7-cara-konfigurasi-aplikasi)
8. [Cara Menjalankan Aplikasi](#8-cara-menjalankan-aplikasi)
9. [Cara Menggunakan Fitur](#9-cara-menggunakan-fitur)
10. [Troubleshooting](#10-troubleshooting)
11. [FAQ](#11-faq)

---

## 1. Apa itu AssistantWhatsapp?

AssistantWhatsapp adalah bot WhatsApp pintar yang membantu Anda:

### 💰 Fitur Keuangan Pribadi
- Mencatat pengeluaran dan pemasukan dengan bahasa sehari-hari
- Melihat laporan keuangan harian, mingguan, bulanan
- Mengatur budget per kategori dengan peringatan otomatis
- Menyimpan catatan cepat

### 📦 Fitur Sistem Penjualan
- Database item/barang dengan harga beli
- Database customer dengan terms pembayaran (7/14/30 hari)
- Harga jual custom per customer
- Tracking hutang ke supplier
- Tracking piutang dari customer
- Reminder otomatis setiap hari pukul 08:00

### ⏰ Fitur Reminder
- Pengingat satu kali atau berulang
- Parsing tanggal/waktu natural
- Diingatkan 3x sehari sampai selesai

---

## 2. Persiapan Sebelum Memulai

### Yang Anda Butuhkan:

| No | Kebutuhan | Keterangan |
|----|-----------|------------|
| 1 | Komputer/Laptop | Windows, Mac, atau Linux |
| 2 | Koneksi Internet | Stabil untuk proses instalasi |
| 3 | Akun Google | Untuk Google Sheets API |
| 4 | Akun OpenAI/LLM | Untuk fitur AI (OpenAI, OpenRouter, dll) |
| 5 | WhatsApp | Di HP yang sama dengan komputer untuk scan QR |
| 6 | Text Editor | VS Code, Notepad++, atau lainnya |

### Software yang Harus Diinstall:

#### A. Install Go (Golang)

**Windows:**
1. Buka https://go.dev/dl/
2. Download file `goX.XX.X.windows-amd64.msi` (pilih versi terbaru)
3. Double-click file yang didownload
4. Klik "Next" sampai selesai
5. Buka Command Prompt, ketik: `go version`
6. Jika muncul versi Go, berarti berhasil

**Mac:**
1. Buka Terminal
2. Ketik: `brew install go`
3. Tunggu sampai selesai
4. Ketik: `go version`

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install golang-go
go version
```

#### B. Install Git

**Windows:**
1. Buka https://git-scm.com/download/win
2. Download dan install seperti biasa

**Mac:**
```bash
brew install git
```

**Linux:**
```bash
sudo apt install git
```

---

## 3. Cara Mendapatkan Google Service Account

Service Account diperlukan agar bot bisa membaca dan menulis ke Google Sheets.

### Langkah-langkah:

#### Step 1: Buka Google Cloud Console
1. Buka browser, akses: https://console.cloud.google.com/
2. Login dengan akun Google Anda

#### Step 2: Buat Project Baru
1. Klik dropdown di header bagian atas (dekat logo Google Cloud)
2. Klik "NEW PROJECT"
3. Isi:
   - **Project name:** `AssistantWhatsapp` (atau nama lain)
   - **Organization:** Biarkan default
4. Klik "CREATE"
5. Tunggu beberapa detik sampai project selesai dibuat
6. Klik "SELECT PROJECT" untuk masuk ke project

#### Step 3: Enable Google Sheets API
1. Di sidebar kiri, klik "APIs & Services" > "Library"
2. Di kolom pencarian, ketik: `Google Sheets API`
3. Klik hasil "Google Sheets API"
4. Klik tombol "ENABLE"
5. Tunggu sampai enabled

#### Step 4: Buat Service Account
1. Di sidebar kiri, klik "APIs & Services" > "Credentials"
2. Klik "+ CREATE CREDENTIALS" di atas
3. Pilih "Service account"
4. Isi:
   - **Service account name:** `assistant-bot` (atau nama lain)
   - **Service account ID:** Akan otomatis terisi
   - **Description:** `Service account for WhatsApp bot`
5. Klik "CONTINUE"
6. Role: Pilih "Editor" (atau "Owner" untuk akses penuh)
7. Klik "CONTINUE"
8. Klik "DONE"

#### Step 5: Buat dan Download Key JSON
1. Klik Service Account yang baru dibuat
2. Klik tab "KEYS"
3. Klik "ADD KEY" > "Create new key"
4. Pilih "JSON"
5. Klik "CREATE"
6. File JSON akan otomatis terdownload
7. **PENTING:** Simpan file ini dengan aman! Rename menjadi `credentials.json`

---

## 4. Cara Membuat Google Spreadsheet

### Langkah-langkah:

#### Step 1: Buat Spreadsheet Baru
1. Buka https://sheets.google.com/
2. Klik "+" (Blank) untuk buat spreadsheet baru
3. Rename spreadsheet menjadi "AssistantWhatsapp Data"

#### Step 2: Dapatkan Spreadsheet ID
1. Lihat URL spreadsheet Anda, contoh:
   ```
   https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjGMUUqpt35/edit
   ```
2. Spreadsheet ID adalah bagian antara `/d/` dan `/edit`:
   ```
   1BxiMVs0XRA5nFMdKvBdBZjGMUUqpt35
   ```
3. **Salin dan simpan ID ini!**

#### Step 3: Share ke Service Account
1. Di spreadsheet, klik tombol "Share" (bagikan) di kanan atas
2. Paste email Service Account Anda
   - Email format: `nama-service-account@project-id.iam.gserviceaccount.com`
   - Bisa ditemukan di file `credentials.json` bagian `client_email`
3. Pilih role "Editor"
4. Uncheck "Notify people"
5. Klik "Send"

---

## 5. Cara Mendapatkan API Key LLM

### Opsi A: Menggunakan OpenAI

1. Buka https://platform.openai.com/
2. Daftar atau login
3. Klik "API keys" di sidebar
4. Klik "+ Create new secret key"
5. Beri nama, lalu klik "Create"
6. **Salin API key** (hanya muncul sekali!)
7. Simpan dengan aman

### Opsi B: Menggunakan OpenRouter (Lebih Murah)

1. Buka https://openrouter.ai/
2. Daftar atau login
3. Klik "Keys" di sidebar
4. Klik "Create Key"
5. Beri nama, lalu klik "Create"
6. **Salin API key**

**Konfigurasi OpenRouter:**
- Base URL: `https://openrouter.ai/api/v1`
- Model: `openai/gpt-4o-mini` (atau model lain)

### Opsi C: Menggunakan LLM Lokal (Gratis)

Jika punya komputer powerful, bisa gunakan Ollama:
1. Install Ollama dari https://ollama.ai/
2. Jalankan: `ollama run llama3`
3. Base URL: `http://localhost:11434/v1`
4. Model: `llama3`
5. API Key: `ollama` (dummy)

---

## 6. Cara Install Aplikasi

### Step 1: Download Kode

Buka Terminal/Command Prompt, jalankan:

```bash
# Pindah ke folder yang diinginkan
cd ~
# Atau di Windows: cd %USERPROFILE%

# Clone repository
git clone https://github.com/winolikemove/AssistantWhatsapp.git

# Masuk ke folder
cd AssistantWhatsapp
```

### Step 2: Pindahkan credentials.json

1. Copy file `credentials.json` yang sudah didownload
2. Paste ke dalam folder `AssistantWhatsapp`

### Step 3: Download Dependencies

```bash
go mod tidy
```

Tunggu sampai proses selesai (perlu koneksi internet).

---

## 7. Cara Konfigurasi Aplikasi

### Step 1: Buat File .env

Buat file baru bernama `.env` di folder `AssistantWhatsapp`:

**Windows (Command Prompt):**
```cmd
copy .env.example .env
```

**Mac/Linux:**
```bash
cp .env.example .env
```

### Step 2: Edit File .env

Buka file `.env` dengan text editor, isi dengan data Anda:

```env
# ===========================================
# KONFIGURASI LLM (WAJIB)
# ===========================================
# Untuk OpenAI:
LLM_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxx
LLM_BASE_URL=https://api.openai.com/v1
LLM_MODEL=gpt-4o-mini

# Untuk OpenRouter:
# LLM_API_KEY=sk-or-xxxxxxxxxxxxxxxxxxxxxx
# LLM_BASE_URL=https://openrouter.ai/api/v1
# LLM_MODEL=openai/gpt-4o-mini

# ===========================================
# KONFIGURASI GOOGLE SHEETS (WAJIB)
# ===========================================
GOOGLE_APPLICATION_CREDENTIALS=./credentials.json
SHEETS_SPREADSHEET_ID=1BxiMVs0XRA5nFMdKvBdBZjGMUUqpt35

# ===========================================
# KONFIGURASI WHATSAPP (WAJIB)
# ===========================================
WHATSAPP_SESSION_DB_PATH=./session.db
OWNER_PHONE_NUMBER=6281234567890

# ===========================================
# KONFIGURASI PENJUALAN (OPSIONAL)
# ===========================================
SUPPLIER_NAME=Toko Supplier
SUPPLIER_PAY_DAY=25
REMINDER_TIME=08:00
WA_REMINDER_TO_CUSTOMER=false
```

### Penjelasan Konfigurasi:

| Variable | Keterangan | Contoh |
|----------|------------|--------|
| `LLM_API_KEY` | API key dari OpenAI/OpenRouter | `sk-proj-xxx` |
| `LLM_BASE_URL` | URL API LLM | `https://api.openai.com/v1` |
| `LLM_MODEL` | Model yang digunakan | `gpt-4o-mini` |
| `GOOGLE_APPLICATION_CREDENTIALS` | Path ke file credentials.json | `./credentials.json` |
| `SHEETS_SPREADSHEET_ID` | ID Google Spreadsheet | `1BxiMV...` |
| `WHATSAPP_SESSION_DB_PATH` | Path untuk menyimpan session | `./session.db` |
| `OWNER_PHONE_NUMBER` | Nomor WhatsApp Anda (tanpa +) | `6281234567890` |
| `SUPPLIER_NAME` | Nama supplier Anda | `Toko ABC` |
| `SUPPLIER_PAY_DAY` | Tanggal bayar supplier | `25` |
| `REMINDER_TIME` | Waktu reminder harian | `08:00` |
| `WA_REMINDER_TO_CUSTOMER` | Kirim reminder ke customer? | `true` atau `false` |

---

## 8. Cara Menjalankan Aplikasi

### Step 1: Jalankan Bot

```bash
go run ./cmd/bot
```

### Step 2: Scan QR Code

1. Buka WhatsApp di HP
2. Ketuk tiga titik (⋮) di kanan atas
3. Pilih "Linked devices"
4. Ketuk "Link a device"
5. Scan QR Code yang muncul di terminal/komputer
6. Tunggu sampai muncul pesan "✅ Bot is running!"

### Step 3: Test Bot

Kirim pesan dari WhatsApp ke nomor sendiri:
```
halo
```

Bot akan membalas dengan sapaan.

### Proses di Background (Opsional)

**Windows:**
```bash
# Build executable
go build -o assistant-bot.exe ./cmd/bot

# Jalankan
./assistant-bot.exe
```

**Mac/Linux:**
```bash
# Build executable
go build -o assistant-bot ./cmd/bot

# Jalankan di background
nohup ./assistant-bot &
```

---

## 9. Cara Menggunakan Fitur

### 💰 Fitur Keuangan

#### Mencatat Pengeluaran
Ketik dengan bahasa sehari-hari:
```
beli makan 25k
bayar listrik 150rb
isi bensin 100k
belanja di indomaret 50rb
```

#### Mencatat Pemasukan
```
gaji 5jt
terima transfer 500rb
freelance 2jt
```

#### Melihat Laporan
```
laporan
laporan hari ini
laporan minggu ini
laporan bulan ini
```

#### Mengatur Budget
```
set budget makanan 1jt
budget transportasi 500rb
```

#### Edit/Hapus Transaksi
```
/edit 20260317-001 jumlah 30000
/hapus 20260317-001
```

### 📦 Fitur Penjualan

#### Menambah Item
```
tambah item aussie bbq 15000 kg
tambah item dori fillet 25000 kg
```

#### Menambah Customer
```
tambah customer ambrogio jakarta 14 hari credit
tambah customer toko budi bandung 7 hari credit
tambah customer warung siti cash
```

#### Set Harga Jual
```
set harga aussie bbq untuk ambrogio 18000
set harga dori fillet untuk toko budi 30000
```

#### Mencatat Penjualan
```
jual aussie bbq 50 kg ke ambrogio
jual dori fillet 20 kg ke toko budi
```

#### Laporan Penjualan
```
laporan profit
laporan profit bulan ini
```

#### Lihat Piutang
```
piutang
siapa yang hutang
```

#### Lihat Hutang ke Supplier
```
hutang
hutang supplier
```

#### Catat Pembayaran
```
bayar piutang dari ambrogio 900000
bayar hutang supplier 500000
```

### ⏰ Fitur Reminder

#### Buat Reminder
```
ingetin tanggal 25 maret bayar listrik
ingatkan jam 9 meeting dengan client
reminder besok pagi cek email
```

#### Tandai Selesai
```
/done RMD-20260317-090000-001
```

Atau natural:
```
udah bayar listrik
sudah meeting dengan client
```

### 📝 Catatan Cepat
```
/notes beli kado ultah minggu depan
/notes ide usaha baru: jual kopi
```

### 📋 Menu Lengkap

Kirim `/help` atau `/menu` untuk melihat semua command:

| Command | Fungsi |
|---------|--------|
| `/start` | Mulai bot |
| `/help` | Bantuan |
| `/kategori` | Daftar kategori |
| `/laporan [periode]` | Laporan keuangan |
| `/budget [kategori] [jumlah]` | Set budget |
| `/notes [teks]` | Simpan catatan |
| `/reminder [teks]` | Buat pengingat |
| `/done [id]` | Tandai selesai |
| `/edit [id] [field] [nilai]` | Edit transaksi |
| `/hapus [id]` | Hapus transaksi |
| `/export` | Link Google Sheets |

---

## 10. Troubleshooting

### Masalah: QR Code Tidak Muncul
**Solusi:**
1. Pastikan terminal mendukung tampilan QR
2. Coba jalankan ulang
3. Gunakan terminal yang berbeda (cmd, PowerShell, Git Bash)

### Masalah: "Config error"
**Solusi:**
1. Pastikan file `.env` ada di folder yang sama
2. Cek format `.env` tidak ada spasi di sekitar `=`
3. Pastikan tidak ada karakter tersembunyi

### Masalah: "Sheets error"
**Solusi:**
1. Cek `credentials.json` sudah benar
2. Cek Spreadsheet ID sudah benar
3. Cek Service Account sudah di-share ke spreadsheet
4. Cek Google Sheets API sudah di-enable

### Masalah: "LLM connectivity check failed"
**Solusi:**
1. Cek API key sudah benar
2. Cek Base URL sudah benar
3. Cek koneksi internet
4. Cek saldo API (untuk OpenAI/OpenRouter)

### Masalah: "WhatsApp connection error"
**Solusi:**
1. Hapus file `session.db`
2. Jalankan ulang bot
3. Scan QR ulang
4. Pastikan WhatsApp di HP tidak sedang digunakan di perangkat lain

### Masalah: Bot Tidak Membalas
**Solusi:**
1. Cek nomor di `OWNER_PHONE_NUMBER` sudah benar
2. Pastikan bot masih berjalan
3. Cek log di terminal untuk error

### Masalah: Tab di Google Sheets Tidak Muncul
**Solusi:**
1. Bot akan otomatis membuat tab saat pertama kali dijalankan
2. Coba kirim pesan untuk trigger pembuatan tab
3. Cek permission spreadsheet

---

## 11. FAQ

### Q: Apakah bot ini berbayar?
**A:** Bot ini gratis. Namun, untuk fitur AI, Anda perlu API key dari OpenAI atau provider LLM yang mungkin berbayar. Alternatif gratis: gunakan Ollama dengan model lokal.

### Q: Apakah data saya aman?
**A:** Data disimpan di Google Sheets milik Anda. Kredensial disimpan lokal di komputer Anda. Tidak ada data yang dikirim ke server pihak ketiga selain API LLM.

### Q: Bisa digunakan di HP?
**A:** Bot ini dirancang untuk dijalankan di komputer. Untuk penggunaan di HP, diperlukan setup tambahan seperti termux (Android) atau server cloud.

### Q: Bagaimana cara backup data?
**A:** Data tersimpan di Google Sheets. Anda bisa export kapan saja dari Google Sheets (File > Download).

### Q: Bisa multi-user?
**A:** Saat ini bot dirancang untuk single-user. Hanya nomor di `OWNER_PHONE_NUMBER` yang bisa mengakses semua fitur.

### Q: Bagaimana cara update?
**A:**
```bash
cd AssistantWhatsapp
git pull
go mod tidy
```

### Q: Bagaimana cara berhenti menggunakan?
**A:**
1. Tekan `Ctrl+C` di terminal untuk stop bot
2. Hapus folder `AssistantWhatsapp`
3. (Opsional) Hapus project di Google Cloud Console
4. (Opsional) Hapus spreadsheet

### Q: Error "module not found"
**A:** Jalankan `go mod tidy` lagi, pastikan koneksi internet stabil.

### Q: Bisa pakai nomor berbeda untuk bot?
**A:** Ya, bot akan menggunakan nomor WhatsApp yang di-scan via QR code, bukan nomor di config. Config `OWNER_PHONE_NUMBER` adalah nomor yang boleh mengirim perintah ke bot.

---

## 📞 Bantuan

Jika mengalami masalah:

1. Baca ulang dokumentasi
2. Cek bagian [Troubleshooting](#10-troubleshooting)
3. Buka issue di GitHub: https://github.com/winolikemove/AssistantWhatsapp/issues

---

## 📄 Lisensi

MIT License - Bebas digunakan dan dimodifikasi.

---

**Selamat menggunakan AssistantWhatsapp! 🎉**
