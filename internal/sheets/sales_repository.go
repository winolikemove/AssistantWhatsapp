package sheets

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// SalesRepository defines the contract for all sales-related Google Sheets operations
type SalesRepository interface {
	// Item operations
	GetItemByNama(ctx context.Context, nama string) (*SalesItem, error)
	GetItemByID(ctx context.Context, id string) (*SalesItem, error)
	GetAllItems(ctx context.Context) ([]SalesItem, error)
	AppendItem(ctx context.Context, item *SalesItem) error

	// Customer operations
	GetCustomerByNama(ctx context.Context, nama string) (*SalesCustomer, error)
	GetCustomerByID(ctx context.Context, id string) (*SalesCustomer, error)
	GetAllCustomers(ctx context.Context) ([]SalesCustomer, error)
	AppendCustomer(ctx context.Context, customer *SalesCustomer) error
	UpdateCustomer(ctx context.Context, rowIndex int, customer *SalesCustomer) error

	// Customer pricing operations
	GetCustomerPricing(ctx context.Context, customerID, itemID string) (*CustomerPricing, error)
	SetCustomerPricing(ctx context.Context, pricing *CustomerPricing) error
	GetAllCustomerPricings(ctx context.Context) ([]CustomerPricing, error)

	// Sales transaction operations
	AppendSalesTransaction(ctx context.Context, tx *SalesTransaction) error
	GetSalesTransactions(ctx context.Context, tabName string) ([]SalesTransaction, error)
	GetSalesTransactionByID(ctx context.Context, id string) (*SalesTransaction, int, string, error)

	// Payable operations
	AppendPayable(ctx context.Context, payable *Payable) error
	GetPendingPayables(ctx context.Context) ([]Payable, error)
	GetPayableByID(ctx context.Context, id string) (*Payable, int, error)
	UpdatePayable(ctx context.Context, rowIndex int, payable *Payable) error

	// Receivable operations
	AppendReceivable(ctx context.Context, receivable *Receivable) error
	GetPendingReceivables(ctx context.Context) ([]Receivable, error)
	GetReceivablesDueToday(ctx context.Context) ([]Receivable, error)
	GetOverdueReceivables(ctx context.Context) ([]Receivable, error)
	GetReceivableByID(ctx context.Context, id string) (*Receivable, int, error)
	UpdateReceivable(ctx context.Context, rowIndex int, receivable *Receivable) error
	GetReceivablesByCustomer(ctx context.Context, customerNama string) ([]Receivable, error)

	// Tab initialization
	InitSalesTabs(ctx context.Context) error
}

// Compile-time interface check
var _ SalesRepository = (*salesRepository)(nil)

type salesRepository struct {
	service       *GoogleSheetRepository
	salesIDGen    *IDGenerator
	pricingIDGen  *IDGenerator
	itemIDCounter int
	custIDCounter int
	mu            sync.Mutex
}

// NewSalesRepository creates a new SalesRepository
func NewSalesRepository(baseRepo *GoogleSheetRepository) SalesRepository {
	return &salesRepository{
		service:    baseRepo,
		salesIDGen: &IDGenerator{},
	}
}

// Tab names
const (
	TabSalesItems          = "SalesItems"
	TabSalesCustomers      = "SalesCustomers"
	TabCustomerPricing     = "CustomerPricing"
	TabSalesTransactions   = "SalesTransactions"
	TabPayables            = "Payables"
	TabReceivables         = "Receivables"
)

// === ITEM OPERATIONS ===

func (r *salesRepository) GetItemByNama(ctx context.Context, nama string) (*SalesItem, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabSalesItems)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet items: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 3 {
			continue
		}
		if matchNama(cellString(row[1]), nama) {
			return SalesItemFromRow(row)
		}
	}

	return nil, fmt.Errorf("item '%s' tidak ditemukan. Tambah dulu dengan: tambah item %s harga [harga]", nama, nama)
}

func (r *salesRepository) GetItemByID(ctx context.Context, id string) (*SalesItem, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabSalesItems)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet items: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}
		if cellString(row[0]) == id {
			return SalesItemFromRow(row)
		}
	}

	return nil, fmt.Errorf("item dengan ID '%s' tidak ditemukan", id)
}

