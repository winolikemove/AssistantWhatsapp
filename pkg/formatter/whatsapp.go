package formatter

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

func FormatExpenseRecorded(id, description, category string, amount float64) string {
	return fmt.Sprintf(
		"✅ *Pengeluaran Dicatat!*\n\n🆔 ID: %s\n📝 Deskripsi: %s\n📂 Kategori: %s\n💰 Jumlah: %s\n\n📅 %s",
		id, safe(description), safe(category), formatIDR(amount), nowWIBString(),
	)
}

func FormatIncomeRecorded(id, description, category string, amount float64) string {
	return fmt.Sprintf(
		"✅ *Pemasukan Dicatat!*\n\n🆔 ID: %s\n📝 Deskripsi: %s\n📂 Kategori: %s\n💵 Jumlah: %s\n\n📅 %s",
		id, safe(description), safe(category), formatIDR(amount), nowWIBString(),
	)
}

func FormatDailyReport(date string, totalIncome, totalExpense float64, categories map[string]float64) string {
	return formatReport("📊 *Laporan Harian*", date, totalIncome, totalExpense, categories)
}

func FormatWeeklyReport(dateRange string, totalIncome, totalExpense float64, categories map[string]float64) string {
	return formatReport("📊 *Laporan Mingguan*", dateRange, totalIncome, totalExpense, categories)
}

func FormatMonthlyReport(month string, totalIncome, totalExpense float64, categories map[string]float64) string {
	return formatReport("📊 *Laporan Bulanan*", month, totalIncome, totalExpense, categories)
}

func FormatWelcome() string {
	return "👋 *Halo! Selamat datang di WA AI Assistant!*\n\n" +
		"Saya asisten keuangan pribadi Anda.\n\n" +
		"✨ Fitur utama:\n" +
		"• Catat pengeluaran/pemasukan pakai bahasa natural\n" +
		"• Laporan harian, mingguan, bulanan\n" +
		"• Budget per kategori + peringatan\n" +
		"• Catatan cepat dan export Google Sheets\n\n" +
		"Ketik */help* untuk melihat semua perintah."
}

func FormatHelp() string {
	return "📖 *Daftar Perintah:*\n\n" +
		"💰 *Keuangan*\n" +
		"• Kirim pesan seperti \"beli ayam crispy 16k\" untuk mencatat pengeluaran\n" +
		"• /laporan [hari ini|minggu ini|bulan ini] — Lihat laporan\n" +
		"• /budget [kategori] [jumlah] — Atur budget\n" +
		"• /edit [ID] [field] [nilai] — Edit transaksi\n" +
		"• /hapus [ID] — Hapus transaksi\n\n" +
		"📝 *Catatan*\n" +
		"• /notes [teks] — Simpan catatan cepat\n\n" +
		"⏰ *Pengingat*\n" +
		"• /reminder [teks] — Buat pengingat (contoh: /reminder tanggal 26 maret bayar vps)\n" +
		"• /done [ID] — Tandai pengingat sudah selesai\n\n" +
		"📂 *Lainnya*\n" +
		"• /kategori — Lihat daftar kategori\n" +
		"• /export — Dapatkan link Google Sheets\n" +
		"• /help — Tampilkan bantuan ini"
}

func FormatBudgetAlert(category string, budget, spent, remaining float64) string {
	if remaining >= 0 {
		return fmt.Sprintf(
			"⚠️ *Peringatan Budget!*\nKategori *%s* hampir habis.\nBudget: %s\nTerpakai: %s\nSisa: %s",
			safe(category), formatIDR(budget), formatIDR(spent), formatIDR(remaining),
		)
	}

	return fmt.Sprintf(
		"🚨 *Peringatan Budget!*\nKategori *%s* sudah melebihi budget!\nBudget: %s\nTerpakai: %s\nLebih: %s",
		safe(category), formatIDR(budget), formatIDR(spent), formatIDR(math.Abs(remaining)),
	)
}

func FormatBudgetSet(category string, amount float64) string {
	return fmt.Sprintf("✅ Budget *%s* diatur ke *%s* per bulan.", safe(category), formatIDR(amount))
}

func FormatNoteSaved(note string) string {
	return fmt.Sprintf("✅ Catatan disimpan!\n📝 \"%s\"", safe(note))
}

func FormatTransactionDeleted(id string) string {
	return fmt.Sprintf("✅ Transaksi *%s* berhasil dihapus.", safe(id))
}

func FormatTransactionEdited(id, field, oldValue, newValue string) string {
	return fmt.Sprintf(
		"✅ Transaksi *%s* diperbarui!\n📝 %s: %s → %s",
		safe(id), safe(field), safe(oldValue), safe(newValue),
	)
}

func FormatCategories(expenseCategories, incomeCategories []string) string {
	return fmt.Sprintf(
		"📂 *Daftar Kategori*\n\n💸 *Pengeluaran:*\n%s\n\n💵 *Pemasukan:*\n%s",
		"• "+strings.Join(expenseCategories, "\n• "),
		"• "+strings.Join(incomeCategories, "\n• "),
	)
}

