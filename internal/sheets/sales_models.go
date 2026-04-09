package sheets

import (
	"fmt"
	"strings"
	"time"
)

// SalesItem represents an item with fixed purchase price from supplier
type SalesItem struct {
	ID        string
	Nama      string
	HargaBeli float64
	Satuan    string
	Catatan   string
	CreatedAt time.Time
}

// ToRow converts SalesItem into Items sheet row format
func (i *SalesItem) ToRow() []interface{} {
	return []interface{}{
		i.ID,
		i.Nama,
		i.HargaBeli,
		i.Satuan,
		i.Catatan,
		i.CreatedAt.In(WIB).Format("02/01/2006 15:04"),
	}
}

// SalesItemFromRow parses a sheet row back into SalesItem
func SalesItemFromRow(row []interface{}) (*SalesItem, error) {
	if len(row) < 5 {
		return nil, fmt.Errorf("invalid row: expected at least 5 columns, got %d", len(row))
	}

	hargaBeli, err := cellFloat64(row[2])
	if err != nil {
		return nil, fmt.Errorf("invalid harga_beli: %w", err)
	}

	return &SalesItem{
		ID:        cellString(row[0]),
		Nama:      cellString(row[1]),
		HargaBeli: hargaBeli,
		Satuan:    cellString(row[3]),
		Catatan:   cellString(row[4]),
	}, nil
}

// SalesCustomer represents a customer with credit terms
type SalesCustomer struct {
	ID         string
	Nama       string
	Alamat     string
	Telepon    string
	JatuhTempo int    // days: 7, 14, 30
	Payment    string // cash or credit
	Catatan    string
	CreatedAt  time.Time
}

// ToRow converts SalesCustomer into Customers sheet row format
func (c *SalesCustomer) ToRow() []interface{} {
	return []interface{}{
		c.ID,
		c.Nama,
		c.Alamat,
		c.Telepon,
		c.JatuhTempo,
		c.Payment,
		c.Catatan,
		c.CreatedAt.In(WIB).Format("02/01/2006 15:04"),
	}
}

// SalesCustomerFromRow parses a sheet row back into SalesCustomer
func SalesCustomerFromRow(row []interface{}) (*SalesCustomer, error) {
	if len(row) < 6 {
		return nil, fmt.Errorf("invalid row: expected at least 6 columns, got %d", len(row))
	}

	jatuhTempo, err := cellInt(row[4])
	if err != nil {
		jatuhTempo = 14 // default
	}

	return &SalesCustomer{
		ID:         cellString(row[0]),
		Nama:       cellString(row[1]),
		Alamat:     cellString(row[2]),
		Telepon:    cellString(row[3]),
		JatuhTempo: jatuhTempo,
		Payment:    cellString(row[5]),
		Catatan:    cellStringOrEmpty(row, 6),
	}, nil
}

// CustomerPricing represents selling price per customer per item
type CustomerPricing struct {
	CustomerID string
	ItemID     string
	HargaJual  float64
}

// ToRow converts CustomerPricing into CustomerPricing sheet row format
func (p *CustomerPricing) ToRow() []interface{} {
	return []interface{}{
		p.CustomerID,
		p.ItemID,
		p.HargaJual,
	}
}

// CustomerPricingFromRow parses a sheet row back into CustomerPricing
func CustomerPricingFromRow(row []interface{}) (*CustomerPricing, error) {
	if len(row) < 3 {
		return nil, fmt.Errorf("invalid row: expected 3 columns, got %d", len(row))
	}

	hargaJual, err := cellFloat64(row[2])
	if err != nil {
		return nil, fmt.Errorf("invalid harga_jual: %w", err)
	}

	return &CustomerPricing{
		CustomerID: cellString(row[0]),
		ItemID:     cellString(row[1]),
		HargaJual:  hargaJual,
	}, nil
}

// SalesTransaction represents a sales transaction
type SalesTransaction struct {
	ID           string
	Tanggal      time.Time
	CustomerID   string
	CustomerNama string
	Alamat       string
	ItemID       string
	ItemNama     string
	Qty          int
	Satuan       string
	HargaBeli    float64
	HargaJual    float64
	TotalBeli    float64
	TotalJual    float64
	Profit       float64
	Payment      string
	Catatan      string
}

