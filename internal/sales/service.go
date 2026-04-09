package sales

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/verssache/AssistantWhatsapp/internal/sheets"
	"github.com/verssache/AssistantWhatsapp/pkg/formatter"
)

// Service provides sales business logic
type Service struct {
	repo   sheets.SalesRepository
	config *sheets.SalesConfig
}

// NewService creates a new sales service
func NewService(repo sheets.SalesRepository, config *sheets.SalesConfig) *Service {
	if config == nil {
		config = sheets.DefaultSalesConfig()
	}
	return &Service{
		repo:   repo,
		config: config,
	}
}

// TransactionRequest represents a sales transaction request
type TransactionRequest struct {
	ItemNama     string
	Qty          int
	CustomerNama string
	Catatan      string
}

// TransactionResult represents the result of a sales transaction
type TransactionResult struct {
	Transaction *sheets.SalesTransaction
	Payable     *sheets.Payable
	Receivable  *sheets.Receivable
}

// CreateTransaction creates a complete sales transaction
func (s *Service) CreateTransaction(ctx context.Context, req *TransactionRequest) (*TransactionResult, error) {
	if s == nil {
		return nil, fmt.Errorf("sales service is nil")
	}
	if s.repo == nil {
		return nil, fmt.Errorf("sales repository is nil")
	}
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	req.ItemNama = strings.TrimSpace(req.ItemNama)
	req.CustomerNama = strings.TrimSpace(req.CustomerNama)
	if req.ItemNama == "" {
		return nil, fmt.Errorf("nama item tidak boleh kosong")
	}
	if req.CustomerNama == "" {
		return nil, fmt.Errorf("nama customer tidak boleh kosong")
	}
	if req.Qty <= 0 {
		req.Qty = 1
	}

	// 1. Get item from database
	item, err := s.repo.GetItemByNama(ctx, req.ItemNama)
	if err != nil {
		return nil, fmt.Errorf("item tidak ditemukan: %s. Tambah dulu dengan: tambah item %s [harga_beli] [satuan]", req.ItemNama, req.ItemNama)
	}

	// 2. Get customer from database
	customer, err := s.repo.GetCustomerByNama(ctx, req.CustomerNama)
	if err != nil {
		return nil, fmt.Errorf("customer tidak ditemukan: %s. Tambah dulu dengan: tambah customer %s [alamat] [jatuh_tempo] [payment]", req.CustomerNama, req.CustomerNama)
	}

	// 3. Get customer pricing
	pricing, err := s.repo.GetCustomerPricing(ctx, customer.ID, item.ID)
	if err != nil {
		return nil, fmt.Errorf("harga jual untuk %s ke %s belum diset. Set dengan: set harga %s untuk %s [harga_jual]", item.Nama, customer.Nama, item.Nama, customer.Nama)
	}

	// 4. Calculate totals
	totalBeli := item.HargaBeli * float64(req.Qty)
	totalJual := pricing.HargaJual * float64(req.Qty)
	profit := totalJual - totalBeli

	now := time.Now().In(sheets.WIB)

	// 5. Create transaction
	tx := &sheets.SalesTransaction{
		Tanggal:      now,
		CustomerID:   customer.ID,
		CustomerNama: customer.Nama,
		Alamat:       customer.Alamat,
		ItemID:       item.ID,
		ItemNama:     item.Nama,
		Qty:          req.Qty,
		Satuan:       item.Satuan,
		HargaBeli:    item.HargaBeli,
		HargaJual:    pricing.HargaJual,
		TotalBeli:    totalBeli,
		TotalJual:    totalJual,
		Profit:       profit,
		Payment:      customer.Payment,
		Catatan:      req.Catatan,
	}

	if err := s.repo.AppendSalesTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("gagal menyimpan transaksi: %w", err)
	}

	result := &TransactionResult{
		Transaction: tx,
	}

	// 6. Create payable (hutang ke supplier)
	jatuhTempoSupplier := s.calculateSupplierDueDate(now)
	payable := &sheets.Payable{
		Tanggal:      now,
		SupplierName: s.config.SupplierName,
		ItemNama:     fmt.Sprintf("%s %d %s", item.Nama, req.Qty, item.Satuan),
		Qty:          req.Qty,
		Jumlah:       totalBeli,
		JatuhTempo:   jatuhTempoSupplier,
		TransID:      tx.ID,
	}
	s.repo.AppendPayable(ctx, payable)
	result.Payable = payable

	// 7. Create receivable if credit payment
	if customer.Payment == "credit" {
		jatuhTempoCustomer := now.AddDate(0, 0, customer.JatuhTempo)
		receivable := &sheets.Receivable{
			Tanggal:      now,
			CustomerID:   customer.ID,
			CustomerNama: customer.Nama,
			ItemNama:     fmt.Sprintf("%s %d %s", item.Nama, req.Qty, item.Satuan),
			Qty:          req.Qty,
			Jumlah:       totalJual,
			JatuhTempo:   jatuhTempoCustomer,
			HariTempo:    customer.JatuhTempo,
			TransID:      tx.ID,
		}
		s.repo.AppendReceivable(ctx, receivable)
		result.Receivable = receivable
	}

	return result, nil
}