func FormatExport(url string) string {
	return fmt.Sprintf("📊 *Data Keuangan Anda*\n🔗 Link Google Sheets:\n%s", safe(url))
}

func FormatError(msg string) string {
	return fmt.Sprintf("❌ *Error:* %s", safe(msg))
}

func FormatConfirmation(description string, txType string, amount float64, category string) string {
	kind := strings.ToLower(strings.TrimSpace(txType))
	if kind == "" {
		kind = "transaksi"
	}
	return fmt.Sprintf(
		"🤔 Saya catat sebagai *%s*:\n📝 %s\n📂 %s\n💰 %s\n\nBenar? Ketik *ya* untuk konfirmasi atau *bukan* untuk membatalkan.",
		kind, safe(description), safe(category), formatIDR(amount),
	)
}

func formatReport(title, period string, totalIncome, totalExpense float64, categories map[string]float64) string {
	net := totalIncome - totalExpense
	top := topCategories(categories, 5)

	var details strings.Builder
	if len(top) == 0 {
		details.WriteString("• Belum ada data kategori")
	} else {
		for _, c := range top {
			details.WriteString(fmt.Sprintf("• %s: %s\n", c.Name, formatIDR(c.Amount)))
		}
	}

	return fmt.Sprintf(
		"%s\n📅 %s\n\n💵 Total Pemasukan: %s\n💸 Total Pengeluaran: %s\n💰 Saldo Bersih: %s\n\n📂 *Rincian Kategori:*\n%s",
		title, safe(period), formatIDR(totalIncome), formatIDR(totalExpense), formatIDR(net), strings.TrimSpace(details.String()),
	)
}

type categoryAmount struct {
	Name   string
	Amount float64
}

func topCategories(m map[string]float64, n int) []categoryAmount {
	if len(m) == 0 || n <= 0 {
		return nil
	}

	items := make([]categoryAmount, 0, len(m))
	for k, v := range m {
		items = append(items, categoryAmount{Name: k, Amount: v})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Amount == items[j].Amount {
			return items[i].Name < items[j].Name
		}
		return items[i].Amount > items[j].Amount
	})

	if len(items) > n {
		items = items[:n]
	}
	return items
}

func formatIDR(amount float64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	amount = math.Round(amount*100) / 100
	intPart := int64(amount)
	fracPart := int(math.Round((amount - float64(intPart)) * 100))

	intText := withThousandDots(fmt.Sprintf("%d", intPart))
	if fracPart == 0 {
		return fmt.Sprintf("%sRp %s", sign, intText)
	}
	return fmt.Sprintf("%sRp %s,%02d", sign, intText, fracPart)
}

func withThousandDots(s string) string {
	if len(s) <= 3 {
		return s
	}
	prefix := len(s) % 3
	if prefix == 0 {
		prefix = 3
	}

	var b strings.Builder
	b.WriteString(s[:prefix])
	for i := prefix; i < len(s); i += 3 {
		b.WriteString(".")
		b.WriteString(s[i : i+3])
	}
	return b.String()
}

func nowWIBString() string {
	wib := time.FixedZone("WIB", 7*60*60)
	return time.Now().In(wib).Format("02 Jan 2006, 15:04 WIB")
}

func safe(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "-"
	}
	return s
}

// Safe is an exported version of safe for use by other packages
func Safe(s string) string {
	return safe(s)
}

// === SALES FORMATTERS ===

func FormatSalesItemAdded(nama, satuan string, hargaBeli float64) string {
	return fmt.Sprintf(
		"✅ *Item Ditambahkan!*\n\n📦 Nama: %s\n💰 Harga Beli: %s/%s\n\nSekarang set harga jual untuk customer dengan:\nset harga %s untuk [customer] [harga_jual]",
		safe(nama), formatIDR(hargaBeli), safe(satuan), safe(nama),
	)
}

func FormatSalesCustomerAdded(nama, alamat string, jatuhTempo int, payment string) string {
	return fmt.Sprintf(
		"✅ *Customer Ditambahkan!*\n\n👤 Nama: %s\n📍 Alamat: %s\n⏰ Jatuh Tempo: %d hari\n💳 Payment: %s",
		safe(nama), safe(alamat), jatuhTempo, safe(payment),
	)
}

func FormatCustomerPricingSet(customerNama, itemNama string, hargaJual float64) string {
	return fmt.Sprintf(
		"✅ *Harga Jual Diset!*\n\n👤 Customer: %s\n📦 Item: %s\n💵 Harga Jual: %s",
		safe(customerNama), safe(itemNama), formatIDR(hargaJual),
	)
}

