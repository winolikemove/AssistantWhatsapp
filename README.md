# AssistantWhatsapp

WhatsApp AI Assistant dengan fitur keuangan pribadi dan sistem penjualan terintegrasi.

> **📖 [Baca Tutorial Lengkap](TUTORIAL.md)** - Panduan step-by-step untuk pemula

## ✨ Fitur

### 💰 Keuangan Pribadi
- Catat pengeluaran/pemasukan dengan bahasa natural
- Laporan harian, mingguan, bulanan
- Budget per kategori dengan peringatan
- Catatan cepat

### 📦 Sistem Penjualan
- Database item dengan harga beli dari supplier
- Database customer dengan terms pembayaran (7/14/30 hari)
- Harga jual custom per customer
- Tracking hutang ke supplier (jatuh tempo tanggal 25)
- Tracking piutang dari customer
- Reminder otomatis setiap hari pukul 08:00

### ⏰ Reminder
- Pengingat satu kali atau berulang
- Parsing tanggal/waktu natural
- Auto-complete dari chat

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Google Cloud Project dengan Sheets API
- LLM API Key (OpenAI/OpenRouter)

### Install & Run

```bash
# Clone repository
git clone https://github.com/winolikemove/AssistantWhatsapp.git
cd AssistantWhatsapp

# Copy environment template
cp .env.example .env

# Edit .env dengan konfigurasi Anda
# Lihat TUTORIAL.md untuk panduan lengkap

# Download dependencies
go mod tidy

# Run
go run ./cmd/bot
```

### Scan QR Code
Buka WhatsApp di HP → Linked devices → Link a device → Scan QR

## 📱 Contoh Penggunaan

### Keuangan
```
beli makan 25k          # Catat pengeluaran
gaji 5jt               # Catat pemasukan
laporan bulan ini      # Lihat laporan
```

### Penjualan
```
tambah item aussie bbq 15000 kg              # Tambah item
tambah customer ambrogio jakarta 14 hari     # Tambah customer
set harga aussie bbq untuk ambrogio 18000    # Set harga
jual aussie bbq 50 kg ke ambrogio            # Catat penjualan
```

### Laporan
```
laporan profit         # Laporan keuntungan
piutang                # Daftar piutang
hutang                 # Daftar hutang supplier
```

### Reminder
```
ingetin tanggal 25 bayar listrik    # Buat reminder
/done RMD-xxx                       # Tandai selesai
```

## ⚙️ Konfigurasi

| Variable | Required | Description |
|----------|----------|-------------|
| `LLM_API_KEY` | ✅ | API key LLM |
| `LLM_BASE_URL` | ✅ | Base URL API |
| `LLM_MODEL` | ✅ | Model name |
| `GOOGLE_APPLICATION_CREDENTIALS` | ✅ | Path credentials.json |
| `SHEETS_SPREADSHEET_ID` | ✅ | Spreadsheet ID |
| `WHATSAPP_SESSION_DB_PATH` | ✅ | Path session.db |
| `OWNER_PHONE_NUMBER` | ✅ | Nomor WhatsApp owner |
| `SUPPLIER_NAME` | ❌ | Nama supplier |
| `SUPPLIER_PAY_DAY` | ❌ | Tanggal bayar supplier (default: 25) |
| `REMINDER_TIME` | ❌ | Waktu reminder (default: 08:00) |
| `WA_REMINDER_TO_CUSTOMER` | ❌ | Kirim reminder ke customer |

## 📚 Dokumentasi

- **[Tutorial Lengkap](TUTORIAL.md)** - Panduan step-by-step untuk pemula
- **[Troubleshooting](TUTORIAL.md#10-troubleshooting)** - Solusi masalah umum
- **[FAQ](TUTORIAL.md#11-faq)** - Pertanyaan yang sering diajukan

## 🔧 Commands

| Command | Description |
|---------|-------------|
| `/start` | Mulai bot |
| `/help` | Bantuan |
| `/kategori` | Daftar kategori |
| `/laporan [periode]` | Laporan keuangan |
| `/budget [kategori] [jumlah]` | Set budget |
| `/notes [teks]` | Simpan catatan |
| `/reminder [teks]` | Buat pengingat |
| `/done [id]` | Tandai selesai |
| `/export` | Link Google Sheets |

## 🤝 Kontribusi

Pull requests welcome! Untuk perubahan besar, buka issue dulu.

## 📄 License

MIT