// ToRow converts SalesTransaction into SalesTransactions sheet row format
func (t *SalesTransaction) ToRow() []interface{} {
	return []interface{}{
		t.ID,
		t.Tanggal.In(WIB).Format("02/01/2006"),
		t.Tanggal.In(WIB).Format("15:04"),
		t.CustomerID,
		t.CustomerNama,
		t.Alamat,
		t.ItemID,
		t.ItemNama,
		t.Qty,
		t.Satuan,
		t.HargaBeli,
		t.HargaJual,
		t.TotalBeli,
		t.TotalJual,
		t.Profit,
		t.Payment,
		t.Catatan,
	}
}

// SalesTransactionFromRow parses a sheet row back into SalesTransaction
func SalesTransactionFromRow(row []interface{}) (*SalesTransaction, error) {
	if len(row) < 16 {
		return nil, fmt.Errorf("invalid row: expected at least 16 columns, got %d", len(row))
	}

	dateStr := cellString(row[1])
	timeStr := cellString(row[2])
	dt, err := time.ParseInLocation("02/01/2006 15:04", strings.TrimSpace(dateStr)+" "+strings.TrimSpace(timeStr), WIB)
	if err != nil {
		return nil, fmt.Errorf("invalid date/time: %w", err)
	}

	qty, err := cellInt(row[8])
	if err != nil {
		qty = 1
	}

	hargaBeli, _ := cellFloat64(row[10])
	hargaJual, _ := cellFloat64(row[11])
	totalBeli, _ := cellFloat64(row[12])
	totalJual, _ := cellFloat64(row[13])
	profit, _ := cellFloat64(row[14])

	return &SalesTransaction{
		ID:           cellString(row[0]),
		Tanggal:      dt,
		CustomerID:   cellString(row[3]),
		CustomerNama: cellString(row[4]),
		Alamat:       cellString(row[5]),
		ItemID:       cellString(row[6]),
		ItemNama:     cellString(row[7]),
		Qty:          qty,
		Satuan:       cellString(row[9]),
		HargaBeli:    hargaBeli,
		HargaJual:    hargaJual,
		TotalBeli:    totalBeli,
		TotalJual:    totalJual,
		Profit:       profit,
		Payment:      cellString(row[15]),
		Catatan:      cellStringOrEmpty(row, 16),
	}, nil
}

// Payable represents debt to supplier
type Payable struct {
	ID           string
	Tanggal      time.Time
	SupplierName string
	ItemNama     string
	Qty          int
	Jumlah       float64
	JatuhTempo   time.Time
	Status       string // pending or lunas
	TanggalBayar *time.Time
	TransID      string
}

// ToRow converts Payable into Payables sheet row format
func (p *Payable) ToRow() []interface{} {
	tanggalBayar := ""
	if p.TanggalBayar != nil {
		tanggalBayar = p.TanggalBayar.In(WIB).Format("02/01/2006")
	}
	return []interface{}{
		p.ID,
		p.Tanggal.In(WIB).Format("02/01/2006"),
		p.SupplierName,
		p.ItemNama,
		p.Qty,
		p.Jumlah,
		p.JatuhTempo.In(WIB).Format("02/01/2006"),
		p.Status,
		tanggalBayar,
		p.TransID,
	}
}

// PayableFromRow parses a sheet row back into Payable
func PayableFromRow(row []interface{}) (*Payable, error) {
	if len(row) < 9 {
		return nil, fmt.Errorf("invalid row: expected at least 9 columns, got %d", len(row))
	}

	tanggal, _ := time.ParseInLocation("02/01/2006", cellString(row[1]), WIB)
	jatuhTempo, _ := time.ParseInLocation("02/01/2006", cellString(row[6]), WIB)

	qty, _ := cellInt(row[4])
	jumlah, _ := cellFloat64(row[5])

	var tanggalBayar *time.Time
	if tb := cellString(row[8]); tb != "" {
		if parsed, err := time.ParseInLocation("02/01/2006", tb, WIB); err == nil {
			tanggalBayar = &parsed
		}
	}

	return &Payable{
		ID:           cellString(row[0]),
		Tanggal:      tanggal,
		SupplierName: cellString(row[2]),
		ItemNama:     cellString(row[3]),
		Qty:          qty,
		Jumlah:       jumlah,
		JatuhTempo:   jatuhTempo,
		Status:       cellString(row[7]),
		TanggalBayar: tanggalBayar,
		TransID:      cellStringOrEmpty(row, 9),
	}, nil
}