func (r *salesRepository) GetAllItems(ctx context.Context) ([]SalesItem, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabSalesItems)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet items: %w", err)
	}

	var items []SalesItem
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}
		item, err := SalesItemFromRow(row)
		if err != nil {
			continue
		}
		items = append(items, *item)
	}

	return items, nil
}

func (r *salesRepository) AppendItem(ctx context.Context, item *SalesItem) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}

	r.mu.Lock()
	r.itemIDCounter++
	item.ID = fmt.Sprintf("ITEM-%03d", r.itemIDCounter)
	r.mu.Unlock()

	item.CreatedAt = time.Now()
	return r.service.appendRow(TabSalesItems, item.ToRow())
}

// === CUSTOMER OPERATIONS ===

func (r *salesRepository) GetCustomerByNama(ctx context.Context, nama string) (*SalesCustomer, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabSalesCustomers)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet customers: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}
		if matchNama(cellString(row[1]), nama) {
			return SalesCustomerFromRow(row)
		}
	}

	return nil, fmt.Errorf("customer '%s' tidak ditemukan. Tambah dulu dengan: tambah customer %s", nama, nama)
}

func (r *salesRepository) GetCustomerByID(ctx context.Context, id string) (*SalesCustomer, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabSalesCustomers)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet customers: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}
		if cellString(row[0]) == id {
			return SalesCustomerFromRow(row)
		}
	}

	return nil, fmt.Errorf("customer dengan ID '%s' tidak ditemukan", id)
}

func (r *salesRepository) GetAllCustomers(ctx context.Context) ([]SalesCustomer, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabSalesCustomers)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet customers: %w", err)
	}

	var customers []SalesCustomer
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}
		cust, err := SalesCustomerFromRow(row)
		if err != nil {
			continue
		}
		customers = append(customers, *cust)
	}

	return customers, nil
}

func (r *salesRepository) AppendCustomer(ctx context.Context, customer *SalesCustomer) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}

	r.mu.Lock()
	r.custIDCounter++
	customer.ID = fmt.Sprintf("CUST-%03d", r.custIDCounter)
	r.mu.Unlock()

	customer.CreatedAt = time.Now()
	return r.service.appendRow(TabSalesCustomers, customer.ToRow())
}

func (r *salesRepository) UpdateCustomer(ctx context.Context, rowIndex int, customer *SalesCustomer) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}
	return r.service.updateRow(TabSalesCustomers, rowIndex, customer.ToRow())
}

// === CUSTOMER PRICING OPERATIONS ===

func (r *salesRepository) GetCustomerPricing(ctx context.Context, customerID, itemID string) (*CustomerPricing, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabCustomerPricing)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet pricing: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}
		if cellString(row[0]) == customerID && cellString(row[1]) == itemID {
			return CustomerPricingFromRow(row)
		}
	}

	return nil, fmt.Errorf("harga jual untuk item ini ke customer belum diset. Set dengan: set harga [item] untuk [customer] [harga]")
}

func (r *salesRepository) SetCustomerPricing(ctx context.Context, pricing *CustomerPricing) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}

	// Check if exists, update or append
	rows, err := r.service.readSheet(TabCustomerPricing)
	if err != nil {
		return r.service.appendRow(TabCustomerPricing, pricing.ToRow())
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}
		if cellString(row[0]) == pricing.CustomerID && cellString(row[1]) == pricing.ItemID {
			// Update existing
			return r.service.updateRow(TabCustomerPricing, i, pricing.ToRow())
		}
	}

	// Append new
	return r.service.appendRow(TabCustomerPricing, pricing.ToRow())
}

func (r *salesRepository) GetAllCustomerPricings(ctx context.Context) ([]CustomerPricing, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabCustomerPricing)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet pricing: %w", err)
	}

	var pricings []CustomerPricing
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}
		p, err := CustomerPricingFromRow(row)
		if err != nil {
			continue
		}
		pricings = append(pricings, *p)
	}

	return pricings, nil
}

// === SALES TRANSACTION OPERATIONS ===

func (r *salesRepository) AppendSalesTransaction(ctx context.Context, tx *SalesTransaction) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}

	tx.ID = r.salesIDGen.Next(time.Now())
	tx.Tanggal = time.Now()
	return r.service.appendRow(TabSalesTransactions, tx.ToRow())
}