// calculateSupplierDueDate calculates the due date for supplier payment (tanggal 25 bulan berjalan)
func (s *Service) calculateSupplierDueDate(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), s.config.SupplierPayDay, 0, 0, 0, 0, now.Location())
}

// AddItem adds a new item to the database
func (s *Service) AddItem(ctx context.Context, nama string, hargaBeli float64, satuan string) (*sheets.SalesItem, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("sales service is nil")
	}

	nama = strings.TrimSpace(nama)
	if nama == "" {
		return nil, fmt.Errorf("nama item tidak boleh kosong")
	}
	if hargaBeli <= 0 {
		return nil, fmt.Errorf("harga beli harus lebih dari 0")
	}
	satuan = strings.TrimSpace(satuan)
	if satuan == "" {
		satuan = "pcs"
	}

	item := &sheets.SalesItem{
		Nama:      nama,
		HargaBeli: hargaBeli,
		Satuan:    satuan,
	}

	if err := s.repo.AppendItem(ctx, item); err != nil {
		return nil, fmt.Errorf("gagal menambah item: %w", err)
	}

	return item, nil
}

// AddCustomer adds a new customer to the database
func (s *Service) AddCustomer(ctx context.Context, nama, alamat, telepon string, jatuhTempo int, payment string) (*sheets.SalesCustomer, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("sales service is nil")
	}

	nama = strings.TrimSpace(nama)
	if nama == "" {
		return nil, fmt.Errorf("nama customer tidak boleh kosong")
	}

	alamat = strings.TrimSpace(alamat)
	payment = strings.ToLower(strings.TrimSpace(payment))
	if payment != "cash" && payment != "credit" {
		payment = "credit"
	}

	if jatuhTempo <= 0 {
		jatuhTempo = s.config.DefaultCreditDays
	}

	customer := &sheets.SalesCustomer{
		Nama:       nama,
		Alamat:     alamat,
		Telepon:    telepon,
		JatuhTempo: jatuhTempo,
		Payment:    payment,
	}

	if err := s.repo.AppendCustomer(ctx, customer); err != nil {
		return nil, fmt.Errorf("gagal menambah customer: %w", err)
	}

	return customer, nil
}

// SetCustomerPricing sets selling price for an item to a customer
func (s *Service) SetCustomerPricing(ctx context.Context, customerNama, itemNama string, hargaJual float64) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("sales service is nil")
	}

	customer, err := s.repo.GetCustomerByNama(ctx, customerNama)
	if err != nil {
		return fmt.Errorf("customer tidak ditemukan: %s", customerNama)
	}

	item, err := s.repo.GetItemByNama(ctx, itemNama)
	if err != nil {
		return fmt.Errorf("item tidak ditemukan: %s", itemNama)
	}

	if hargaJual <= 0 {
		return fmt.Errorf("harga jual harus lebih dari 0")
	}

	pricing := &sheets.CustomerPricing{
		CustomerID: customer.ID,
		ItemID:     item.ID,
		HargaJual:  hargaJual,
	}

	if err := s.repo.SetCustomerPricing(ctx, pricing); err != nil {
		return fmt.Errorf("gagal menyimpan harga: %w", err)
	}

	return nil
}