func FormatProfitReport(period string, totalProfit, totalSales, totalCost float64, transactionCount int, items, customers map[string]float64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📊 *Laporan Profit - %s*\n\n", safe(period)))
	sb.WriteString(fmt.Sprintf("✨ Total Profit: %s\n", formatIDR(totalProfit)))
	sb.WriteString(fmt.Sprintf("💵 Total Penjualan: %s\n", formatIDR(totalSales)))
	sb.WriteString(fmt.Sprintf("💰 Total Modal: %s\n", formatIDR(totalCost)))
	sb.WriteString(fmt.Sprintf("📦 Jumlah Transaksi: %d\n\n", transactionCount))

	if len(items) > 0 {
		sb.WriteString("📂 *Profit per Item:*\n")
		for item, profit := range topCategories(items, 5) {
			sb.WriteString(fmt.Sprintf("• %s: %s\n", item, formatIDR(profit)))
		}
		sb.WriteString("\n")
	}

	if len(customers) > 0 {
		sb.WriteString("👤 *Profit per Customer:*\n")
		for cust, profit := range topCategories(customers, 5) {
			sb.WriteString(fmt.Sprintf("• %s: %s\n", cust, formatIDR(profit)))
		}
	}

	return sb.String()
}

func FormatReceivableSummary(dueToday, totalOverdue, totalPending float64, dueTodayList, overdueList []ReceivableItem) string {
	var sb strings.Builder
	sb.WriteString("💰 *Ringkasan Piutang*\n\n")

	if len(dueTodayList) > 0 {
		sb.WriteString("📌 *Jatuh Tempo Hari Ini:*\n")
		for _, r := range dueTodayList {
			sb.WriteString(fmt.Sprintf("• %s: %s (%s)\n", r.Customer, formatIDR(r.Jumlah), r.Item))
		}
		sb.WriteString(fmt.Sprintf("   *Total: %s*\n\n", formatIDR(dueToday)))
	}

	if len(overdueList) > 0 {
		sb.WriteString("⚠️ *Overdue (Terlambat):*\n")
		for _, r := range overdueList {
			sb.WriteString(fmt.Sprintf("• %s: %s (telat %d hari)\n", r.Customer, formatIDR(r.Jumlah), r.DaysLate))
		}
		sb.WriteString(fmt.Sprintf("   *Total: %s*\n\n", formatIDR(totalOverdue)))
	}

	sb.WriteString(fmt.Sprintf("💰 *TOTAL PIUTANG: %s*", formatIDR(totalPending)))

	return sb.String()
}

func FormatPayableSummary(totalPending float64, payables []PayableItem, dueDate string) string {
	var sb strings.Builder
	sb.WriteString("💸 *Ringkasan Hutang ke Supplier*\n\n")

	if len(payables) > 0 {
		sb.WriteString("📋 *Daftar Hutang:*\n")
		for _, p := range payables {
			sb.WriteString(fmt.Sprintf("• %s: %s\n", p.Item, formatIDR(p.Jumlah)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("💰 *Total Hutang: %s*\n", formatIDR(totalPending)))
	sb.WriteString(fmt.Sprintf("📅 Jatuh Tempo: %s", safe(dueDate)))

	return sb.String()
}

func FormatReceivablePaid(customerNama string, jumlah float64) string {
	return fmt.Sprintf("✅ *Piutang Dibayar!*\n\n👤 Customer: %s\n💵 Jumlah: %s", safe(customerNama), formatIDR(jumlah))
}

func FormatPayablePaid(jumlah float64) string {
	return fmt.Sprintf("✅ *Hutang Dibayar!*\n\n💵 Jumlah: %s", formatIDR(jumlah))
}

func FormatWAReminderToggled(enabled bool) string {
	status := "nonaktif"
	if enabled {
		status = "aktif"
	}
	return fmt.Sprintf("✅ *Reminder WA ke Customer %s*", status)
}

func FormatSalesItemsList(items []SalesItemInfo) string {
	if len(items) == 0 {
		return "📦 *Daftar Item*\n\nBelum ada item. Tambah dengan:\ntambah item [nama] [harga_beli] [satuan]"
	}

	var sb strings.Builder
	sb.WriteString("📦 *Daftar Item*\n\n")
	for _, i := range items {
		sb.WriteString(fmt.Sprintf("• %s: %s/%s\n", i.Nama, formatIDR(i.HargaBeli), i.Satuan))
	}
	return sb.String()
}

func FormatSalesCustomersList(customers []SalesCustomerInfo) string {
	if len(customers) == 0 {
		return "👤 *Daftar Customer*\n\nBelum ada customer. Tambah dengan:\ntambah customer [nama] [alamat] [jatuh_tempo] [payment]"
	}

	var sb strings.Builder
	sb.WriteString("👤 *Daftar Customer*\n\n")
	for _, c := range customers {
		sb.WriteString(fmt.Sprintf("• %s (%s, %d hari)\n", c.Nama, c.Payment, c.JatuhTempo))
	}
	return sb.String()
}

// Helper types for sales formatters
type ReceivableItem struct {
	Customer string
	Item     string
	Jumlah   float64
	DaysLate int
}

type PayableItem struct {
	Item   string
	Jumlah float64
}

type SalesItemInfo struct {
	Nama      string
	HargaBeli float64
	Satuan    string
}

type SalesCustomerInfo struct {
	Nama       string
	JatuhTempo int
	Payment    string
}