// Receivable represents customer debt
type Receivable struct {
	ID            string
	Tanggal       time.Time
	CustomerID    string
	CustomerNama  string
	ItemNama      string
	Qty           int
	Jumlah        float64
	JatuhTempo    time.Time
	HariTempo     int
	Status        string // pending or lunas
	TanggalBayar  *time.Time
	TransID       string
	ReminderSent  bool
}

// ToRow converts Receivable into Receivables sheet row format
func (r *Receivable) ToRow() []interface{} {
	tanggalBayar := ""
	if r.TanggalBayar != nil {
		tanggalBayar = r.TanggalBayar.In(WIB).Format("02/01/2006")
	}
	return []interface{}{
		r.ID,
		r.Tanggal.In(WIB).Format("02/01/2006"),
		r.CustomerID,
		r.CustomerNama,
		r.ItemNama,
		r.Qty,
		r.Jumlah,
		r.JatuhTempo.In(WIB).Format("02/01/2006"),
		r.HariTempo,
		r.Status,
		tanggalBayar,
		r.TransID,
		r.ReminderSent,
	}
}

// ReceivableFromRow parses a sheet row back into Receivable
func ReceivableFromRow(row []interface{}) (*Receivable, error) {
	if len(row) < 11 {
		return nil, fmt.Errorf("invalid row: expected at least 11 columns, got %d", len(row))
	}

	tanggal, _ := time.ParseInLocation("02/01/2006", cellString(row[1]), WIB)
	jatuhTempo, _ := time.ParseInLocation("02/01/2006", cellString(row[7]), WIB)

	qty, _ := cellInt(row[5])
	jumlah, _ := cellFloat64(row[6])
	hariTempo, _ := cellInt(row[8])

	var tanggalBayar *time.Time
	if tb := cellString(row[10]); tb != "" {
		if parsed, err := time.ParseInLocation("02/01/2006", tb, WIB); err == nil {
			tanggalBayar = &parsed
		}
	}

	reminderSent := false
	if len(row) > 12 {
		if rs, err := cellBool(row[12]); err == nil {
			reminderSent = rs
		}
	}

	return &Receivable{
		ID:           cellString(row[0]),
		Tanggal:      tanggal,
		CustomerID:   cellString(row[2]),
		CustomerNama: cellString(row[3]),
		ItemNama:     cellString(row[4]),
		Qty:          qty,
		Jumlah:       jumlah,
		JatuhTempo:   jatuhTempo,
		HariTempo:    hariTempo,
		Status:       cellString(row[9]),
		TanggalBayar: tanggalBayar,
		TransID:      cellStringOrEmpty(row, 11),
		ReminderSent: reminderSent,
	}, nil
}

// SalesConfig holds configuration for sales module
type SalesConfig struct {
	SupplierName              string
	SupplierPayDay            int    // 25 = tanggal 25
	DefaultCreditDays         int    // fallback jatuh tempo
	ReminderTime              string // "08:00"
	ReminderEnabled           bool
	WhatsAppToCustomerEnabled bool
	WhatsAppReminderDays      int // H-1, H-3, dll
}

// DefaultSalesConfig returns default configuration
func DefaultSalesConfig() *SalesConfig {
	return &SalesConfig{
		SupplierName:              "Toko Supplier",
		SupplierPayDay:            25,
		DefaultCreditDays:         14,
		ReminderTime:              "08:00",
		ReminderEnabled:           true,
		WhatsAppToCustomerEnabled: false,
		WhatsAppReminderDays:      1,
	}
}

// Helper functions

func cellBool(v interface{}) (bool, error) {
	switch x := v.(type) {
	case bool:
		return x, nil
	case string:
		s := strings.ToLower(strings.TrimSpace(x))
		return s == "true" || s == "1" || s == "ya", nil
	default:
		return false, fmt.Errorf("unsupported type %T", v)
	}
}

func cellStringOrEmpty(row []interface{}, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return cellString(row[idx])
}