// GetProfitReport generates profit report for a period
func (s *Service) GetProfitReport(ctx context.Context, period string) (*ProfitReport, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("sales service is nil")
	}

	transactions, err := s.repo.GetSalesTransactions(ctx, sheets.TabSalesTransactions)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil transaksi: %w", err)
	}

	now := time.Now().In(sheets.WIB)
	filtered := filterSalesByPeriod(transactions, now, period)

	report := &ProfitReport{
		Period:    period,
		Items:     make(map[string]float64),
		Customers: make(map[string]float64),
	}

	for _, tx := range filtered {
		report.TotalProfit += tx.Profit
		report.TotalSales += tx.TotalJual
		report.TotalCost += tx.TotalBeli
		report.Items[tx.ItemNama] += tx.Profit
		report.Customers[tx.CustomerNama] += tx.Profit
		report.TransactionCount++
	}

	return report, nil
}

// ProfitReport represents profit report data
type ProfitReport struct {
	Period           string
	TotalProfit      float64
	TotalSales       float64
	TotalCost        float64
	TransactionCount int
	Items            map[string]float64
	Customers        map[string]float64
}

// GetReceivableSummary returns summary of receivables
func (s *Service) GetReceivableSummary(ctx context.Context) (*ReceivableSummary, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("sales service is nil")
	}

	dueToday, _ := s.repo.GetReceivablesDueToday(ctx)
	overdue, _ := s.repo.GetOverdueReceivables(ctx)
	allPending, _ := s.repo.GetPendingReceivables(ctx)

	var totalDueToday, totalOverdue, totalPending float64
	for _, r := range dueToday {
		totalDueToday += r.Jumlah
	}
	for _, r := range overdue {
		totalOverdue += r.Jumlah
	}
	for _, r := range allPending {
		totalPending += r.Jumlah
	}

	return &ReceivableSummary{
		DueToday:      dueToday,
		Overdue:       overdue,
		TotalDueToday: totalDueToday,
		TotalOverdue:  totalOverdue,
		TotalPending:  totalPending,
	}, nil
}

// ReceivableSummary represents receivable summary
type ReceivableSummary struct {
	DueToday      []sheets.Receivable
	Overdue       []sheets.Receivable
	TotalDueToday float64
	TotalOverdue  float64
	TotalPending  float64
}

// GetPayableSummary returns summary of payables
func (s *Service) GetPayableSummary(ctx context.Context) (*PayableSummary, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("sales service is nil")
	}

	pending, _ := s.repo.GetPendingPayables(ctx)

	var totalPending float64
	for _, p := range pending {
		totalPending += p.Jumlah
	}

	return &PayableSummary{
		Payables:     pending,
		TotalPending: totalPending,
		DueDate:      s.calculateSupplierDueDate(time.Now()),
	}, nil
}

// PayableSummary represents payable summary
type PayableSummary struct {
	Payables     []sheets.Payable
	TotalPending float64
	DueDate      time.Time
}

// PayReceivable records a payment from customer
func (s *Service) PayReceivable(ctx context.Context, customerNama string, jumlah float64) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("sales service is nil")
	}

	receivables, err := s.repo.GetReceivablesByCustomer(ctx, customerNama)
	if err != nil {
		return fmt.Errorf("tidak ada piutang untuk customer %s", customerNama)
	}

	var total float64
	for _, r := range receivables {
		total += r.Jumlah
	}

	if total <= 0 {
		return fmt.Errorf("tidak ada piutang untuk customer %s", customerNama)
	}

	// Mark as paid (simplified - marks all receivables as paid)
	now := time.Now()
	for _, r := range receivables {
		if jumlah <= 0 {
			break
		}
		rec, idx, err := s.repo.GetReceivableByID(ctx, r.ID)
		if err != nil {
			continue
		}
		rec.Status = "lunas"
		rec.TanggalBayar = &now
		s.repo.UpdateReceivable(ctx, idx, rec)
		jumlah -= rec.Jumlah
	}

	return nil
}

// PayPayable records a payment to supplier
func (s *Service) PayPayable(ctx context.Context, jumlah float64) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("sales service is nil")
	}

	payables, err := s.repo.GetPendingPayables(ctx)
	if err != nil {
		return fmt.Errorf("tidak ada hutang")
	}

	now := time.Now()
	for _, p := range payables {
		if jumlah <= 0 {
			break
		}
		pay, idx, err := s.repo.GetPayableByID(ctx, p.ID)
		if err != nil {
			continue
		}
		pay.Status = "lunas"
		pay.TanggalBayar = &now
		s.repo.UpdatePayable(ctx, idx, pay)
		jumlah -= pay.Jumlah
	}

	return nil
}