func (r *salesRepository) GetSalesTransactions(ctx context.Context, tabName string) ([]SalesTransaction, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(tabName)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet transactions: %w", err)
	}

	var transactions []SalesTransaction
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 15 {
			continue
		}
		tx, err := SalesTransactionFromRow(row)
		if err != nil {
			continue
		}
		transactions = append(transactions, *tx)
	}

	return transactions, nil
}

func (r *salesRepository) GetSalesTransactionByID(ctx context.Context, id string) (*SalesTransaction, int, string, error) {
	if r == nil || r.service == nil {
		return nil, 0, "", fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabSalesTransactions)
	if err != nil {
		return nil, 0, "", fmt.Errorf("gagal membaca sheet: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 1 {
			continue
		}
		if cellString(row[0]) == id {
			tx, err := SalesTransactionFromRow(row)
			if err != nil {
				return nil, i, TabSalesTransactions, err
			}
			return tx, i, TabSalesTransactions, nil
		}
	}

	return nil, 0, "", fmt.Errorf("transaksi dengan ID '%s' tidak ditemukan", id)
}

// === PAYABLE OPERATIONS ===

func (r *salesRepository) AppendPayable(ctx context.Context, payable *Payable) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}

	payable.ID = fmt.Sprintf("PAY-%s", time.Now().Format("20060102-150405"))
	payable.Status = "pending"
	return r.service.appendRow(TabPayables, payable.ToRow())
}

func (r *salesRepository) GetPendingPayables(ctx context.Context) ([]Payable, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabPayables)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet payables: %w", err)
	}

	var payables []Payable
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 8 {
			continue
		}
		if cellString(row[7]) == "pending" {
			p, err := PayableFromRow(row)
			if err != nil {
				continue
			}
			payables = append(payables, *p)
		}
	}

	return payables, nil
}

func (r *salesRepository) GetPayableByID(ctx context.Context, id string) (*Payable, int, error) {
	if r == nil || r.service == nil {
		return nil, 0, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabPayables)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal membaca sheet: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 1 {
			continue
		}
		if cellString(row[0]) == id {
			p, err := PayableFromRow(row)
			return p, i, err
		}
	}

	return nil, 0, fmt.Errorf("hutang dengan ID '%s' tidak ditemukan", id)
}

func (r *salesRepository) UpdatePayable(ctx context.Context, rowIndex int, payable *Payable) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}
	return r.service.updateRow(TabPayables, rowIndex, payable.ToRow())
}

// === RECEIVABLE OPERATIONS ===

func (r *salesRepository) AppendReceivable(ctx context.Context, receivable *Receivable) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}

	receivable.ID = fmt.Sprintf("REC-%s", time.Now().Format("20060102-150405"))
	receivable.Status = "pending"
	return r.service.appendRow(TabReceivables, receivable.ToRow())
}

func (r *salesRepository) GetPendingReceivables(ctx context.Context) ([]Receivable, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabReceivables)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet receivables: %w", err)
	}

	var receivables []Receivable
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 10 {
			continue
		}
		if cellString(row[9]) == "pending" {
			rec, err := ReceivableFromRow(row)
			if err != nil {
				continue
			}
			receivables = append(receivables, *rec)
		}
	}

	return receivables, nil
}

func (r *salesRepository) GetReceivablesDueToday(ctx context.Context) ([]Receivable, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabReceivables)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet receivables: %w", err)
	}

	today := time.Now().In(WIB).Format("02/01/2006")
	var receivables []Receivable

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 10 {
			continue
		}
		if cellString(row[9]) == "pending" && cellString(row[7]) == today {
			rec, err := ReceivableFromRow(row)
			if err != nil {
				continue
			}
			receivables = append(receivables, *rec)
		}
	}

	return receivables, nil
}

func (r *salesRepository) GetOverdueReceivables(ctx context.Context) ([]Receivable, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabReceivables)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet receivables: %w", err)
	}

	today := time.Now().In(WIB)
	var receivables []Receivable

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 10 {
			continue
		}
		if cellString(row[9]) == "pending" {
			rec, err := ReceivableFromRow(row)
			if err != nil {
				continue
			}
			if rec.JatuhTempo.Before(today) {
				receivables = append(receivables, *rec)
			}
		}
	}

	return receivables, nil
}

