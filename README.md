# AssistantWhatsapp

WhatsApp AI Assistant dengan fitur keuangan pribadi dan sistem penjualan terintegrasi.

## Fitur

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

## Setup

### 1. Prerequisites
- Go 1.25.5+
- Google Cloud Project dengan Sheets API enabled
- Service Account credentials
- OpenAI API key atau compatible LLM

### 2. Configuration
```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

### 3. Google Sheets Setup
1. Buat Google Spreadsheet baru
2. Share spreadsheet ke Service Account email
3. Copy Spreadsheet ID ke `SHEETS_SPREADSHEET_ID`

### 4. Run
```bash
go mod tidy
go run ./cmd/bot
```

## Commands

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

## Contoh Penggunaan

### Keuangan
```
beli makan 25k
gaji 5jt
bayar listrik 150rb
```

### Penjualan
```
tambah item aussie bbq 15000 kg
tambah customer ambrogio jakarta 14 hari credit
set harga aussie bbq untuk ambrogio 18000
jual aussie bbq 50 kg ke ambrogio
```

### Laporan
```
laporan profit bulan ini
piutang
hutang
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `LLM_API_KEY` | Yes | API key untuk LLM |
| `LLM_BASE_URL` | Yes | Base URL LLM API |
| `LLM_MODEL` | Yes | Model name (e.g., gpt-4o-mini) |
| `GOOGLE_APPLICATION_CREDENTIALS` | Yes | Path ke credentials.json |
| `SHEETS_SPREADSHEET_ID` | Yes | Google Spreadsheet ID |
| `WHATSAPP_SESSION_DB_PATH` | Yes | Path untuk session database |
| `OWNER_PHONE_NUMBER` | Yes | Nomor WhatsApp owner (digits only) |
| `SUPPLIER_NAME` | No | Nama supplier (default: Toko Supplier) |
| `SUPPLIER_PAY_DAY` | No | Tanggal bayar supplier (default: 25) |
| `REMINDER_TIME` | No | Waktu reminder (default: 08:00) |
| `WA_REMINDER_TO_CUSTOMER` | No | Kirim reminder ke customer (default: false) |

## License

MIT