// ToggleWAReminder toggles WhatsApp reminder to customer
func (s *Service) ToggleWAReminder(enabled bool) {
	if s != nil && s.config != nil {
		s.config.WhatsAppToCustomerEnabled = enabled
	}
}

// IsWAReminderEnabled returns whether WA reminder is enabled
func (s *Service) IsWAReminderEnabled() bool {
	if s == nil || s.config == nil {
		return false
	}
	return s.config.WhatsAppToCustomerEnabled
}

// ListItems returns all items
func (s *Service) ListItems(ctx context.Context) ([]sheets.SalesItem, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("sales service is nil")
	}
	return s.repo.GetAllItems(ctx)
}

// ListCustomers returns all customers
func (s *Service) ListCustomers(ctx context.Context) ([]sheets.SalesCustomer, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("sales service is nil")
	}
	return s.repo.GetAllCustomers(ctx)
}

// GetConfig returns current sales config
func (s *Service) GetConfig() *sheets.SalesConfig {
	if s == nil {
		return nil
	}
	return s.config
}

func filterSalesByPeriod(txs []sheets.SalesTransaction, now time.Time, period string) []sheets.SalesTransaction {
	period = strings.ToLower(strings.TrimSpace(period))

	switch period {
	case "hari ini", "harian", "daily":
		todayKey := now.Format("20060102")
		var out []sheets.SalesTransaction
		for _, tx := range txs {
			if tx.Tanggal.Format("20060102") == todayKey {
				out = append(out, tx)
			}
		}
		return out

	case "minggu ini", "mingguan", "weekly":
		start := now.AddDate(0, 0, -6)
		startKey := start.Format("20060102")
		endKey := now.Format("20060102")
		var out []sheets.SalesTransaction
		for _, tx := range txs {
			k := tx.Tanggal.Format("20060102")
			if k >= startKey && k <= endKey {
				out = append(out, tx)
			}
		}
		return out

	case "bulan ini", "bulanan", "monthly":
		monthKey := now.Format("200601")
		var out []sheets.SalesTransaction
		for _, tx := range txs {
			if tx.Tanggal.Format("200601") == monthKey {
				out = append(out, tx)
			}
		}
		return out

	default:
		return txs
	}
}

// FormatTransactionResult formats transaction result for WhatsApp
func FormatTransactionResult(result *TransactionResult) string {
	if result == nil || result.Transaction == nil {
		return formatter.FormatError("Transaksi tidak valid.")
	}

	tx := result.Transaction
	var sb strings.Builder

	sb.WriteString("✅ *Transaksi Tercatat!*\n\n")
	sb.WriteString(fmt.Sprintf("📦 Item: %s %d %s\n", tx.ItemNama, tx.Qty, tx.Satuan))
	sb.WriteString(fmt.Sprintf("👤 Customer: %s\n", tx.CustomerNama))
	sb.WriteString(fmt.Sprintf("📍 Alamat: %s\n", formatter.Safe(tx.Alamat)))
	sb.WriteString("─────────────────────\n")
	sb.WriteString(fmt.Sprintf("💰 Harga Beli: %s\n", formatter.FormatIDR(tx.TotalBeli)))
	sb.WriteString(fmt.Sprintf("💵 Harga Jual: %s\n", formatter.FormatIDR(tx.TotalJual)))
	sb.WriteString(fmt.Sprintf("✨ *Profit: %s*\n", formatter.FormatIDR(tx.Profit)))
	sb.WriteString("─────────────────────\n")
	sb.WriteString(fmt.Sprintf("💳 Payment: %s\n\n", tx.Payment))

	sb.WriteString(fmt.Sprintf("📌 Hutang ke Supplier: +%s\n", formatter.FormatIDR(tx.TotalBeli)))

	if result.Receivable != nil {
		sb.WriteString(fmt.Sprintf("📌 Piutang %s: +%s (jth tempo: %s)",
			result.Receivable.CustomerNama,
			formatter.FormatIDR(result.Receivable.Jumlah),
			result.Receivable.JatuhTempo.Format("02/01/2006")))
	}

	return sb.String()
}