func (r *salesRepository) GetReceivableByID(ctx context.Context, id string) (*Receivable, int, error) {
	if r == nil || r.service == nil {
		return nil, 0, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabReceivables)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal membaca sheet: %w", err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 1 {
			continue
		}
		if cellString(row[0]) == id {
			rec, err := ReceivableFromRow(row)
			return rec, i, err
		}
	}

	return nil, 0, fmt.Errorf("piutang dengan ID '%s' tidak ditemukan", id)
}

func (r *salesRepository) UpdateReceivable(ctx context.Context, rowIndex int, receivable *Receivable) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}
	return r.service.updateRow(TabReceivables, rowIndex, receivable.ToRow())
}

func (r *salesRepository) GetReceivablesByCustomer(ctx context.Context, customerNama string) ([]Receivable, error) {
	if r == nil || r.service == nil {
		return nil, fmt.Errorf("repository is nil")
	}

	rows, err := r.service.readSheet(TabReceivables)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca sheet receivables: %w", err)
	}

	var receivables []Receivable
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 10 {
			continue
		}
		if matchNama(cellString(row[3]), customerNama) && cellString(row[9]) == "pending" {
			rec, err := ReceivableFromRow(row)
			if err != nil {
				continue
			}
			receivables = append(receivables, *rec)
		}
	}

	return receivables, nil
}

// === TAB INITIALIZATION ===

func (r *salesRepository) InitSalesTabs(ctx context.Context) error {
	if r == nil || r.service == nil {
		return fmt.Errorf("repository is nil")
	}

	// Init SalesItems tab
	if err := r.service.ensureTabExists(TabSalesItems); err != nil {
		return fmt.Errorf("gagal membuat tab SalesItems: %w", err)
	}
	if err := r.service.formatSalesHeaders(TabSalesItems, []string{"ID", "Nama", "Harga_Beli", "Satuan", "Catatan", "Created_At"}); err != nil {
		return err
	}

	// Init SalesCustomers tab
	if err := r.service.ensureTabExists(TabSalesCustomers); err != nil {
		return fmt.Errorf("gagal membuat tab SalesCustomers: %w", err)
	}
	if err := r.service.formatSalesHeaders(TabSalesCustomers, []string{"ID", "Nama", "Alamat", "Telepon", "Jatuh_Tempo", "Payment", "Catatan", "Created_At"}); err != nil {
		return err
	}

	// Init CustomerPricing tab
	if err := r.service.ensureTabExists(TabCustomerPricing); err != nil {
		return fmt.Errorf("gagal membuat tab CustomerPricing: %w", err)
	}
	if err := r.service.formatSalesHeaders(TabCustomerPricing, []string{"Customer_ID", "Item_ID", "Harga_Jual"}); err != nil {
		return err
	}

	// Init SalesTransactions tab
	if err := r.service.ensureTabExists(TabSalesTransactions); err != nil {
		return fmt.Errorf("gagal membuat tab SalesTransactions: %w", err)
	}
	if err := r.service.formatSalesHeaders(TabSalesTransactions, []string{"ID", "Tanggal", "Waktu", "Customer_ID", "Customer", "Alamat", "Item_ID", "Item", "Qty", "Satuan", "Harga_Beli", "Harga_Jual", "Total_Beli", "Total_Jual", "Profit", "Payment", "Catatan"}); err != nil {
		return err
	}

	// Init Payables tab
	if err := r.service.ensureTabExists(TabPayables); err != nil {
		return fmt.Errorf("gagal membuat tab Payables: %w", err)
	}
	if err := r.service.formatSalesHeaders(TabPayables, []string{"ID", "Tanggal", "Supplier", "Item", "Qty", "Jumlah", "Jatuh_Tempo", "Status", "Tanggal_Bayar", "Trans_ID"}); err != nil {
		return err
	}

	// Init Receivables tab
	if err := r.service.ensureTabExists(TabReceivables); err != nil {
		return fmt.Errorf("gagal membuat tab Receivables: %w", err)
	}
	if err := r.service.formatSalesHeaders(TabReceivables, []string{"ID", "Tanggal", "Customer_ID", "Customer", "Item", "Qty", "Jumlah", "Jatuh_Tempo", "Hari_Tempo", "Status", "Tanggal_Bayar", "Trans_ID", "Reminder_Sent"}); err != nil {
		return err
	}

	return nil
}

// Helper function to match names case-insensitively
func matchNama(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
